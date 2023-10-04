package parser

import (
	"errors"
	"fmt"

	"github.com/taubyte/go-seer"
	"github.com/taubyte/tcc/object"
)

func (n *Node) Map() map[string]any {
	return map[string]any{
		"group":      n.Group,
		"match":      stringify(n.Match),
		"attributes": n.attributesToMap(),
		"children":   n.childrenToSlice(),
	}
}

func (n *Node) attributesToMap() map[string]any {
	ret := make(map[string]any, len(n.Attributes))
	for _, attr := range n.Attributes {
		m := map[string]any{
			"type": attr.Type.String(),
		}
		if attr.Key {
			m["key"] = true
		}
		if attr.Default != nil {
			m["default"] = attr.Default
		}
		if len(attr.Path) > 0 {
			m["path"] = stringify(attr.Path)
		}
		if len(attr.Compat) > 0 {
			m["compat"] = stringify(attr.Compat)
		}
		ret[attr.Name] = m
	}
	return ret
}

func (n *Node) childrenToSlice() []any {
	ret := make([]any, len(n.Children))
	for i, node := range n.Children {
		ret[i] = node.Map()
	}
	return ret
}

func buildPathQuery(path []StringMatch, query *seer.Query) (*seer.Query, error) {
	query = query.Fork()

	for _, itm := range path {
		switch pitm := itm.(type) {
		case string:
			query.Get(pitm)
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
				}
			}
			if !found {
				return nil, errors.New("can't find match for path")
			}
		}
	}

	return query, nil
}

func (n *Node) hasRequiredAttributes() bool {
	for _, attr := range n.Attributes {
		if attr.Required {
			return true
		}
	}
	return false
}

func (n *Node) setAttributes(obj object.Object[object.Refrence], query *seer.Query) error {
	if len(n.Attributes) == 0 {
		return nil
	}
	for _, attr := range n.Attributes {
		if len(attr.Path) == 0 {
			attr.Path = []StringMatch{attr.Name}
		}
		aq, err := buildPathQuery(attr.Path, query)
		if err != nil && len(attr.Compat) > 0 {
			aq, err = buildPathQuery(attr.Compat, query)
		}

		if err != nil {
			if attr.Default != nil {
				obj.Set(attr.Name, attr.Default)
			}
			continue // no attribute is required
		}

		var val any
		switch attr.Type {
		case TypeInt:
			var v int
			err = aq.Value(&v)
			val = v
		case TypeBool:
			var v bool
			err = aq.Value(&v)
			val = v
		case TypeFloat:
			var v float64
			err = aq.Value(&v)
			val = v
		case TypeString:
			var v string
			err = aq.Value(&v)
			val = v
		case TypeStringSlice:
			var v []string
			err = aq.Value(&v)
			val = v
		}

		if err == nil && attr.Validator != nil {
			err = attr.Validator(val)
		}

		if err != nil {
			if attr.Default != nil {
				val = attr.Default
			} else if attr.Required {
				return err
			} else {
				continue
			}
		}

		obj.Set(attr.Name, val)
	}

	return nil
}

func (n *Node) Object(query *seer.Query) (object.Object[object.Refrence], error) {
	obj := object.New[object.Refrence]()

	if !n.Group {
		if err := n.setAttributes(obj, query); err != nil {
			return nil, err
		}
		return obj, nil
	}

	// file might or might not have config, so we ignore error
	err := n.setAttributes(obj, query.Fork().Get(NodeDefaultSeerLeaf))
	if err != nil && n.hasRequiredAttributes() {
		return nil, err
	}

	list, _ := query.Fork().List()
	for _, itm := range n.Children {
		for _, l := range list {
			if l == NodeDefaultSeerLeaf {
				continue
			}

			var (
				match   string
				matched bool
			)

			switch i := itm.Match.(type) {
			case string:
				if i == l {
					match = l
					matched = true
				}
			case StringMatcher:
				if i.Match(l) {
					match = l
					matched = true
				}
			}

			if !matched {
				continue
			}

			cobj, err := itm.Object(query.Fork().Get(match))
			if err != nil {
				return nil, err
			}

			err = obj.Child(match).Add(cobj)
			if err != nil {
				return nil, err
			}
		}

	}

	return obj, nil
}
