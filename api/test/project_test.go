package api_test

import (
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

		assert.Equalf(t, http.StatusNotFound, w.Code, "wanted %s got %s", http.StatusNotFound, w.Code)
	})
}
