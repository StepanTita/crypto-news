package crawler

import (
	"common/data/model"
)

func ToModelBatch[T model.Model](bodies []ParsedBody) []T {
	newsBatch := make([]T, len(bodies))
	for i := range bodies {
		newsBatch[i] = bodies[i].ToModel().(T)
	}
	return newsBatch
}
