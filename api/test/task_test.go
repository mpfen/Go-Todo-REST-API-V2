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

// Test for Route POST /projects/:name/tasks
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

// makes a new json request body for POST /projects/:name/tasks
// if !valid a invalid requestBody is returned
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
