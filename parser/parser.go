package parser

import (
	"github.com/taubyte/go-seer"
	"github.com/taubyte/tcc/object"
)

func New(schema Schema) Parser {
	return &instance{
		schema: schema.(*schemaDef),
	}
}

func (s *instance) Parse(sr *seer.Seer) (object.Object[object.Refrence], error) {
	return s.schema.root.Object(sr.Query())
}

func (s *instance) Schema() Schema {
	return s.schema
}
