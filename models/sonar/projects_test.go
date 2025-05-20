package sonar

import (
	"testing"
)

func TestProject_TableName(t *testing.T) {
	p := &Project{}
	if p.TableName() != "projects" {
		t.Errorf("Expected 'projects', got '%s'", p.TableName())
	}
}
