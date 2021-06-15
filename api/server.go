package api

import (
	"github.com/gin-gonic/gin"
	"github.com/mpfen/Go-Todo-REST-API-V2/api/handler"
	"github.com/mpfen/Go-Todo-REST-API-V2/api/store"
)

type TodoServer struct {
	Router *gin.Engine
	Store  store.TodoStore
}

// Initialize TodoServer and create a gin router
func NewTodoServer(store store.TodoStore) *TodoServer {
	t := new(TodoServer)
	t.Store = store
	t.Router = gin.Default()

	// Project routes
	t.Router.GET("/projects/:name", t.GetProject)
	t.Router.POST("/projects/", t.PostProject)
	t.Router.GET("/projects/", t.GetAllProjects)

	return t
}

// Project Handlers
func (t *TodoServer) GetProject(c *gin.Context) {
	handler.GetProjectHandler(t.Store, c)
}

func (t *TodoServer) PostProject(c *gin.Context) {
	handler.PostProjectHandler(t.Store, c)
}

func (t *TodoServer) GetAllProjects(c *gin.Context) {
	handler.GetAllProjectsHandler(t.Store, c)
}
