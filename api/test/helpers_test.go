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
	if index >= 0 {
		s.Projects[index].Name = project.Name
		s.Projects[index].Archived = project.Archived
	}

	return nil
}

func (s *StubTodoStore) DeleteProject(projectName string) error {
	for i, project := range s.Projects {
		if project.Name == projectName {
			s.Projects = append(s.Projects[:i], s.Projects[(i+1):]...)
			return nil
		}
	}
	return gorm.ErrRecordNotFound
}

func (s *StubTodoStore) PostTask(task model.Task) error {
	time := time.Time{}
	deletedAT := gorm.DeletedAt{}

	newTask := model.Task{
		Model: gorm.Model{
			ID:        uint(42),
			CreatedAt: time,
			UpdatedAt: time,
			DeletedAt: deletedAT,
		},
		Name:      task.Name,
		Priority:  task.Priority,
		Deadline:  &time,
		Done:      false,
		ProjectID: task.ProjectID,
	}

	s.Tasks = append(s.Tasks, newTask)
	return nil
}

func (s *StubTodoStore) GetTask(projectName, taskName string) model.Task {
	var projectID uint
	for _, project := range s.Projects {
		if project.Name == projectName {
			projectID = project.ID
		}
	}
	for i, task := range s.Tasks {
		if task.Name == taskName && task.ProjectID == projectID {
			return s.Tasks[i]
		}
	}
	return model.Task{}
}

func (s *StubTodoStore) GetAllProjectTasks(project model.Project) []model.Task {
	projects := []model.Task{}

	for _, projectIter := range s.Tasks {
		if projectIter.ProjectID == project.ID {
			projects = append(projects, projectIter)
		}
	}
	return projects
}

func (s *StubTodoStore) UpdateTask(task model.Task) error {
	index := int(task.ID) - 1

	if index >= 0 {
		s.Tasks[index].Name = task.Name
		s.Tasks[index].Deadline = task.Deadline
		s.Tasks[index].Priority = task.Priority
		s.Tasks[index].Done = task.Done
	}
	return nil
}

func (s *StubTodoStore) DeleteTask(task model.Task) error {
	index := int(task.ID) - 1
	s.Tasks = append(s.Tasks[:index], s.Tasks[(index+1):]...)

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

// Converts a task struct to json
func taskToJSON(t *testing.T, task model.Task) string {
	t.Helper()
	want, err := json.Marshal(task)
	if err != nil {
		t.Errorf("Error parsing task to json: %s", err)
	}
	return string(want[:])
}

func tasksToJson(t *testing.T, tasks []model.Task) string {
	t.Helper()
	want, err := json.Marshal(tasks)
	if err != nil {
		t.Errorf("Error parsing task to json: %s", err)
	}
	return string(want[:])
}
