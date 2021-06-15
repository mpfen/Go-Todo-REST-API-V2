package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mpfen/Go-Todo-REST-API-V2/api/store"
)

// Handler for GET /projects/:name
func GetProjectHandler(t store.TodoStore, c *gin.Context) {
	projectName := c.Param("name")

	project := t.GetProject(projectName)

	// Check if no project with that name was found
	if project.Name == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "project not found",
		})
		return
	}
	c.JSON(http.StatusOK, project)
}
