package parser

import (
	"errors"
	"fmt"

	"github.com/taubyte/go-seer"
)

// TODO: test of not a group should not have children
// or replace Group bu len(children)==0

func buildPathQuery(path []StringMatch, node *Node, query *seer.Query) (*seer.Query, error) {
	var err error

	query = query.Fork()

	for _, itm := range path {
		switch pitm := itm.(type) {
		case string:
			query.Get(pitm)
			node, err = node.ChildMatch(pitm)
			if err != nil {
				return nil, err
			}
		case StringMatcher:
			list, err := query.Fork().List()
			if err != nil {
				return nil, fmt.Errorf("list path matches failed with %w", err)
			}

			var found bool
			for _, l := range list {
				if pitm.Match(l) {
					found = true
					query.Get(l)
					break
				}
			}

			if !found {
				return nil, errors.New("can't find match for path")
			}
		}
	}

	return query, nil
}

func (a *at) attribute(name string) (*seer.Query, *Attribute, error) {
	if len(a.node.Attributes) == 0 {
		return nil, nil, errors.New("no attributes")
	}

	for _, attr := range a.node.Attributes {
		if attr.Name != name {
			continue
		}

		if len(attr.Path) == 0 {
			attr.Path = []StringMatch{attr.Name}
		}
		aq, err := buildPathQuery(attr.Path, a.query)
		if err != nil && len(attr.Compat) > 0 {
			aq, err = buildPathQuery(attr.Compat, a.query)
		}

		if err != nil {
			return nil, attr, errors.New("attribute not set")
		}

		return aq, attr, nil
	}

	return nil, nil, errors.New("attribute no found")
}

func (a *at) At(name string) (At, error) {
	q, attr, err := a.attribute(name)
	if err != nil {
		if attr.Default != nil {
			return &at{defaultValue: attr.Default}
		}
		return nil, errors.New("not set")
	}

	if d != nil {
		return &at{
			node: nil,
		}
	}

	if q != nil {

	}

	if !a.node.Group {
		return nil, errors.New("not a group")
	}

	// TODO need to get attributes

	for _, itm := range a.node.Children {
		var matched bool

		switch i := itm.Match.(type) {
		case string:
			if i == name {
				matched = true
			}
		case StringMatcher:
			if i.Match(name) {
				matched = true
			}
		}

		if matched {
			return &at{
				node:  itm,
				query: a.query.Fork().Get(name),
			}, nil
		}
	}

	return nil, fmt.Errorf("no match found for `%s`", name)
}

func (a *at) Get(v any) error { // converts to schema expected type + validate type and value
	// validate
	return a.query.Fork().Value(v)
}

func (a *at) Set(v any) error { // validate type and value
	//validate
	return a.query.Fork().Set(v).Commit()
}

// helpers
func Get[T SupportedTypes](a At) (v T, err error) {
	err = a.Get(&v)
	return
}

func Set[T SupportedTypes](a At, v T) error {
	return a.Set(v)
}
