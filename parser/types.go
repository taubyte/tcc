package parser

import (
	"github.com/taubyte/go-seer"
	"github.com/taubyte/tcc/object"
)

var NodeDefaultSeerLeaf = "config"

type instance struct {
	schema *schemaDef
	seer   *seer.Seer
}

type schemaDef struct {
	root *Node
}

type Schema interface {
	Yaml() string
	Json() string
	Map() map[string]any
}

type Parser interface {
	Schema() Schema
	Parse() (object.Object[object.Refrence], error)
	Root() At // flows schema not yaml
}

type At interface {
	At(string) (At, error)
	Get(any) error // converts to schema expected type + validate type and value
	Set(any) error // validate type and value
}

type Type int

type SupportedTypes interface {
	int | bool | float64 | string | []string
}

const (
	TypeInt Type = iota
	TypeBool
	TypeFloat
	TypeString
	TypeStringSlice
)

type StringMatch any // string or PathMatcher

type AttributeValidator func(any) error

type Attribute struct {
	Name      string
	Type      Type
	Required  bool
	Key       bool // means the value is the key of a map
	Default   any
	Path      []StringMatch
	Compat    []StringMatch
	Validator AttributeValidator
}

type Option func(*Attribute)

type Node struct {
	Group      bool
	Match      StringMatch
	Attributes []*Attribute
	Children   []*Node
}

type at struct {
	node         *Node
	query        *seer.Query
	defaultValue any
}
