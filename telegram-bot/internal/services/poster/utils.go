package poster

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

func escapeKeepingHTML(text string) string {
	replacer := strings.NewReplacer(
		"&", "&amp;",
	)

	return replacer.Replace(text)
}

func downloadFile(fileUrl string) ([]byte, error) {
	//Get the response bytes from the url
	response, err := http.Get(fileUrl)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get file from url: %s", fileUrl)
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("Received %d response code", response.StatusCode))
	}

	fileBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read file body")
	}

	return fileBytes, nil
}
