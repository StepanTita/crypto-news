package url_crawler

import (
	"common/data/model"
)

type body struct {
	text string
}

func (b body) ToModel() any {
	return model.RawNewsWebpage{
		Body: &b.text,
	}
}
