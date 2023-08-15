package browse_ai_crawler

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"common"
	"common/convert"
	"common/data/model"
	"common/hash"
	"parser/internal/services/crawler"
)

var _ crawler.ParsedBody = body{}

const (
	TaskStatusSuccessful = "successful"
	TaskStatusInProgress = "in-progress"
	TaskStatusFailed     = "failed"
)

const coinTelegraphTimeLayout = "Jan 02, 2006"

type rawBody struct {
	Status  int32  `json:"statusCode"`
	Message string `json:"messageCode"`
	Result  struct {
		Id                string  `json:"id"`
		RobotId           string  `json:"robotId"`
		Status            string  `json:"status"`
		CreatedAt         *int64  `json:"createdAt"`
		StartedAt         *int64  `json:"startedAt"`
		FinishedAt        *int64  `json:"finishedAt"`
		UserFriendlyError *string `json:"userFriendlyError"`
		CapturedLists     *struct {
			Releases []body `json:"coin_telegraph_press_releases"`
		} `json:"capturedLists"`
	} `json:"result"`
}

type body struct {
	ReleaseDate  string `json:"release_date"`
	ReleaseURL   string `json:"release_url"`
	ReleaseTitle string `json:"release_title"`
	ReleaseDesc  string `json:"release_desc"`
}

func (b body) ToModel() any {
	return model.Title{
		Title:       &b.ReleaseTitle,
		Summary:     &b.ReleaseDesc,
		Hash:        convert.ToPtr(hash.Hash(b.ReleaseTitle, b.ReleaseDesc)),
		URL:         &b.ReleaseURL,
		ReleaseDate: processReleaseDate(b.ReleaseDate),
		Status:      convert.ToPtr(model.StatusPending),
	}
}

func processReleaseDate(s string) *time.Time {
	t, err := time.Parse(coinTelegraphTimeLayout, s)
	if err != nil {
		perr := err.(*time.ParseError)
		logrus.WithError(perr).Debug("Parsing failed...")

		tu, err := matchTimePassed(s)
		if err != nil {
			logrus.WithError(err).Debug("Match time parsing failed...")
			return nil
		}
		return convert.ToPtr(common.CurrentTimestamp().Add(-tu).Truncate(tu))
	}

	return &t
}

func convertTimeUnit(unit string) (time.Duration, error) {
	if strings.HasPrefix(unit, "HOUR") {
		return time.Hour, nil
	} else if strings.HasPrefix(unit, "MINUTE") {
		return time.Minute, nil
	} else if strings.HasPrefix(unit, "SECOND") {
		return time.Second, nil
	}
	return 0, errors.Errorf("Unknown time unit: %s", unit)
}

func matchTimePassed(releaseDate string) (time.Duration, error) {
	r := regexp.MustCompile(`(?P<num>\d+)\s*(?P<timeunit>HOUR(S)?|MINUTE(S)?|SECOND(S)?)\s*AGO`)

	matches := r.FindStringSubmatch(releaseDate)

	if len(matches) < 2 {
		return 0, errors.New("failed to match pattern")
	}
	num := matches[r.SubexpIndex("num")]
	timeUnit := matches[r.SubexpIndex("timeunit")]

	n, err := strconv.Atoi(num)
	if err != nil {
		return 0, errors.Wrap(err, "failed to convert time to int")
	}

	tu, err := convertTimeUnit(timeUnit)
	if err != nil {
		return 0, errors.Wrap(err, "failed to convert time unit")
	}
	return time.Duration(n) * tu, nil
}
