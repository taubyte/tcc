package object

import (
	"testing"

	cbor "github.com/fxamacker/cbor/v2"
	"gotest.tools/v3/assert"
)

func TestCborTranstyper(t *testing.T) {
	type sampleData struct {
		Name  string
		Value int
	}

	// Create a sample object
	origObj := New[sampleData]()
	data := sampleData{Name: "test", Value: 42}
	origObj.Set("key1", data)

	// Create child and set its data
	childSelector := origObj.Child("child")
	childData := sampleData{Name: "childTest", Value: 24}
	childSelector.Set("childKey", childData)

	// Use the cborOpaque transtyper to transform the object
	transtyper := Cbor[sampleData]()
	cborObj, err := transtyper.Process(origObj)
	assert.NilError(t, err)

	// Verify the CBOR encoding of the transformed object
	var decodedData sampleData
	err = cbor.Unmarshal(cborObj.Get("key1"), &decodedData)
	assert.NilError(t, err)
	assert.DeepEqual(t, decodedData, data)

	// Verify the CBOR encoding of the transformed child object
	cborChildSelector := cborObj.Child("child")
	var decodedChildData sampleData
	childObj, err := cborChildSelector.Object()
	assert.NilError(t, err)
	err = cbor.Unmarshal(childObj.Get("childKey"), &decodedChildData)
	assert.NilError(t, err)
	assert.DeepEqual(t, decodedChildData, childData)
}
