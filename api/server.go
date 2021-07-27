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
	t.Router.GET("/projects/:projectName", t.GetProject)
	t.Router.POST("/projects/", t.PostProject)
	t.Router.GET("/projects/", t.GetAllProjects)
	t.Router.PUT("/projects/:projectName", t.PutProject)
	t.Router.DELETE("/projects/:projectName", t.DeleteProject)
	t.Router.DELETE("/projects/:projectName/archive", t.ArchiveProject)
	t.Router.PUT("/projects/:projectName/archive", t.ArchiveProject)

	// Task routes
	t.Router.POST("projects/:projectName/tasks", t.PostTask)
	t.Router.GET("projects/:projectName/tasks/:taskName", t.GetTask)
	t.Router.GET("projects/:projectName/tasks", t.GetAllTasks)
	t.Router.PUT("projects/:projectName/tasks/:taskName", t.PutTask)
	t.Router.DELETE("projects/:projectName/tasks/:taskName", t.DeleteTask)

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

func (t *TodoServer) PutProject(c *gin.Context) {
	handler.PutProjectHandler(t.Store, c)
}

func (t *TodoServer) DeleteProject(c *gin.Context) {
	handler.DeleteProjectHandler(t.Store, c)
}

func (t *TodoServer) ArchiveProject(c *gin.Context) {
	handler.ArchiveProjectHandler(t.Store, c)
}

// Task Handlers
func (t *TodoServer) PostTask(c *gin.Context) {
	handler.PostTaskHandler(t.Store, c)
}

func (t *TodoServer) GetTask(c *gin.Context) {
	handler.GetTaskHandler(t.Store, c)
}

func (t *TodoServer) GetAllTasks(c *gin.Context) {
	handler.GetAllTasksHandler(t.Store, c)
}

func (t *TodoServer) PutTask(c *gin.Context) {
	handler.PutTaskHandler(t.Store, c)
}

func (t *TodoServer) DeleteTask(c *gin.Context) {
	handler.DeleteTaskHandler(t.Store, c)
}
