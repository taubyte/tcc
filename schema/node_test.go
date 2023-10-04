package parser

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestNodeToMap(t *testing.T) {
	node := &Node{
		Match: "test",
		Attributes: []*Attribute{
			{
				Name:    "attr1",
				Type:    TypeInt,
				Key:     true,
				Default: 123,
				Path:    []StringMatch{StringMatchAll{}},
				Compat:  []StringMatch{StringMatchAll{}},
			},
			{
				Name: "attr2",
				Type: TypeString,
				Path: []StringMatch{StringMatchAll{}},
			},
		},
		Children: []*Node{
			{
				Match: "child1",
			},
		},
	}

	expected := map[string]any{
		"match": "test",
		"group": false,
		"attributes": map[string]any{
			"attr1": map[string]any{
				"type":    "Int",
				"key":     true,
				"default": 123,
				"path":    stringify([]StringMatch{StringMatchAll{}}),
				"compat":  stringify([]StringMatch{StringMatchAll{}}),
			},
			"attr2": map[string]any{
				"type": "String",
				"path": stringify(StringMatchAll{}),
			},
		},
		"children": []any{
			map[string]any{
				"match":      "child1",
				"group":      false,
				"attributes": map[string]any{},
				"children":   []any{},
			},
		},
	}

	assert.DeepEqual(t, node.Map(), expected)
}
