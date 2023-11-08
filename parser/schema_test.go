package parser

import (
	"encoding/json"
	"testing"

	"github.com/taubyte/go-seer"
	"gopkg.in/yaml.v2"
	"gotest.tools/v3/assert"
)

func setupTestSchema() *schemaDef {
	return &schemaDef{
		root: &Node{
			Match: "sample",
			Attributes: []*Attribute{
				{Name: "sampleAttr", Type: TypeString},
			},
		},
	}
}

func TestSchemaToMap(t *testing.T) {
	parser := setupTestSchema()

	expectedMap := map[string]any{
		"root": map[string]any{
			"match": "sample",
			"group": false,
			"attributes": map[string]any{
				"sampleAttr": map[string]any{
					"type": "String",
				},
			},
			"children": []any{},
		},
	}
	assert.DeepEqual(t, parser.Map(), expectedMap)
}

func TestSchemaJson(t *testing.T) {
	parser := setupTestSchema()

	expectedMap := map[string]any{
		"root": map[string]any{
			"match": "sample",
			"group": false,
			"attributes": map[string]any{
				"sampleAttr": map[string]any{
					"type": "String",
				},
			},
			"children": []any{},
		},
	}
	expectedJson, err := json.Marshal(expectedMap)
	assert.NilError(t, err)
	assert.Equal(t, parser.Json(), string(expectedJson))
}

func TestSchemaYaml(t *testing.T) {
	parser := setupTestSchema()

	expectedMap := map[string]any{
		"root": map[string]any{
			"match": "sample",
			"group": false,
			"attributes": map[string]any{
				"sampleAttr": map[string]any{
					"type": "String",
				},
			},
			"children": []any{},
		},
	}
	expectedYaml, err := yaml.Marshal(expectedMap)
	assert.NilError(t, err)
	assert.Equal(t, parser.Yaml(), string(expectedYaml))
}

var resSchemaDef = []*Node{
	DefineGroup("type1",
		DefineIter(
			Attributes(
				String("name"),
				Bool("really", Path("question", "really"), Required()),
				Int("count", Default(1)),
				Int("size", Default(10)),
			),
		),
	),
	DefineGroup("type2",
		DefineIter(
			Attributes(
				String("fqdn", IsFqdn()),
				String("type", Path("object", "type"), InSet("oval", "rect"), Default("x509")),
				String("type", Path("object", "size"), InSet(0, 16, 32)),
			),
		),
	),
}

var testSchemaDef = SchemaDefinition(
	Root(
		Attributes(
			String("email", Path("notification", "email"), IsEmail(), Required()),
		),
		append(
			resSchemaDef,
			DefineGroup("sub",
				resSchemaDef...,
			),
		)...,
	),
)

func TestSchemaParse(t *testing.T) {
	p, err := New(testSchemaDef, seer.SystemFS("fixtures/config"))
	assert.NilError(t, err)

	obj, err := p.Parse()
	assert.NilError(t, err)

	// Top-level attributes
	assert.Equal(t, obj.Get("email"), "yo@yo.com")

	// type1 children and its attributes
	objType1, err := obj.Child("type1").Object()
	assert.NilError(t, err)
	it1, err := objType1.Child("it1").Object()
	assert.NilError(t, err)

	// Check attributes of 'it1' under 'type1'
	assert.Equal(t, it1.Get("name"), "it1")
	assert.Equal(t, it1.Get("count"), 1)
	assert.Equal(t, it1.Get("really"), true)
	assert.Equal(t, it1.Get("size"), 10)

	// sub children and its attributes
	subObj, err := obj.Child("sub").Object()
	assert.NilError(t, err)

	// type2 children and its attributes
	_, err = subObj.Child("type2").Object()
	assert.NilError(t, err)

	type1Obj, err := obj.Child("type1").Object()
	assert.NilError(t, err)

	t1_2, err := type1Obj.Child("it1").Object()
	assert.NilError(t, err)

	assert.Equal(t, t1_2.Get("name"), it1.Get("name"))
}

func TestBadSchemaParse(t *testing.T) {
	for _, i := range []string{"1", "2", "3"} {
		p, err := New(testSchemaDef, seer.SystemFS("fixtures/config_bad_"+i))
		assert.NilError(t, err)
		_, err = p.Parse()
		if err == nil {
			t.Error("should have failed")
			t.FailNow()
		}
	}
}
