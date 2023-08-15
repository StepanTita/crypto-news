package deep_map

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type DeepMap struct {
	node *yaml.Node
}

// NewDeepMap creates a new instance of DeepMap from raw YAML bytes.
func NewDeepMap(node yaml.Node) *DeepMap {
	return &DeepMap{
		node: &node,
	}
}

// Get returns a child DeepMap for a given key.
func (y *DeepMap) Get(key string) (*DeepMap, error) {
	if y.IsMap() {
		for i := 0; i < len(y.node.Content); i += 2 {
			k := y.node.Content[i]
			v := y.node.Content[i+1]
			if k.Value == key {
				return &DeepMap{node: v}, nil
			}
		}
		return nil, errors.New("key not found")
	}
	return nil, errors.New("current node is not a YAML map")
}

// IsMap checks if the current node is a YAML map.
func (y *DeepMap) IsMap() bool {
	return y.node.Kind == yaml.MappingNode
}

// IsScalar checks if the current node is a scalar value.
func (y *DeepMap) IsScalar() bool {
	return y.node.Kind == yaml.ScalarNode
}

// IsList checks if the current node is a YAML list.
func (y *DeepMap) IsList() bool {
	return y.node.Kind == yaml.SequenceNode
}

func (y *DeepMap) GetScalar() string {
	return y.node.Value
}
