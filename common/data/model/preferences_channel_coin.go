package model

type PreferencesChannelCoin struct {
	ChannelID int64  `db:"channel_id"`
	CoinCode  string `db:"coin_code"`
}

func (n PreferencesChannelCoin) TableName() string {
	return PREFERENCES_CHANNEL_COINS
}
