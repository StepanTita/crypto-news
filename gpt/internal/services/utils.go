package services

import (
	"regexp"

	"github.com/google/uuid"

	"common/data/model"
)

const keyPrevDigest = "news/prev-digest/<locale:/%s>"

const maxInputChars = 4000

var coinsRegex = regexp.MustCompile(`\<coins\>\[([A-Z1-9\,\s]+)\]\<\/coins\>`)

func toNewsChannelsBatch(news *model.News, channels []model.Channel) []model.NewsChannel {
	newsChannels := make([]model.NewsChannel, len(channels))
	for i, c := range channels {
		newsChannels[i] = model.NewsChannel{
			ChannelID: c.ChannelID,
			NewsID:    news.ID,
		}
	}
	return newsChannels
}

func createCoinsNewsCoinsBatch(newsID uuid.UUID, codes []string) ([]model.Coin, []model.NewsCoin) {
	coins := make([]model.Coin, len(codes))
	newsCoins := make([]model.NewsCoin, len(codes))
	for i, code := range codes {
		newsCoins[i] = model.NewsCoin{
			Code:   code,
			NewsID: newsID,
		}

		coins[i] = model.Coin{
			Code: code,
			Slug: code,
		}
	}
	return coins, newsCoins
}
