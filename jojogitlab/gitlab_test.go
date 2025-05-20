package jojogitlab

import (
	"testing"

	"github.com/xanzy/go-gitlab"
)

func TestChangeGitlabProjectsToMap(t *testing.T) {
	// This is a simple test for the ChangeGitlabProjectsToMap function
	// Since it requires *gitlab.Project, we only test the mapping logic with minimal struct
	// For real tests, use mocks or interfaces
	projects := []*ProjectMock{
		{Namespace: NamespaceMock{FullPath: "foo/bar"}, Path: "baz", ID: 1},
		{Namespace: NamespaceMock{FullPath: "abc/def"}, Path: "xyz", ID: 2},
	}
	var converted []*gitlab.Project
	for _, p := range projects {
		converted = append(converted, p.ToGitlabProject())
	}
	m := ChangeGitlabProjectsToMap(converted)
	if len(m) != 2 {
		t.Errorf("Expected 2, got %d", len(m))
	}
}

// Mocks for testing

type NamespaceMock struct {
	FullPath string
}

type ProjectMock struct {
	Namespace NamespaceMock
	Path      string
	ID        int
}

func (p *ProjectMock) ToGitlabProject() *gitlab.Project {
	return &gitlab.Project{
		Namespace: &gitlab.ProjectNamespace{FullPath: p.Namespace.FullPath},
		Path:      p.Path,
		ID:        p.ID,
	}
}
