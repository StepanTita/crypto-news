package services

import (
	"regexp"
	"strings"

	"github.com/google/uuid"

	"common/data/model"
)

const maxInputChars = 50_000

var coinsRegex = regexp.MustCompile(`\<coins\>\[([A-Z1-9\,\s]+)\]\<\/coins\>`)

func parseCoins(content string) (string, []model.Coin) {
	coinsSet := make(map[string]bool)
	for _, match := range coinsRegex.FindAllStringSubmatch(content, -1) {
		for _, coin := range strings.Split(match[1], ",") {
			coinsSet[strings.TrimSpace(coin)] = true
		}
	}

	coins := make([]model.Coin, 0, len(coinsSet))

	for k := range coinsSet {
		coins = append(coins, model.Coin{
			Code: k,
			Slug: k,
		})
	}

	return coinsRegex.ReplaceAllString(content, ""), coins
}

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

func createCoinsNewsCoinsBatch(newsID uuid.UUID, coins []model.Coin) []model.NewsCoin {
	newsCoins := make([]model.NewsCoin, len(coins))
	for i, c := range coins {
		newsCoins[i] = model.NewsCoin{
			Code:   c.Code,
			NewsID: newsID,
		}
	}
	return newsCoins
}
