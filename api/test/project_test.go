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

func setupProjectTests() (server *api.TodoServer, store *StubTodoStore) {
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
		[]model.Task{},
	}

	server = api.NewTodoServer(store)
	return server, store
}

// Tests for route GET /project/:name
func TestGetProject(t *testing.T) {
	server, store := setupProjectTests()

	t.Run("Get Project homework", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/projects/homework", nil)
		w := httptest.NewRecorder()
		server.Router.ServeHTTP(w, req)

		assert.Equalf(t, http.StatusOK, w.Code, "wanted %s got %s", http.StatusOK, w.Code)
		want := projectToJson(t, store.Projects[0])
		assert.JSONEq(t, want, w.Body.String())
	})

	t.Run("Get Project cleaning", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/projects/cleaning", nil)
		w := httptest.NewRecorder()
		server.Router.ServeHTTP(w, req)

		assert.Equalf(t, http.StatusOK, w.Code, "wanted %s got %s", http.StatusOK, w.Code)
		want := projectToJson(t, store.Projects[1])
		assert.JSONEq(t, want, w.Body.String())
	})

	t.Run("returns http.StatusNotFound on nonexistent projects", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/projects/work", nil)
		w := httptest.NewRecorder()
		server.Router.ServeHTTP(w, req)

		assert.Equalf(t, http.StatusNotFound, w.Code, "wanted http.StatusNotFound got %s", w.Code)
	})
}

// Test for Route POST /Projects/
func TestPostProject(t *testing.T) {
	server, store := setupProjectTests()

	t.Run("Create project exams", func(t *testing.T) {
		requestBody := makeNewPostProjectBody(t, "exams", true)
		req, _ := http.NewRequest("POST", "/projects/", requestBody)
		w := httptest.NewRecorder()
		server.Router.ServeHTTP(w, req)

		if assert.Equalf(t, http.StatusCreated, w.Code, "wanted http.StatusCreated got %s", w.Code) {
			assert.Equalf(t, "exams", store.Projects[3].Name, "project was not created")
		}

	})

	t.Run("Try to create project with invalid request body", func(t *testing.T) {
		requestBody := makeNewPostProjectBody(t, "exams", false)
		req, _ := http.NewRequest("POST", "/projects/", requestBody)
		w := httptest.NewRecorder()
		server.Router.ServeHTTP(w, req)

		assert.Equalf(t, http.StatusBadRequest, w.Code, "wanted http.StatusBadRequest got: %s", w.Code)
		assert.Len(t, store.Projects, 4)
	})

	t.Run("Try to create an exiting project", func(t *testing.T) {
		requestBody := makeNewPostProjectBody(t, "homework", true)
		req, _ := http.NewRequest("POST", "/projects/", requestBody)
		w := httptest.NewRecorder()
		server.Router.ServeHTTP(w, req)

		assert.Equalf(t, http.StatusBadRequest, w.Code, "wanted http.StatusBadRequest got: %s", w.Code)
		assert.Len(t, store.Projects, 4)
	})

}

// Test for Route GET /projects/
func TestGetAllProjects(t *testing.T) {
	server, store := setupProjectTests()

	req, _ := http.NewRequest("GET", "/projects/", nil)
	w := httptest.NewRecorder()
	server.Router.ServeHTTP(w, req)

	assert.Equalf(t, http.StatusOK, w.Code, "wanted http.StatusOK got %s", w.Code)

	gotJSON := w.Body.String()
	want := projectsToJson(t, store.Projects)
	assert.JSONEqf(t, want, gotJSON, "wanted %s got %s", want, gotJSON)
}

// makes a new json request body for POST /projects/
// if valid == false a invalid requestBody is returned
func makeNewPostProjectBody(t *testing.T, name string, vaild bool) *bytes.Buffer {
	if vaild {
		requestBody, err := json.Marshal(map[string]string{
			"name": name,
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
