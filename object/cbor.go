package object

import (
	"errors"

	cbor "github.com/fxamacker/cbor/v2"
)

type cborOpaque[F DataTypes] struct{}

func Cbor[F DataTypes]() Transtyper[F, Opaque] {
	return &cborOpaque[F]{}
}

func (c *cborOpaque[F]) Process(o Object[F]) (Object[Opaque], error) {
	switch obj := o.(type) {
	case *object[F]:
		return c.encode(obj)
	default:
		return nil, errors.New("not supported object type")
	}
}

func (c *cborOpaque[F]) encode(obj *object[F]) (*object[Opaque], error) {
	cobj := &object[Opaque]{
		children: make(map[string]*object[Opaque], len(obj.children)),
		data:     make(map[string]Opaque, len(obj.data)),
	}

	for k, v := range obj.data {
		d, err := cbor.Marshal(v)
		if err != nil {
			return nil, err
		}
		cobj.data[k] = d
	}

	for k, v := range obj.children {
		cv, err := c.encode(v)
		if err != nil {
			return nil, err
		}

		cobj.children[k] = cv
	}

	return cobj, nil
}
