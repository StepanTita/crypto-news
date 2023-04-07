package model

import (
	"fmt"
	"strings"

	"github.com/google/uuid"

	"common"
)

type Model interface {
	News | Coin | Channel | NewsCoin | NewsChannel | PreferencesChannelCoin | AuthorizationKeys
	TableName() string
}

// ToMap used to translate struct to map, where keys - are db tags values, and values - corresponding fields' values
// it omits passed fields, for the sake of valid insertion and generation, it also omits Nil fields,
// if some value is intended to be set Nil, use sql.Null[Type]
func ToMap[T Model](v T) map[string]any {
	return common.StructTagsMap(v, true)
}

// skipNil should be false by default here
func Columns[T Model](v T, skipNil bool) []string {
	columnsMap := common.StructTagsMap(v, skipNil)
	columns := make([]string, 0, len(columnsMap))
	for k := range columnsMap {
		columns = append(columns, k)
	}
	return columns
}

func NamedBinding[T Model](v T) ([]string, []string) {
	columns := Columns(v, true)
	namedColumns := make([]string, len(columns))

	for i, v := range columns {
		namedColumns[i] = fmt.Sprintf(":%s", v)
	}
	return columns, namedColumns
}

func ToKey[T Model](i T, unique bool) string {
	basicTypeKey := strings.ReplaceAll(common.GetTypeName(i), ".", "/")
	if unique {
		return fmt.Sprintf("%s/%s", strings.ReplaceAll(basicTypeKey, ".", "/"), uuid.NewString())
	}
	return basicTypeKey
}
