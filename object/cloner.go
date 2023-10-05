package object

import "errors"

type cloner struct{}

func Cloner() Transformer[Opaque] {
	return &cloner{}
}

func (c *cloner) Process(o Object[Opaque]) (Object[Opaque], error) {
	switch obj := o.(type) {
	case *object[Opaque]:
		return c.clone(obj)
	default:
		return nil, errors.New("not supported object type")
	}
}

func (c *cloner) clone(obj *object[Opaque]) (*object[Opaque], error) {
	cobj := &object[Opaque]{
		children: make(map[string]*object[Opaque], len(obj.children)),
		data:     make(map[string]Opaque, len(obj.data)),
	}

	for k, v := range obj.data {
		cobj.data[k] = c.copy(v)
	}

	for k, v := range obj.children {
		cv, err := c.clone(v)
		if err != nil {
			return nil, err
		}

		cobj.children[k] = cv
	}

	return cobj, nil
}

func (c *cloner) copy(d Opaque) Opaque {
	b := make(Opaque, len(d))
	copy(b, d)
	return b
}
