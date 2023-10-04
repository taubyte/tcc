package object

import (
	"bytes"
	"os"
	"testing"

	"gotest.tools/v3/assert"
)

func TestNew(t *testing.T) {
	o := New[Opaque]()
	if _, ok := o.(*object[Opaque]); !ok {
		t.Errorf("Expected type *object[Opaque], got %T", o)
	}
}

func TestObjectSetGet(t *testing.T) {
	o := New[Opaque]()
	data := Opaque{1, 2, 3}
	o.Set("attr", data)
	if !bytes.Equal(o.Get("attr"), data) {
		t.Errorf("Expected data %v, got %v", data, o.Get("attr"))
	}
}

func TestMatch(t *testing.T) {
	o := New[Opaque]()
	o.Child("apple").Set("attr", Opaque{})
	o.Child("applepie").Set("attr", Opaque{})
	o.Child("pieapple").Set("attr", Opaque{})

	tests := []struct {
		expr  string
		mtype MatchType
		want  int
	}{
		{"apple", ExactMatch, 1},
		{"apple", PrefixMatch, 2},
		{"apple", SuffixMatch, 2},
		{"apple", SubMatch, 3},
		{"app.*pie", RegExMatch, 1},
	}

	for _, test := range tests {
		got, err := o.Match(test.expr, test.mtype)
		if err != nil {
			t.Errorf("Error on matching %s: %s", test.expr, err)
			continue
		}
		if len(got) != test.want {
			t.Errorf("For %s and type %v expected %d matches, got %d", test.expr, test.mtype, test.want, len(got))
		}
	}
}

func TestSelector(t *testing.T) {
	o := New[Opaque]()
	name := "testName"
	data := Opaque{4, 5, 6}
	selector := o.Child(name)
	if err := selector.Set("attr", data); err != nil {
		t.Errorf("Failed to set data on selector: %s", err)
	}
	if selector.Name() != name {
		t.Errorf("Expected name %s, got %s", name, selector.Name())
	}
	if !selector.Exists() {
		t.Error("Expected the selector to exist")
	}
	obj, err := selector.Object()
	if err != nil {
		t.Errorf("Failed to get object from selector: %s", err)
	}
	if !bytes.Equal(obj.Get("attr"), data) {
		t.Errorf("Expected data %v, got %v", data, obj.Get("attr"))
	}
	getData, err := selector.Get("attr")
	if err != nil {
		t.Errorf("Failed to get data from selector: %s", err)
	}
	if !bytes.Equal(getData, data) {
		t.Errorf("Expected data %v, got %v", data, getData)
	}
}

func TestSelectorFailures(t *testing.T) {
	o := New[Opaque]()
	selector := o.Child("nonExistent")
	if selector.Exists() {
		t.Error("Expected the selector to not exist")
	}
	if _, err := selector.Object(); err == nil {
		t.Error("Expected an error when getting object from non-existent selector")
	}
	if _, err := selector.Get("attr"); err == nil {
		t.Error("Expected an error when getting data from non-existent selector")
	}
}

func TestChildPanic(t *testing.T) {
	o := New[Opaque]()
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic when providing unknown object type")
		}
	}()
	o.Child(42) // This should panic
}

func TestGetObjectByName(t *testing.T) {
	o := New[Opaque]()
	childName := "child"
	child := New[Opaque]()
	o.Child(childName).Set("attr", child.Get("attr"))
	_, err := o.(*object[Opaque]).getObjectByName(childName)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	_, err = o.(*object[Opaque]).getObjectByName("nonexistent")
	if err == nil {
		t.Error("Expected an error when getting object by nonexistent name")
	}
}

func TestObjectExists(t *testing.T) {
	o := New[Opaque]()
	childName := "child"
	child := New[Opaque]()
	o.Child(childName).Set("attr", child.Get("attr"))

	if !o.(*object[Opaque]).exists(childName, nil) {
		t.Errorf("Expected child %s to exist", childName)
	}

	if o.(*object[Opaque]).exists("nonexistent", nil) {
		t.Error("Expected child nonexistent to not exist")
	}

	otherChild := New[Opaque]()
	if o.(*object[Opaque]).exists("", otherChild.(*object[Opaque])) {
		t.Error("Expected otherChild to not exist")
	}
}

func TestSelectorExistingObject(t *testing.T) {
	o := New[Opaque]()
	childName := "child"
	childData := Opaque{7, 8, 9}
	child := o.Child(childName)
	child.Set("attr", childData)
	selectorForExisting := o.Child(childName)
	if selectorForExisting.Name() != childName {
		t.Errorf("Expected name %s, got %s", childName, selectorForExisting.Name())
	}
	if !selectorForExisting.Exists() {
		t.Error("Expected the selector for existing object to exist")
	}
	retrievedData, err := selectorForExisting.Get("attr")
	if err != nil {
		t.Errorf("Failed to get data from selector: %s", err)
	}
	if !bytes.Equal(retrievedData, childData) {
		t.Errorf("Expected data %v, got %v", childData, retrievedData)
	}
}

func TestObjectSetName(t *testing.T) {
	o := New[Opaque]()
	childName := "child"
	childData := Opaque{10, 11, 12}
	child := o.Child(childName)
	child.Set("attr", childData)
	newChildName := "renamedChild"
	o.Child(newChildName).Set("attr", childData)
	if o.(*object[Opaque]).exists(newChildName, nil) {
		ncdata, err := o.Child(newChildName).Get("attr")
		if err != nil {
			t.Error(err)
		}
		if !bytes.Equal(ncdata, childData) {
			t.Errorf("Expected data %v, got %v", childData, ncdata)
		}
	} else {
		t.Errorf("Child with name %s doesn't exist", newChildName)
	}
}

func TestChildWithPointer(t *testing.T) {
	o := New[Opaque]()
	childName := "child"
	child := o.Child(childName)
	child.Set("attr", Opaque{7, 8, 9})
	selectorFromPointer := o.Child(child.(*selector[Opaque]).obj)
	if selectorFromPointer.Name() != childName {
		t.Errorf("Expected name %s, got %s", childName, selectorFromPointer.Name())
	}
}

func TestSetExistingChild(t *testing.T) {
	o := New[Opaque]()
	childName := "child"
	childData1 := Opaque{7, 8, 9}
	childData2 := Opaque{10, 11, 12}
	o.Child(childName).Set("attr", childData1)
	o.Child(childName).Set("attr", childData2) // Overwriting
	retrievedData, _ := o.Child(childName).Get("attr")
	if !bytes.Equal(retrievedData, childData2) {
		t.Errorf("Expected data %v, got %v", childData2, retrievedData)
	}
}

func TestMatchUnknownType(t *testing.T) {
	o := New[Opaque]()
	_, err := o.Match("apple", MatchType(999)) // Unknown MatchType
	if err == nil {
		t.Error("Expected an error for unknown match type")
	}
}

func TestSetNameForExistingObject(t *testing.T) {
	o := New[Opaque]()
	childName := "child"
	childData := Opaque{7, 8, 9}
	o.Child(childName).Set("attr", childData)
	newChild := o.Child(childName).(*selector[Opaque]).obj
	o.Child(newChild).Set("attr", childData)
	if !o.(*object[Opaque]).exists(childName, newChild) {
		t.Errorf("Failed to set child with existing object")
	}
}

func TestExistsWithObjectReference(t *testing.T) {
	o := New[Opaque]()
	childName := "child"
	childData := Opaque{10, 11, 12}
	child := o.Child(childName)
	child.Set("attr", childData)
	if !o.(*object[Opaque]).exists("", child.(*selector[Opaque]).obj) {
		t.Errorf("Expected child %s to exist by object reference", childName)
	}
}

func TestChildPanicOnUnknownType(t *testing.T) {
	o := New[Opaque]()
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected Child() to panic for unknown object type")
		}
	}()
	o.Child(123) // This should panic
}

func TestGetObjectByNameNonExistent(t *testing.T) {
	o := New[Opaque]()
	_, err := o.(*object[Opaque]).getObjectByName("nonExistentName")
	if err == nil || err != os.ErrNotExist {
		t.Errorf("Expected os.ErrNotExist for non-existent name")
	}
}

func TestExistsByName(t *testing.T) {
	o := New[Opaque]()
	childName := "child"
	if o.(*object[Opaque]).exists(childName, nil) {
		t.Errorf("Expected child %s to not exist", childName)
	}
}

func TestObjectSetGetByName(t *testing.T) {
	o := New[Opaque]()

	name1 := "testName1"
	data1 := Opaque{1, 2, 3}
	o.Set(name1, data1)
	if got := o.Get(name1); !bytes.Equal(got, data1) {
		t.Errorf("For name %s, expected %v, got %v", name1, data1, got)
	}

	name2 := "testName2"
	data2 := Opaque("refValue")
	o.Set(name2, data2)
	if got := o.Get(name2); !bytes.Equal(got, data2) {
		t.Errorf("For name %s, expected %v, got %v", name2, data2, got)
	}
}

func TestObjectReference(t *testing.T) {
	o := New[Refrence]()

	refName := "refName"
	data := Refrence("Hello World")
	o.Set(refName, data)

	got := o.Get(refName)
	if got != data {
		t.Errorf("Expected reference %v, got %v", data, got)
	}
}

func TestSelectorForReference(t *testing.T) {
	o := New[Refrence]()
	name := "refName"
	data := Refrence("Hello Again")
	selector := o.Child(name)
	if err := selector.Set(name, data); err != nil {
		t.Errorf("Failed to set data on selector: %s", err)
	}
	if selector.Name() != name {
		t.Errorf("Expected name %s, got %s", name, selector.Name())
	}
	if !selector.Exists() {
		t.Error("Expected the selector to exist")
	}
	getData, err := selector.Get(name)
	if err != nil {
		t.Errorf("Failed to get data from selector: %s", err)
	}
	if getData != data {
		t.Errorf("Expected data %v, got %v", data, getData)
	}
}

func TestSelectorAdd(t *testing.T) {
	o := New[Opaque]()
	name := "newChild"
	selector := o.Child(name)

	// Confirm the child doesn't exist yet
	if selector.Exists() {
		t.Error("Child shouldn't exist yet")
	}

	// Add the child using the selector's Add() method
	err := selector.Add(o)
	assert.NilError(t, err, "Expected no error when adding a child")

	// Check if the child now exists
	if !selector.Exists() {
		t.Error("Child should exist after calling Add()")
	}

	// Fetch the object from the selector to ensure it was correctly added
	_, err = selector.Object()
	assert.NilError(t, err, "Expected no error when fetching the object after adding it")
}

func TestObjectChildren(t *testing.T) {
	o := New[Opaque]()
	childNames := []string{"child1", "child2", "child3"}

	// Add children to the object
	for _, name := range childNames {
		o.Child(name).Set("test", Opaque{})
	}

	// Retrieve child names using the Children() method
	gotChildren := o.Children()

	// Verify we got the correct number of children
	assert.Equal(t, len(gotChildren), len(childNames))

	// Verify each child name is present in the returned slice
	for _, name := range childNames {
		found := false
		for _, child := range gotChildren {
			if child == name {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected to find child name %s, but didn't", name)
		}
	}
}
