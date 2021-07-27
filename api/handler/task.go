package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mpfen/Go-Todo-REST-API-V2/api/model"
	"github.com/mpfen/Go-Todo-REST-API-V2/api/store"
)

// Handler for POST /projects/:projectName/tasks
func PostTaskHandler(t store.TodoStore, c *gin.Context) {
	// validate json requestBody
	var json Task
	if err := c.ShouldBindJSON(&json); err != nil {
		sendJSONResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	projectName := c.Param("projectName")

	// Check if project exists
	project := checkIfProjectExistsOr404(t, c, projectName)
	if project.Name == "" {
		return
	}

	// Create Task
	task := model.Task{}
	task.Name = json.Name
	task.Priority = json.Priority

	// Parse the deadline string as a time.Time{}
	deadline, errTime := time.Parse("2006-01-02 15:04:05 +0000 UTC", json.Deadline)
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

// Handler for Route GET /projects/:projectName/tasks/:taskName
func GetTaskHandler(t store.TodoStore, c *gin.Context) {
	projectName := c.Param("projectName")
	taskName := c.Param("taskName")

	// Check if project exists
	project := checkIfProjectExistsOr404(t, c, projectName)
	if project.Name == "" {
		return
	}

	// Check if task exists
	task := checkIfTaskExistsOr404(t, c, projectName, taskName)
	if task.Name == "" {
		return
	}
	c.JSON(http.StatusOK, task)
}

// Handler for Route GET /projects/:projectName/tasks
func GetAllTasksHandler(t store.TodoStore, c *gin.Context) {
	projectName := c.Param("projectName")

	// Check if project exists
	project := checkIfProjectExistsOr404(t, c, projectName)
	if project.Name == "" {
		return
	}

	// Get all Tasks of the project
	tasks := t.GetAllProjectTasks(project)
	c.JSON(http.StatusOK, tasks)
}

// Handler for Route PUT /projects/:projectName/tasks/:taskName
func PutTaskHandler(t store.TodoStore, c *gin.Context) {
	// validate json requestBody
	var jsonTask Task
	if err := c.ShouldBindJSON(&jsonTask); err != nil {
		sendJSONResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	projectName := c.Param("projectName")

	// Check if project exists
	project := checkIfProjectExistsOr404(t, c, projectName)
	if project.Name == "" {
		return
	}

	// Check if task exists
	oldTaskName := c.Param("taskName")
	oldTask := checkIfTaskExistsOr404(t, c, projectName, oldTaskName)
	if oldTask.Name == "" {
		return
	}

	// Update task
	oldTask.Name = jsonTask.Name
	oldTask.Priority = jsonTask.Priority
	// Parse the deadline string as a time.Time{}
	deadline, errTime := time.Parse("2006-01-02 15:04:05 +0000 UTC", jsonTask.Deadline)
	if errTime != nil {
		sendJSONResponse(c, http.StatusBadRequest, errTime.Error())
		return
	}
	oldTask.Deadline = &deadline

	err := t.UpdateTask(oldTask)

	if err != nil {
		sendJSONResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	sendJSONResponse(c, http.StatusOK, "task updated")
}

// Handler for Route DELETE /projects/:projectName/tasks/:taskName
func DeleteTaskHandler(t store.TodoStore, c *gin.Context) {
	projectName := c.Param("projectName")
	taskName := c.Param("taskName")

	// Check if task exists
	task := checkIfTaskExistsOr404(t, c, projectName, taskName)
	if task.Name == "" {
		return
	}

	err := t.DeleteTask(task)
	if err != nil {
		sendJSONResponse(c, http.StatusInternalServerError, err.Error())
	}

	sendJSONResponse(c, http.StatusOK, "task deleted")
}
