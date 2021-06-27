package api_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/mpfen/Go-Todo-REST-API-V2/api/model"
	"gorm.io/gorm"
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

func (s *StubTodoStore) PostProject(name string) error {
	// Find out which ID the project gets
	lastProjectIndex := (len(s.Projects) - 1)
	id := s.Projects[lastProjectIndex].ID

	project := model.Project{Model: gorm.Model{
		ID:        uint(id),
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
		DeletedAt: gorm.DeletedAt{},
	},
		Name:     name,
		Archived: false,
		Tasks:    []model.Task{}}

	s.Projects = append(s.Projects, project)
	return nil
}

func (s *StubTodoStore) GetAllProjects() []model.Project {
	return s.Projects
}

func (s *StubTodoStore) UpdateProject(project model.Project) error {
	index := int(project.ID) - 1
	s.Projects[index].Name = project.Name
	return nil
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

// Converts an array of projects to json
func projectsToJson(t *testing.T, projects []model.Project) string {
	t.Helper()
	want, err := json.Marshal(projects)
	if err != nil {
		t.Errorf("Error parsing project to json: %s", err)
	}
	return string(want[:])
}
