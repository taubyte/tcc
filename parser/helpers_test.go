package parser

import (
	"testing"
)

func TestIntAttribute(t *testing.T) {
	name := "testInt"
	attr := Int(name)
	if attr.Name != name || attr.Type != TypeInt {
		t.Errorf("Expected Int attribute with name %s and type %d, got %s and %d", name, TypeInt, attr.Name, attr.Type)
	}
}

func TestBoolAttribute(t *testing.T) {
	name := "testBool"
	attr := Bool(name)
	if attr.Name != name || attr.Type != TypeBool {
		t.Errorf("Expected Bool attribute with name %s and type %d, got %s and %d", name, TypeBool, attr.Name, attr.Type)
	}
}

func TestFloatAttribute(t *testing.T) {
	name := "testFloat"
	attr := Float(name)
	if attr.Name != name || attr.Type != TypeFloat {
		t.Errorf("Expected Float attribute with name %s and type %d, got %s and %d", name, TypeFloat, attr.Name, attr.Type)
	}
}

func TestStringAttribute(t *testing.T) {
	name := "testString"
	attr := String(name)
	if attr.Name != name || attr.Type != TypeString {
		t.Errorf("Expected String attribute with name %s and type %d, got %s and %d", name, TypeString, attr.Name, attr.Type)
	}
}

func TestStringSliceAttribute(t *testing.T) {
	name := "testStringSlice"
	attr := StringSlice(name)
	if attr.Name != name || attr.Type != TypeStringSlice {
		t.Errorf("Expected StringSlice attribute with name %s and type %d, got %s and %d", name, TypeStringSlice, attr.Name, attr.Type)
	}
}

func TestDefine(t *testing.T) {
	match := "testDefine"
	attrs := Attributes(Int("testInt"), String("testString"))
	node := Define(match, attrs)
	if node.Match != match || len(node.Attributes) != 2 || len(node.Children) != 0 {
		t.Errorf("Expected node with match %s, 2 attributes, and 0 children. Got match %s, %d attributes, and %d children", match, node.Match, len(node.Attributes), len(node.Children))
	}
}

func TestDefineGroup(t *testing.T) {
	match := "testDefineGroup"
	node := DefineGroup(match)
	if node.Match != match || len(node.Attributes) != 0 || len(node.Children) != 0 {
		t.Errorf("Expected node with match %s, 0 attributes, and 0 children. Got match %s, %d attributes, and %d children", match, node.Match, len(node.Attributes), len(node.Children))
	}
}
