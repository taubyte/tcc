package parser

import (
	"encoding/json"

	"github.com/taubyte/go-seer"
	"github.com/taubyte/tcc/object"
	"gopkg.in/yaml.v2"
)

func (s *schema) Yaml() string {
	out, err := yaml.Marshal(s.Map())
	if err != nil {
		panic(err)
	}
	return string(out)
}

func (s *schema) Json() string {
	out, err := json.Marshal(s.Map())
	if err != nil {
		panic(err)
	}
	return string(out)
}

func (s *schema) Map() map[string]any {
	return map[string]any{
		"root": s.root.Map(),
	}
}

func (s *schema) Parse(sr *seer.Seer) (object.Object[object.Refrence], error) {
	return s.root.Object(sr.Query())
}

func SchemaDefinition(root *Node) Schema {
	return &schema{
		root: root,
	}
}
