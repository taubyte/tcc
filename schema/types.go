package parser

import (
	"github.com/taubyte/go-seer"
	"github.com/taubyte/tcc/object"
)

var NodeDefaultSeerLeaf = "config"

type schema struct {
	root *Node
}

type Schema interface {
	Yaml() string
	Json() string
	Map() map[string]any
	Parse(*seer.Seer) (object.Object[object.Refrence], error)
}

type Type int

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

type Engine interface {
}

type engine struct {
	_map Schema
}
