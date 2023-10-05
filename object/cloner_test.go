package object

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestCloner(t *testing.T) {
	// Setup a sample object
	origObj := New[Opaque]()
	origObj.Set("key1", Opaque("value1"))
	child := New[Opaque]()
	child.Set("childKey", Opaque("childValue"))
	origObj.Child("child").Add(child)

	// Clone the object
	transformer := Cloner()
	clonedObj, err := transformer.Process(origObj)
	assert.NilError(t, err)

	clonedChild, err := clonedObj.Child("child").Object()
	assert.NilError(t, err)

	// Check that original and clone have the same values
	assert.DeepEqual(t, origObj.Get("key1"), clonedObj.Get("key1"))
	assert.DeepEqual(t, child.Get("childKey"), clonedChild.Get("childKey"))

	// Modify the clone and verify that the original is unaffected
	clonedObj.Set("key1", Opaque("modifiedValue"))
	clonedChild.Set("childKey", Opaque("modifiedChildValue"))

	assert.DeepEqual(t, origObj.Get("key1"), Opaque("value1"))
	assert.DeepEqual(t, child.Get("childKey"), Opaque("childValue"))

	// Check that modifications were correctly applied to the clone
	assert.DeepEqual(t, clonedObj.Get("key1"), Opaque("modifiedValue"))
	assert.DeepEqual(t, clonedChild.Get("childKey"), Opaque("modifiedChildValue"))
}
