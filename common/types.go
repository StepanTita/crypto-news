package common

import (
	"fmt"
	"reflect"
	"strings"

	"golang.org/x/exp/slices"
)

const (
	omitemptyTag = "omitempty"
	skipTag      = "-"
)

// StructTagsMap is a necessary workaround, since fields like uuid.UUID do not process IsZero correctly
// omits - are omitted always, omitempty are omitted only if IsZero -> true
func StructTagsMap(i any, skipNil bool) map[string]any {
	v := reflect.Indirect(reflect.ValueOf(i))
	if v.Kind() != reflect.Struct {
		panic(fmt.Sprintf("expected struct or pointer to struct, %v, received", v.Kind()))
	}
	t := v.Type()
	res := make(map[string]any)
	for i := 0; i < t.NumField(); i++ {
		keys := strings.Split(t.Field(i).Tag.Get("db"), ",")
		if slices.Contains(keys, skipTag) {
			continue
		}
		if skipNil && isNillable(v.Field(i)) && v.Field(i).IsNil() {
			continue
		}
		if skipNil && slices.Contains(keys, omitemptyTag) && v.Field(i).IsZero() {
			continue
		}
		res[keys[0]] = v.Field(i).Interface()
	}
	return res
}

func isNillable(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Chan, reflect.Map, reflect.Interface, reflect.Pointer, reflect.Func, reflect.Slice:
		return true
	default:
		return false
	}
}
