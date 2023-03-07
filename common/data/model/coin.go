package model

type Coin struct {
	Code  string `db:"code"`
	Title string `db:"title"`
	Slug  string `db:"slug"`
}

func (u Coin) TableName() string {
	return COINS
}
