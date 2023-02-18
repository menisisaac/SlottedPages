package main

import (
	"fmt"
	"testing"
)

const slottedPageSize = 2048

func TestAdd(t *testing.T) {
	p := createPage(0, slottedPageSize)
	add(t, p, "123", 0)
	add(t, p, "456", 1)
	add(t, p, "789", 2)
}

func TestGet(t *testing.T) {
	p := createPage(0, slottedPageSize)
	p.add("123")
	p.add("456")
	p.add("789")
	p.saveLocation(1, 0, 0)
	if 3 != p.entryCount() {
		t.Errorf("Entry count should be 3, got %q", p.entryCount())
	}
	if p.get(-1) != nil {
		t.Error("Should be out of bounds")
	}
	if p.get(0) != "123" {
		t.Errorf("%s should be %s", p.get(0), "123")
	}
	println("Hi")
	if p.get(1) != nil {
		t.Error("Should be nil")
	}
	if p.get(2) != "789" {
		t.Errorf("%s should be %s", p.get(2), "789")
	}
	if p.get(3) != nil {
		t.Error("Should be out of bounds")
	}

}

func add(t *testing.T, p slottedPage, s string, i int) {
	index := p.add(s)
	if i != index {
		t.Errorf("Expected Index: %q Actual Index: %q", i, index)
	}
	if i+1 != p.entryCount() {
		t.Errorf("Expected Entry Count: %q Actual Entry Count: %q", i+1, p.entryCount())
	}
	object := fmt.Sprintf("%s", p.get(index))
	if s != object {
		t.Errorf("Expected Object: %sActual Object: %s", s, object)
	}
}
