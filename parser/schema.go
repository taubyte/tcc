package parser

import (
	"encoding/json"

	"gopkg.in/yaml.v2"
)

func (s *schemaDef) Yaml() string {
	out, err := yaml.Marshal(s.Map())
	if err != nil {
		panic(err)
	}
	return string(out)
}

func (s *schemaDef) Json() string {
	out, err := json.Marshal(s.Map())
	if err != nil {
		panic(err)
	}
	return string(out)
}

func (s *schemaDef) Map() map[string]any {
	return map[string]any{
		"root": s.root.Map(),
	}
}

func SchemaDefinition(root *Node) Schema {
	return &schemaDef{
		root: root,
	}
}
