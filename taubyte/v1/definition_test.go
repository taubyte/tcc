package parser

import (
	"encoding/json"
	"fmt"
	"testing"

	_ "embed"

	"github.com/k0kubun/pp/v3"
	"github.com/taubyte/go-seer"
	"github.com/taubyte/tcc/parser"
	"gotest.tools/v3/assert"
)

//go:embed fixtures/schema.json
var parserJson string

func TestSchema(t *testing.T) {
	taubyteJson := TaubyteProject.Json()

	fmt.Println(taubyteJson)

	var taubyteData, parserData interface{}

	err := json.Unmarshal([]byte(taubyteJson), &taubyteData)
	if err != nil {
		t.Fatalf("Failed to unmarshal TaubyteProject's JSON: %v", err)
	}

	err = json.Unmarshal([]byte(parserJson), &parserData)
	if err != nil {
		t.Fatalf("Failed to unmarshal embedded parser JSON: %v", err)
	}

	assert.DeepEqual(t, taubyteData, parserData)
}

func TestConfigSchema(t *testing.T) {
	sr, err := seer.New(seer.SystemFS("fixtures/config"))
	assert.NilError(t, err)

	obj, err := parser.New(TaubyteProject).Parse(sr)
	assert.NilError(t, err)

	pp.Print(obj.Map())

	appObj, err := obj.Child("applications").Object()
	assert.NilError(t, err)

	app1Obj, err := appObj.Child("test_app1").Object()
	assert.NilError(t, err)

	app2Obj, err := appObj.Child("test_app2").Object()
	assert.NilError(t, err)

	app1funcsObj, err := app1Obj.Child("functions").Object()
	assert.NilError(t, err)

	app2funcsObj, err := app2Obj.Child("functions").Object()
	assert.NilError(t, err)

	funcsObj, err := obj.Child("functions").Object()
	assert.NilError(t, err)

	app1funcs2Obj, err := app1funcsObj.Child("test_function2").Object()
	assert.NilError(t, err)

	app2funcs2Obj, err := app2funcsObj.Child("test_function2").Object()
	assert.NilError(t, err)

	funcs2Obj, err := funcsObj.Child("test_function2").Object()
	assert.NilError(t, err)

	assert.DeepEqual(t, app1funcs2Obj.Map(), app2funcs2Obj.Map())
	assert.DeepEqual(t, funcs2Obj.Map(), app2funcs2Obj.Map())

}
