package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mpfen/Go-Todo-REST-API-V2/api/store"
	"gorm.io/gorm"
)

// Handler for GET /projects/:projectName
func GetProjectHandler(t store.TodoStore, c *gin.Context) {
	projectName := c.Param("projectName")

	project := checkIfProjectExistsOr404(t, c, projectName)
	if project.Name == "" {
		return
	}

	c.JSON(http.StatusOK, project)
}

// Handler for POST /projects/
func PostProjectHandler(t store.TodoStore, c *gin.Context) {
	// validate json requestBody
	var json Post
	if err := c.ShouldBindJSON(&json); err != nil {
		sendJSONResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	projectName := json.Name

	// Check if project already exists
	project := t.GetProject(projectName)
	if project.Name != "" {
		sendJSONResponse(c, http.StatusBadRequest, "project already existing")
		return
	}

	// Create project
	err := t.PostProject(projectName)
	if err != nil {
		sendJSONResponse(c, http.StatusInternalServerError, err.Error())
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

// Handler for PUT /projects/:projectName
func PutProjectHandler(t store.TodoStore, c *gin.Context) {
	// validate json requestBody
	var json Post
	if err := c.ShouldBindJSON(&json); err != nil {
		sendJSONResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	oldProjectName := c.Param("projectName")
	newProjectName := json.Name

	// Check if project exists
	project := checkIfProjectExistsOr404(t, c, oldProjectName)

	// Update Project
	project.Name = newProjectName
	err := t.UpdateProject(project)
	if err != nil {
		sendJSONResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "project updated",
	})
}

// Handler for DELETE /projects/:projectName
func DeleteProjectHandler(t store.TodoStore, c *gin.Context) {
	// Try to delete project
	projectName := c.Param("projectName")
	err := t.DeleteProject(projectName)

	// Check error if no project was found
	if err == gorm.ErrRecordNotFound {
		sendJSONResponse(c, http.StatusNotFound, "project not found")
		return
	} else if err != nil {
		sendJSONResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "project deleted",
	})
}

// Handler for PUT/DELETE /projects/:projectName/archive
func ArchiveProjectHandler(t store.TodoStore, c *gin.Context) {
	projectName := c.Param("projectName")

	// Check if project exists
	project := checkIfProjectExistsOr404(t, c, projectName)
	if project.Name == "" {
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
		sendJSONResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": responseText,
	})
}
