package api_test

import (
	"encoding/json"
	"testing"

	"github.com/mpfen/Go-Todo-REST-API-V2/api/model"
)

type StubTodoStore struct {
	Projects []model.Project
	Tasks    []model.Task
}

func (s *StubTodoStore) GetProject(name string) model.Project {
	for _, p := range s.Projects {
		if p.Name == name {
			return p
		}
	}
	return model.Project{}
}

// Converts a project struct to json
func projectToJson(t *testing.T, project model.Project) string {
	t.Helper()
	want, err := json.Marshal(project)
	if err != nil {
		t.Errorf("Error parsing project to json: %s", err)
	}
	return string(want[:])
}
