package crypto_panic_crawler

import (
	"time"

	"common/convert"
	"common/data/model"
	"parser/internal/utils"
)

type Body struct {
	Kind   *string `json:"kind"`
	Domain *string `json:"domain"`
	Source *struct {
		Title  *string `json:"title"`
		Region *string `json:"region"`
		Domain *string `json:"domain"`
		Path   *string `json:"path"`
		Url    *string `json:"url"`
	} `json:"source"`
	Title       *string    `json:"title"`
	PublishedAt *time.Time `json:"published_at"`
	Slug        *string    `json:"slug"`
	Currencies  []Currency `json:"currencies"`
	Id          *int       `json:"id"`
	Url         *string    `json:"url"`
	CreatedAt   *time.Time `json:"created_at"`
	Votes       *struct {
		Negative  *int `json:"negative"`
		Positive  *int `json:"positive"`
		Important *int `json:"important"`
		Liked     *int `json:"liked"`
		Disliked  *int `json:"disliked"`
		Lol       *int `json:"lol"`
		Toxic     *int `json:"toxic"`
		Saved     *int `json:"saved"`
		Comments  *int `json:"comments"`
	} `json:"votes"`
	Metadata *struct {
		Image       *string `json:"image"`
		Description *string `json:"description"`
	} `json:"metadata"`
}

// TODO: don't like the fact that code can be some trash data, need to handle this in the future
func toCoin(c Currency) model.Coin {
	return model.Coin{
		Code:  convert.FromPtr(c.Code),
		Title: convert.FromPtr(c.Title),
		Slug:  convert.FromPtr(c.Slug),
	}
}

func (b Body) ToNews() model.News {
	var imageUrl *string
	var description *string
	if b.Metadata != nil {
		imageUrl = b.Metadata.Image
		description = b.Metadata.Description
	}

	var url *string
	if b.Source != nil {
		url = b.Source.Url
	}

	return model.News{
		Media: &model.NewsMedia{
			Title: utils.StripHtmlRegex(b.Title),
			Text:  utils.StripHtmlRegex(description),
			Resources: []model.NewsMediaResource{
				{
					Type: convert.ToPtr(utils.ImageType),
					URL:  imageUrl,
				},
			},
		},
		PublishedAt:    b.PublishedAt,
		Url:            url,
		Source:         convert.ToPtr(utils.CryptoPanic),
		Status:         convert.ToPtr(model.StatusPending),
		OriginalSource: b.Source.Url,
		Coins:          utils.Map(b.Currencies, toCoin),
	}
}
