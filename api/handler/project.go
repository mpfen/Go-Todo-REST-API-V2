package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mpfen/Go-Todo-REST-API-V2/api/store"
	"gorm.io/gorm"
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

// Handler for GET /projects/
func GetAllProjectsHandler(t store.TodoStore, c *gin.Context) {
	projects := t.GetAllProjects()
	c.JSON(http.StatusOK, projects)
}

// Handler for PUT /projects/:name
func PutProjectHandler(t store.TodoStore, c *gin.Context) {
	// validate json requestBody
	var json Post
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	oldProjectName := c.Param("name")
	newProjectName := json.Name

	// Check if project exists
	project := t.GetProject(oldProjectName)
	if project.Name == "" {
		c.JSON(http.StatusNotFound, gin.H{"message": "project not found"})
		return
	}

	// Update Project
	project.Name = newProjectName
	err := t.UpdateProject(project)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "project updated",
	})
}

// Handler for DELETE /projects/:name
func DeleteProjectHandler(t store.TodoStore, c *gin.Context) {
	projectName := c.Param("name")
	err := t.DeleteProject(projectName)

	// Check error if no project was found
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "project not found",
		})
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "project deleted",
	})
}

// Handler for PUT/DELETE /projects/:name/archive
func ArchiveProjectHandler(t store.TodoStore, c *gin.Context) {
	projectName := c.Param("name")

	// Check if project exists
	project := t.GetProject(projectName)
	if project.Name == "" {
		c.JSON(http.StatusNotFound, gin.H{"message": "project not found"})
		return
	}

	// Depending on method archive or unarchive project
	var responseText string
	if c.Request.Method == "PUT" {
		project.ArchiveProject()
		responseText = "project archived"

	} else {
		project.UnArchiveProject()
		responseText = "project unarchived"
	}

	// Update project
	err := t.UpdateProject(project)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": responseText,
	})
}
