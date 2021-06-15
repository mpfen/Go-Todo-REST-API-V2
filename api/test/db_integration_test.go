package api_test

// Integration tests for Database struct that implements TodoStore interface functions:
//
// GetProject(name string) model.Project
// PostProject(name string) error
// GetAllProjects() []model.Project
// DeleteProject(name string) error
// UpdateProject(project model.Project) error
//
// GetTask(projectID string, taskName string) model.Task
// PostTask(task model.Task) error
// DeleteTask(task model.Task) error
// UpdateTask(task model.Task) error

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/mpfen/Go-Todo-REST-API-V2/api/model"
	"github.com/mpfen/Go-Todo-REST-API-V2/api/store"
	"github.com/stretchr/testify/assert"
)

const testdbfile = "testdb.db"

// create a testdb file
func createTestDB(t *testing.T) {
	testDB := []byte("")
	err := ioutil.WriteFile(testdbfile, testDB, 0644)

	if err != nil {
		t.Fatalf("could not create test database %v", err)
	}
}

// delete the testdb file
func deleteTestDB(t *testing.T) {
	err := os.Remove(testdbfile)

	if err != nil {
		t.Fatalf("Could not remove test databse %v", err)
	}
}

// populate test database with projects
func populateTestDatabaseProjects(t *testing.T, db *store.Database) {
	err := db.PostProject("homework")

	if err != nil {
		t.Fatalf("Error populating test database with projects: %v", err)
		return
	}

	err = db.PostProject("cleaning")

	if err != nil {
		t.Fatalf("Error populating test database with projects: %v", err)
	}
}

// populate test database with tasks
func populateTestDatabaseTasks(t *testing.T, db *store.Database) {
	taskBiology := model.Task{Name: "biology", ProjectID: uint(2)}
	err := db.PostTask(taskBiology)

	if err != nil {
		t.Fatalf("Error populating test database with tasks: %v", err)
		return
	}

	taskPhysics := model.Task{Name: "physics", ProjectID: uint(2)}
	err = db.PostTask(taskPhysics)

	if err != nil {
		t.Fatalf("Error populating test database with tasks: %v", err)
	}
}

// Integration tests for database
// uses own database file
func TestDatabase(t *testing.T) {
	createTestDB(t)
	defer deleteTestDB(t)

	db := store.NewDatabaseConnection(testdbfile)

	// PostProject(name string) error
	t.Run("Create a new project in database", func(t *testing.T) {
		want := "TestDatabase"
		err := db.PostProject(want)

		assert.NoError(t, err, "Create new project in db")

		got := db.GetProject(want).Name

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}

	})

	// projects/homework
	// projects/cleaning
	populateTestDatabaseProjects(t, db)

	t.Run("Try to create an already existing project", func(t *testing.T) {
		err := db.PostProject("TestDatabase")

		if err == nil {
			t.Errorf("Project should not have been created")
		}
	})

	// DeleteProject(name string) error
	t.Run("Delete a project", func(t *testing.T) {
		err := db.DeleteProject("TestDatabase")

		assert.NoError(t, err, "Project should have been deleted")
	})

	// GetProject(name string) model.Project
	t.Run("Try to get project TestDatabse", func(t *testing.T) {
		projectName := "homework"
		project := db.GetProject(projectName)

		if project.Name != projectName {
			t.Errorf("Project not found: got '%v' wanted %v", project.Name, projectName)
		}
	})

	t.Run("Try to get non existent project", func(t *testing.T) {
		projectName := "NotTestDatabase"
		project := db.GetProject(projectName)

		if project.Name != "" {
			t.Error("No Project should have been found")
		}
	})

	// GetAllProject() []model.Projects
	t.Run("Get all projects in the database", func(t *testing.T) {
		projects := db.GetAllProjects()

		if i := len(projects); i != 2 {
			t.Errorf("Not the right number of projects found: Found %v wanted 2", len(projects))
		}
	})

	// UpdateProject(project model.project) error
	t.Run("Update the name of project cleaning", func(t *testing.T) {
		project := db.GetProject("cleaning")
		project.Name = "springCleaning"
		err := db.UpdateProject(project)

		assert.NoErrorf(t, err, "update name of project failed: %v", err)

		updatedProject := db.GetProject("springCleaning")

		if updatedProject.Name == "" {
			t.Error("Project was not updated")
		}
	})

	// PostTask(task model.Task) error
	// GetTask(projectID string, taskName string) model.Task
	t.Run("Create a new task math for project homework", func(t *testing.T) {
		taskMath := model.Task{Name: "math", ProjectID: uint(2)}

		err := db.PostTask(taskMath)

		assert.NoError(t, err, "Task creation failed")

		task := db.GetTask("homework", "math")

		if task.Name == "" {
			t.Error("Newly created Task not found")
		}
	})

	// homework/task/biology
	// homework/task/physics
	populateTestDatabaseTasks(t, db)

	// GetAllProjectTasks(projectName string) []model.Task
	t.Run("Get all tasks from project homework", func(t *testing.T) {
		project := db.GetProject("homework")
		tasks := db.GetAllProjectTasks(project)

		if len(tasks) != 3 {
			t.Errorf("Not the right numbers of tasks found: got %v want %v", len(tasks), 3)
		}
	})

	t.Run("Get all tasks from a project without tasks", func(t *testing.T) {
		project := db.GetProject("cleaning")
		tasks := db.GetAllProjectTasks(project)

		if len(tasks) != 0 {
			t.Errorf("Not the right numbers of tasks found: got %v want %v", len(tasks), 3)
		}
	})

	// DeleteTask(task model.Task) error
	t.Run("Delete a task", func(t *testing.T) {
		task := db.GetTask("homework", "math")
		err := db.DeleteTask(task)

		assert.NoError(t, err, "Task should have been deleted")
	})

	// UpdateTask(task model.Task) error
	t.Run("Update a task", func(t *testing.T) {
		// Get task to update and change the name
		task := db.GetTask("homework", "physics")
		task.Name = "newtonsLaw"

		err := db.UpdateTask(task)
		assert.NoError(t, err, "tried to update task")

		// check if task was updated
		wasUpdated := db.GetTask("homework", "newtonsLaw")
		if wasUpdated.Name == "" {
			t.Error("Task was not updated")
		}
	})
}
