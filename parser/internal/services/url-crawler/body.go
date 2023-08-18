package url_crawler

import (
	"github.com/google/uuid"

	"common/data/model"
)

type body struct {
	titleID uuid.UUID
	text    string
}

func (b body) ToModel() any {
	return model.RawNews{
		TitleID: b.titleID,
		Body:    &b.text,
	}
}
