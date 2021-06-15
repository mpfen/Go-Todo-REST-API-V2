package store

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	model "github.com/mpfen/Go-Todo-REST-API-V2/api/model"
)

// TodoStore interface for testing
// Tests use own implementation with
// StubTodoStore instead of a real database
type TodoStore interface {
	GetProject(name string) model.Project
	PostProject(name string) error
	//GetAllProjects() []model.Project
	//DeleteProject(name string) error
	//UpdateProject(project model.Project) error

	//GetTask(projectName, taskName string) model.Task
	//PostTask(task model.Task) error
	//GetAllProjectTasks(project model.Project) []model.Task
	//DeleteTask(task model.Task) error
	//UpdateTask(task model.Task) error
}

type Database struct {
	DB *gorm.DB
}

// Gets project by name
func (d *Database) GetProject(name string) model.Project {
	project := model.Project{}
	err := d.DB.Find(&project, "Name = ?", name).Error

	if err != nil {
		return model.Project{}
	}

	return project
}

// Creates a new project
func (d *Database) PostProject(name string) error {
	project := model.Project{}
	project.Name = name
	project.Archived = false

	err := d.DB.Create(&project).Error

	return err
}

// Return an array of all projects
func (d *Database) GetAllProjects() []model.Project {
	projects := []model.Project{}

	d.DB.Find(&projects)

	return projects
}

// Delete a project
func (d *Database) DeleteProject(name string) error {
	project := model.Project{}

	// Unscoped to delete project permanently
	err := d.DB.Unscoped().Where("Name = ?", name).Delete(&project).Error
	return err
}

// Update a project
func (d *Database) UpdateProject(project model.Project) error {
	err := d.DB.Save(&project).Error
	return err
}

// Get project by ID
func (d *Database) GetProjectByID(id uint) model.Project {
	project := model.Project{}
	err := d.DB.Find(&project, id).Error

	if err != nil {
		return model.Project{}
	}

	return project
}

// Get a task
func (d *Database) GetTask(projectName, taskName string) model.Task {
	project := model.Project{}
	err := d.DB.Find(&project, "Name = ?", projectName).Error

	if err != nil {
		return model.Task{}
	}

	task := model.Task{}
	err = d.DB.Find(&task, "Name = ? AND Project_ID = ?", taskName, project.ID).Error

	if err != nil {
		return model.Task{}
	}

	return task
}

// Create a Task
func (d *Database) PostTask(task model.Task) error {
	err := d.DB.Create(&task).Error
	return err
}

// Returns an array of all tasks belonging to a project
func (d *Database) GetAllProjectTasks(project model.Project) []model.Task {
	tasks := []model.Task{}

	d.DB.Find(&tasks, "Project_ID = ?", project.ID)

	return tasks
}

// Deletes a Task
func (d *Database) DeleteTask(task model.Task) error {
	err := d.DB.Unscoped().Where("Name = ? AND Project_ID = ?", task.Name, task.ProjectID).Delete(&task).Error
	return err
}

// Updates a task
func (d *Database) UpdateTask(task model.Task) error {
	err := d.DB.Save(&task).Error
	return err
}

// creates database struct and runs automigrate
func NewDatabaseConnection(name string) *Database {
	db, err := gorm.Open(sqlite.Open(name), &gorm.Config{})

	if err != nil {
		log.Fatalf("Can not open Database %s", err)
	}

	db = model.DbMigrate(db)

	return &Database{DB: db}
}
