package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mpfen/Go-Todo-REST-API-V2/api/model"
	"github.com/mpfen/Go-Todo-REST-API-V2/api/store"
)

// Handler for POST /projects/:name/tasks
func PostTaskHandler(t store.TodoStore, c *gin.Context) {
	// validate json requestBody
	var json Task
	if err := c.ShouldBindJSON(&json); err != nil {
		sendJSONResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	projectName := c.Param("name")

	// Check if project exists
	project := checkIfProjectExistsOr404(t, c, projectName)
	if project.Name == "" {
		return
	}

	// Create Task
	task := model.Task{}
	task.Name = json.Name
	task.Priority = json.Priority
	deadline, errTime := time.Parse("2006-01-02 15:04:05 +0000 UTC", json.Deadline)

	// Parse the deadline string as a time.Time
	if errTime != nil {
		sendJSONResponse(c, http.StatusBadRequest, errTime.Error())
		return
	}

	task.Deadline = &deadline
	task.ProjectID = project.ID

	err := t.PostTask(task)

	if err != nil {
		sendJSONResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	sendJSONResponse(c, http.StatusCreated, "task created")
}