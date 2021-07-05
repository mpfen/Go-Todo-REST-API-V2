package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mpfen/Go-Todo-REST-API-V2/api/model"
	"github.com/mpfen/Go-Todo-REST-API-V2/api/store"
)

// For json validation of POST /projects/:name
type Post struct {
	Name string `json:"name" binding:"required"`
}

// For json validation of POST /projects/:name/tasks
type Task struct {
	Name     string `json:"name" binding:"required"`
	Priority string `json:"priority" binding:"required"`
	Deadline string `json:"deadline" binding:"required"`
}

// Checks if a project with that name exist.
// If no project is found the context is aborted and a response is send
func checkIfProjectExistsOr404(t store.TodoStore, c *gin.Context, projectName string) model.Project {
	project := t.GetProject(projectName)

	// If not project was found the struct is "empty"
	if project.Name == "" {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "project not found",
		})
	}
	return project
}

// Checks if a task with that name exist in the project
// If no task in that project is found the context is aborted and a response is send
func checkIfTaskExistsOr404(t store.TodoStore, c *gin.Context, projectName, taskName string) model.Task {
	task := t.GetTask(projectName, taskName)

	// If no task was found the struct is "empty"
	if task.Name == "" {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "task not found",
		})
	}
	return task
}

// Send a JSON response and ends the context
func sendJSONResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{
		"message": message,
	})
}
