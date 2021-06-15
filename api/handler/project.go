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

// For json validation of POST /projects/
type Post struct {
	Name string `json:"name" binding:"required"`
}

// Handler for POST /projects/
func PostProjectHandler(t store.TodoStore, c *gin.Context) {
	// validate json requestBody
	var json Post
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	projectName := json.Name

	// Check if project already exists
	project := t.GetProject(projectName)
	if project.Name != "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "project already existing"})
		return
	}

	// Create project
	err := t.PostProject(projectName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "project created",
	})
}
