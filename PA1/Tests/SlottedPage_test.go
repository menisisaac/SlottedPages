package Tests

import (
	. "PA1/Structure"
	"fmt"
	"testing"
)

const slottedPageSize = 2048

func TestAdd(t *testing.T) {
	p := CreatePage(0, slottedPageSize)
	add(t, p, "123", 0)
	add(t, p, "456", 1)
	add(t, p, "789", 2)
}

func TestGet(t *testing.T) {
	p := CreatePage(0, slottedPageSize)
	p.Add("123")
	p.Add("456")
	p.Add("789")
	p.SaveLocation(1, -1, -1)
	if 3 != p.EntryCount() {
		t.Errorf("Entry count should be 3, got %q", p.EntryCount())
	}
	if p.Get(-1) != nil {
		t.Error("Should be out of bounds")
	}
	object := fmt.Sprintf("%s", p.Get(0))
	if object != "123" {
		t.Errorf("%s should be %s", object, "123")
	}
	if p.Get(1) != nil {
		t.Error("Should be nil")
	}
	if p.Get(2) != "789" {
		t.Errorf("%s should be %s", p.Get(2), "789")
	}
	if p.Get(3) != nil {
		t.Error("Should be out of bounds")
	}
}
func TestRemove(t *testing.T) {
	p := CreatePage(0, slottedPageSize)
	p.Add("123")
	p.Add("456")
	p.Add("789")
	temp := p.Remove(1)
	if "456" != temp {
		t.Errorf("Removing index 1 should return 456, got %s", temp)
	}
	if nil != p.Remove(1) {
		t.Errorf("Removing index again should remove nil")
	}
	if p.EntryCount() != 3 {
		t.Error("Entry count should be 3")
	}
	if p.Get(0) != "123" {
		t.Error("Index one should not have changed")
	}
	if p.Get(2) != "789" {
		t.Error("Index 3 should lead to 789")
	}
}
func TestCompact(t *testing.T) {
	p := CreatePage(0, slottedPageSize)
	for x := 0; x < 50; x += 1 {
		p.Add(56)
		if x%2 != 0 {
			p.Remove(x)
		}
	}

	for x := 0; x < p.EntryCount(); x += 1 {
		temp := p.Get(x)
		if temp != nil {
			fmt.Printf("%s\n", temp)
		}
	}
}

func add(t *testing.T, p SlottedPage, s string, i int) {
	index := p.Add(s)
	if i != index {
		t.Errorf("Expected Index: %q Actual Index: %q", i, index)
	}
	if i+1 != p.EntryCount() {
		t.Errorf("Expected Entry Count: %q Actual Entry Count: %q", i+1, p.EntryCount())
	}
	object := fmt.Sprintf("%s", p.Get(index))
	if s != object {
		t.Errorf("Expected Object: %sActual Object: %s", s, object)
	}
}
