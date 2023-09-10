package url_crawler

import (
	"bytes"
	"context"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/html"

	"common/convert"
	"common/data/model"
	"common/data/store"
	"parser/internal/config"
	"parser/internal/services/connector"
	"parser/internal/services/crawler"
)

type UrlCrawler struct {
	log *logrus.Entry

	conn connector.Connector

	dataProvider store.DataProvider
}

func NewCrawler(cfg config.Config) crawler.MultiCrawler[model.Title] {
	return UrlCrawler{
		log:          cfg.Logging().WithField("service", "[URL-CRAWLER]"),
		conn:         connector.New(cfg),
		dataProvider: store.New(cfg),
	}
}

func hasMatchingID(attrs []html.Attribute) bool {
	for _, attr := range attrs {
		if attr.Key == "id" && strings.HasPrefix(attr.Val, "article") {
			return true
		}
	}
	return false
}

func extractArticle(doc *html.Node) (*html.Node, error) {
	var body *html.Node
	var htmlCrawler func(*html.Node)
	htmlCrawler = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "article" && hasMatchingID(node.Attr) {
			body = node
			return
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			htmlCrawler(child)
		}
	}
	htmlCrawler(doc)
	if body != nil {
		return body, nil
	}
	return nil, errors.New("Missing <article> in the node tree")
}

func collectText(n *html.Node, buf *bytes.Buffer) {
	if n.Type == html.TextNode {
		buf.WriteString(n.Data)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		collectText(c, buf)
	}
}

func (u UrlCrawler) Crawl(ctx context.Context, pendingTitles []model.Title) ([]crawler.ParsedBody, []int, []error) {

	outBodies := make([]crawler.ParsedBody, 0, len(pendingTitles))
	statusCodes := make([]int, 0, len(pendingTitles))
	errs := make([]error, 0, len(pendingTitles))

	for _, t := range pendingTitles {
		respBody, statusCode, err := u.conn.Poll(ctx, connector.PollParams{
			Url: convert.FromPtr(t.URL),
		})

		if err != nil {
			errs = append(errs, errors.Wrapf(err, "failed to poll url %s", convert.FromPtr(t.URL)))
			statusCodes = append(statusCodes, statusCode)
			continue
		}

		if statusCode != http.StatusOK {
			errs = append(errs, nil)
			statusCodes = append(statusCodes, statusCode)
			continue
		}

		rawHtml, err := html.Parse(respBody)
		if err != nil {
			errs = append(errs, errors.Wrap(err, "failed to parse response body"))
			statusCodes = append(statusCodes, statusCode)
			continue
		}

		rawArticle, err := extractArticle(rawHtml)
		if err != nil {
			errs = append(errs, errors.Wrap(err, "failed to extract article from webpage"))
			statusCodes = append(statusCodes, statusCode)
			continue
		}

		textBuf := bytes.NewBuffer(make([]byte, 0, 1000))
		collectText(rawArticle, textBuf)

		outBodies = append(outBodies, body{text: textBuf.String(), titleID: t.ID})
		statusCodes = append(statusCodes, statusCode)
		errs = append(errs, nil)
	}

	return outBodies, statusCodes, errs
}
