package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mpfen/Go-Todo-REST-API-V2/api"
	"github.com/mpfen/Go-Todo-REST-API-V2/api/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupTaskTests() (server *api.TodoServer, store *StubTodoStore) {
	time := time.Time{}
	deletedAT := gorm.DeletedAt{}

	store = &StubTodoStore{
		[]model.Project{{
			Model: gorm.Model{
				ID:        uint(1),
				CreatedAt: time,
				UpdatedAt: time,
				DeletedAt: deletedAT,
			},
			Name:     "homework",
			Archived: false,
			Tasks:    []model.Task{},
		},
			{Model: gorm.Model{
				ID:        uint(2),
				CreatedAt: time,
				UpdatedAt: time,
				DeletedAt: deletedAT,
			},
				Name:     "cleaning",
				Archived: true,
				Tasks:    []model.Task{},
			},
			{Model: gorm.Model{
				ID:        uint(3),
				CreatedAt: time,
				UpdatedAt: time,
				DeletedAt: deletedAT,
			},
				Name:     "school",
				Archived: false,
				Tasks:    []model.Task{},
			}},
		[]model.Task{{
			Model: gorm.Model{
				ID:        uint(1),
				CreatedAt: time,
				UpdatedAt: time,
				DeletedAt: deletedAT,
			},
			Name:      "math",
			Priority:  "1",
			Deadline:  &time,
			Done:      false,
			ProjectID: uint(1),
		}, {
			Model: gorm.Model{
				ID:        uint(2),
				CreatedAt: time,
				UpdatedAt: time,
				DeletedAt: deletedAT,
			},
			Name:      "kitchen",
			Priority:  "1",
			Deadline:  &time,
			Done:      false,
			ProjectID: uint(2),
		}, {
			Model: gorm.Model{
				ID:        uint(3),
				CreatedAt: time,
				UpdatedAt: time,
				DeletedAt: deletedAT,
			},
			Name:      "pyhsics",
			Priority:  "1",
			Deadline:  &time,
			Done:      false,
			ProjectID: uint(1),
		},
		}}

	server = api.NewTodoServer(store)
	return server, store
}

// Test for Route POST /projects/:projectName/tasks
func TestPostTask(t *testing.T) {
	server, store := setupTaskTests()

	t.Run("Create new task biology for project homework", func(t *testing.T) {
		requestBody := makeNewPostTaskBody(t, "biology", true)
		req, _ := http.NewRequest("POST", "/projects/homework/tasks", requestBody)
		w := httptest.NewRecorder()
		server.Router.ServeHTTP(w, req)

		assert.Equalf(t, http.StatusCreated, w.Code, "wanted http.StatusCreated got %s", w.Code)
		assert.Len(t, store.Tasks, 4)
		assert.Equalf(t, store.Tasks[3].Name, "biology", "task was not created")
	})

	t.Run("Try to create a task for nonexisting project", func(t *testing.T) {
		requestBody := makeNewPostTaskBody(t, "biology", true)
		req, _ := http.NewRequest("POST", "/projects/homwork/tasks", requestBody)
		w := httptest.NewRecorder()
		server.Router.ServeHTTP(w, req)

		assert.Equalf(t, http.StatusNotFound, w.Code, "wanted http.StatusNotFound got %s, error: %s", w.Code)
		assert.Len(t, store.Tasks, 4)
	})

	t.Run("Try to create a task with an invalid requestBody", func(t *testing.T) {
		requestBody := makeNewPostTaskBody(t, "biology", false)
		req, _ := http.NewRequest("POST", "/projects/homework/tasks", requestBody)
		w := httptest.NewRecorder()
		server.Router.ServeHTTP(w, req)

		assert.Equalf(t, http.StatusBadRequest, w.Code, "wanted http.StatusBadRequest got %s, error: %s", w.Code)
		assert.Len(t, store.Tasks, 4)
	})
}

// Tests for Route GET /projects/:projectName/tasks/:taskname
func TestGetTask(t *testing.T) {
	server, store := setupTaskTests()

	t.Run("Get Task 'math' from project 'homework'", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/projects/homework/tasks/math", nil)
		w := httptest.NewRecorder()
		server.Router.ServeHTTP(w, req)

		assert.Equalf(t, http.StatusOK, w.Code, "wanted http.StatusOK got %s", w.Code)
		want := taskToJSON(t, store.Tasks[0])
		assert.JSONEqf(t, want, w.Body.String(), "wanted %s, got %s", want, w.Body.String())
	})

	t.Run("Try to get task from nonexisting project", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/projects/homwork/tasks/math", nil)
		w := httptest.NewRecorder()
		server.Router.ServeHTTP(w, req)

		assert.Equalf(t, http.StatusNotFound, w.Code, "wanted http.StatusNotFound got %s", w.Code)

		// Check for correct message
		want := `{"message": "project not found"}`
		assert.JSONEqf(t, want, w.Body.String(), "wanted %s, got %s", want, w.Body.String())
	})

	t.Run("Try to get nonexistent task from project", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/projects/homework/tasks/math2", nil)
		w := httptest.NewRecorder()
		server.Router.ServeHTTP(w, req)

		assert.Equalf(t, http.StatusNotFound, w.Code, "wanted http.StatusNotFound got %s", w.Code)

		// Check for correct message
		want := `{"message": "task not found"}`
		assert.JSONEqf(t, want, w.Body.String(), "wanted %s, got %s", want, w.Body.String())
	})

}

// Test for Route GET /projects/:projectName/tasks
func TestGetAllTasks(t *testing.T) {
	server, store := setupTaskTests()

	t.Run("Get all task from project ’homework’", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/projects/homework/tasks", nil)
		w := httptest.NewRecorder()
		server.Router.ServeHTTP(w, req)

		assert.Equalf(t, http.StatusOK, w.Code, "wanted http.StatusOK got %s", w.Code)

		homeworkTasks := append([]model.Task{}, store.Tasks[0], store.Tasks[2])
		want := tasksToJson(t, homeworkTasks)
		assert.JSONEqf(t, want, w.Body.String(), "wanted %s, got %s", want, w.Body.String())
	})

	t.Run("Get all Tasks from a project without tasks", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/projects/school/tasks", nil)
		w := httptest.NewRecorder()
		server.Router.ServeHTTP(w, req)

		assert.Equalf(t, http.StatusOK, w.Code, "wanted http.StatusOK got %s", w.Code)
		assert.Equalf(t, "[]", w.Body.String(), "wanted empty response, got %s", w.Body.String())
	})
}

// Test for Route PUt /projects/:projectName/task/:taskName
func TestUpdateTask(t *testing.T) {
	server, store := setupTaskTests()

	t.Run("Update task 'math' from project 'homework' to 'mathexam'", func(t *testing.T) {
		putRequestBody := makeNewPostTaskBody(t, "mathexam", true)
		req, _ := http.NewRequest("PUT", "/projects/homework/tasks/math", putRequestBody)
		w := httptest.NewRecorder()
		server.Router.ServeHTTP(w, req)

		assert.Equalf(t, http.StatusOK, w.Code, "wanted http.StatusOK got %s", w.Code)
		assert.Equalf(t, store.Tasks[0].Name, "mathexam", "wanted 'mathexam', got %s", store.Tasks[0].Name)
	})

	t.Run("Try to update nonexistent task", func(t *testing.T) {
		putRequestBody := makeNewPostTaskBody(t, "mathexam", true)
		req, _ := http.NewRequest("PUT", "/projects/homework/tasks/math2", putRequestBody)
		w := httptest.NewRecorder()
		server.Router.ServeHTTP(w, req)

		assert.Equalf(t, http.StatusNotFound, w.Code, "wanted http.StatusNotFound got %s", w.Code)
	})

	t.Run("Try to update task of nonexistent project", func(t *testing.T) {
		putRequestBody := makeNewPostTaskBody(t, "mathexam", true)
		req, _ := http.NewRequest("PUT", "/projects/homework2/tasks/math", putRequestBody)
		w := httptest.NewRecorder()
		server.Router.ServeHTTP(w, req)

		assert.Equalf(t, http.StatusNotFound, w.Code, "wanted http.StatusNotFound got %s", w.Code)
	})
}

// Tests for Route DELETE /projects/:projectName/tasks/:taskName
func TestDeleteTask(t *testing.T) {
	server, store := setupTaskTests()

	t.Run("Delete task 'math' from project 'homework'", func(t *testing.T) {
		mathTask := store.Tasks[0]

		req, _ := http.NewRequest("DELETE", "/projects/homework/tasks/math", nil)
		w := httptest.NewRecorder()
		server.Router.ServeHTTP(w, req)

		assert.Equalf(t, http.StatusOK, w.Code, "wanted http.StatusOK, got %s", w.Code)
		assert.NotContains(t, store.Tasks, mathTask, "Task was not deleted")
	})

	t.Run("Try to delete nonexistent task", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/projects/homework/tasks/math2", nil)
		w := httptest.NewRecorder()
		server.Router.ServeHTTP(w, req)

		assert.Equalf(t, http.StatusNotFound, w.Code, "wanted http.StatusNotFound, got %s", w.Code)
		assert.Lenf(t, store.Tasks, 2, "No task should have been deleted")
	})
}

// Tests for Route PUT/DELETE /projects/:projectName/tasks/:taskName/complete
// Task completion and undoing
func TestCompleteAndUndoTask(t *testing.T) {
	server, store := setupTaskTests()

	t.Run("Complete task 'math' from project 'homework'", func(t *testing.T) {
		req, _ := http.NewRequest("PUT", "/projects/homework/tasks/math/complete", nil)
		w := httptest.NewRecorder()
		server.Router.ServeHTTP(w, req)

		assert.Equalf(t, http.StatusOK, w.Code, "wanted http.StatusOK, got %s", w.Code)
		assert.Equal(t, true, store.Tasks[0].Done, "task was not completed")
	})

	t.Run("Try to complete nonexistent task", func(t *testing.T) {
		req, _ := http.NewRequest("PUT", "/projects/homework2/tasks/math2/complete", nil)
		w := httptest.NewRecorder()
		server.Router.ServeHTTP(w, req)

		assert.Equalf(t, http.StatusNotFound, w.Code, "wanted http.StatusNotFound, got %s", w.Code)
	})

	t.Run("Undo task 'math' from project 'homework'", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/projects/homework/tasks/math/complete", nil)
		w := httptest.NewRecorder()
		server.Router.ServeHTTP(w, req)

		assert.Equalf(t, http.StatusOK, w.Code, "wanted http.StatusOK, got %s", w.Code)
		assert.Equal(t, false, store.Tasks[0].Done, "task was not undone")
	})

	t.Run("Try to undo nonexistent task", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/projects/homework2/tasks/math2/complete", nil)
		w := httptest.NewRecorder()
		server.Router.ServeHTTP(w, req)

		assert.Equalf(t, http.StatusNotFound, w.Code, "wanted http.StatusNotFound, got %s", w.Code)
	})
}

// makes a new json request body for POST /projects/:projectName/tasks
// if !valid an invalid requestBody is returned
func makeNewPostTaskBody(t *testing.T, name string, vaild bool) *bytes.Buffer {
	if vaild {
		requestBody, err := json.Marshal(map[string]string{
			"name":     name,
			"priority": "1",
			"deadline": time.Time{}.String(),
		})
		if err != nil {
			t.Errorf("Failed to make requestBody: %s", err)
		}

		return bytes.NewBuffer(requestBody)
	} else {
		requestBody, err := json.Marshal(map[string]string{
			"invaild": name,
		})
		if err != nil {
			t.Errorf("Failed to make requestBody: %s", err)
		}

		return bytes.NewBuffer(requestBody)
	}
}
