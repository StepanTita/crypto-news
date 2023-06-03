package model

import "fmt"

func PrependTableName(tableName string, columns []string) []string {
	for i := range columns {
		columns[i] = fmt.Sprintf("%s.%s", tableName, columns[i])
	}
	return columns
}
