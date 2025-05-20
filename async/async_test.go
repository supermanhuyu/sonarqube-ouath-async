package async

import (
	"testing"
)

func TestIsContain(t *testing.T) {
	items := []string{"a", "b", "c"}
	if !IsContain(items, "b") {
		t.Errorf("Expected true, got false")
	}
	if IsContain(items, "d") {
		t.Errorf("Expected false, got true")
	}
}

func TestGetAllGroupName(t *testing.T) {
	group := "foo.bar.baz"
	subGroups := GetAllGroupName(group)
	expected := []string{"foo", "foo.bar", "foo.bar.baz"}
	if len(subGroups) != len(expected) {
		t.Errorf("Expected %d, got %d", len(expected), len(subGroups))
	}
	for i, v := range expected {
		if subGroups[i] != v {
			t.Errorf("Expected %s, got %s", v, subGroups[i])
		}
	}
}
