package parser

import (
	"fmt"

	"github.com/taubyte/go-seer"
	"github.com/taubyte/tcc/object"
)

func New(schema Schema, options ...seer.Option) (Parser, error) {
	sr, err := seer.New(options...)
	if err != nil {
		return nil, fmt.Errorf("parser failed to created seer with %w", err)
	}

	return &instance{
		schema: schema.(*schemaDef),
		seer:   sr,
	}, nil
}

func (s *instance) Parse() (object.Object[object.Refrence], error) {
	return s.schema.root.Object(s.seer.Query())
}

func (s *instance) Schema() Schema {
	return s.schema
}

func (s *instance) Root() At {
	return &at{node: s.schema.root}
}
