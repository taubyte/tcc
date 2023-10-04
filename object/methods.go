package object

import (
	"errors"
	"os"
	"regexp"
	"strings"
)

type selector[T DataTypes] struct {
	parent *object[T]
	obj    *object[T]
	name   string
}

type object[T DataTypes] struct {
	children map[string]*object[T]
	data     map[string]T
}

func New[T DataTypes]() Object[T] {
	return &object[T]{
		children: make(map[string]*object[T]),
		data:     make(map[string]T),
	}
}

func (o *object[T]) Map() map[string]any {
	m := map[string]any{
		"attributes": o.data,
	}
	for n, o := range o.children {
		m[n] = o.Map()
	}
	return m
}

func (o *object[T]) Children() []string {
	ret := make([]string, 0, len(o.children))
	for k := range o.children {
		ret = append(ret, k)
	}
	return ret
}

func (o *object[T]) Child(sel any) Selector[T] {
	switch _sel := sel.(type) {
	case string:
		return o.newNameSelector(_sel)
	case *object[T]: // save interface conv
		return o.newPtrSelector(_sel)
	case Object[T]:
		return o.newPtrSelector(_sel.(*object[T]))
	default:
		panic("unknown object type")
	}
}

func (o *object[T]) Set(name string, data T) {
	o.data[name] = data
}

func (o *object[T]) Get(name string) T {
	return o.data[name]
}

func (o *object[T]) getObjectName(obj *object[T]) string {
	for name, _obj := range o.children {
		if obj == _obj {
			return name
		}
	}
	return ""
}

func (o *object[T]) getObjectByName(name string) (*object[T], error) {
	if obj, ok := o.children[name]; ok {
		return obj, nil
	}
	return nil, os.ErrNotExist
}

func (o *object[T]) exists(name string, obj *object[T]) bool {
	if _, ok := o.children[name]; ok {
		return true
	}

	for _, _obj := range o.children {
		if obj == _obj {
			return true
		}
	}

	return false
}

func (o *object[T]) set(name string, obj *object[T]) error {
	if obj == nil {
		return errors.New("nil object")
	}
	o.children[name] = obj
	return nil
}

func exactMatch(a, b string) bool {
	return a == b
}

func (o *object[T]) Match(expr string, mtype MatchType) ([]Object[T], error) {
	var cmp func(string, string) bool
	switch mtype {
	case ExactMatch:
		cmp = exactMatch
	case PrefixMatch:
		cmp = strings.HasPrefix
	case SuffixMatch:
		cmp = strings.HasSuffix
	case SubMatch:
		cmp = strings.Contains
	case RegExMatch:
		// compile once
		cexpr, err := regexp.Compile(expr)
		if err != nil {
			return nil, err
		}
		cmp = func(str, _ string) bool {
			return cexpr.MatchString(str)
		}
	default:
		return nil, errors.New("unknown match type")
	}

	matches := make([]Object[T], 0)
	for name, _obj := range o.children {
		if cmp(name, expr) {
			matches = append(matches, _obj)
		}
	}

	return matches, nil
}

func (o *object[T]) newNameSelector(name string) *selector[T] {
	return &selector[T]{parent: o, name: name}
}

func (o *object[T]) newPtrSelector(obj *object[T]) *selector[T] {
	return &selector[T]{parent: o, obj: obj}
}

func (s *selector[T]) Name() string {
	if len(s.name) == 0 {
		s.name = s.parent.getObjectName(s.obj)
	}
	return s.name
}

func (s *selector[T]) Exists() bool {
	return s.parent.exists(s.name, s.obj)
}

func (s *selector[T]) Object() (Object[T], error) {
	var err error
	if s.obj == nil {
		s.obj, err = s.parent.getObjectByName(s.name)
	}
	return s.obj, err
}

func (s *selector[T]) Set(name string, data T) error {
	if s.obj == nil {
		s.obj = New[T]().(*object[T])
	}
	s.obj.Set(name, data)
	return s.parent.set(s.name, s.obj)
}

func (s *selector[T]) Add(o Object[T]) error {
	var ok bool
	s.obj, ok = o.(*object[T])
	if !ok {
		return errors.New("unkown object type")
	}
	return s.parent.set(s.name, s.obj)
}

func (s *selector[T]) Get(name string) (ret T, err error) {
	if s.obj == nil {
		if _, err = s.Object(); err != nil {
			return
		}
	}

	ret = s.obj.Get(name)
	return
}
