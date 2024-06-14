package service

import (
	"errors"
	"log"
	"todolist/internal/database"
)

type TaskService interface {
	AddTask(database.Task, string) (database.Task, error)
	UpdateTask(database.Task, string) (database.Task, error)
	DeleteTask(database.Task, string) error
	RelocateTask(database.Task, string) (database.Task, error)
	AddCategory(database.Categories, string) (database.Categories, error)
	UpdateCategory(database.Categories, string) (database.Categories, error)
	DeleteCategory(database.Categories, string) error
	RelocateCategory(database.Categories, string) ([]database.Categories, error)
	GetAllTasksAndCategories(string) ([]database.Categories, []database.Task)
	checkPermissionTask(int64, int64, string) bool
	checkPermissionCategory(int64, string) bool
}

var (
	ErrForbidden error = errors.New("this user is not permitted to modify this entity")
)

type taskService struct {
}

func NewTaskService() TaskService {
	return &taskService{}
}

func (t *taskService) AddTask(task database.Task, username string) (database.Task, error) {

	if !t.checkPermissionCategory(task.Belongs_to, username) {
		return database.Task{}, ErrForbidden
	}

	task = database.AddTask(task)

	return task, nil
}

func (t *taskService) UpdateTask(task database.Task, username string) (database.Task, error) {

	if !t.checkPermissionTask(task.Belongs_to, task.Id, username) {
		return database.Task{}, ErrForbidden
	}
	newTask := database.UpdateTask(task)

	return newTask, nil
}

func (t *taskService) DeleteTask(task database.Task, username string) error {

	if !t.checkPermissionTask(task.Belongs_to, task.Id, username) {
		return ErrForbidden
	}
	database.DeleteTask(task)

	return nil
}

// This function just deletes the old task and creates a new one at the right place. It returns the new task and an error
func (t *taskService) RelocateTask(task database.Task, username string) (database.Task, error) {

	if !t.checkPermissionTask(task.Belongs_to, task.Id, username) {
		return database.Task{}, ErrForbidden
	}

	database.DeleteTask(task)
	newTask := database.AddTask(task)

	return newTask, nil
}

func (t *taskService) AddCategory(category database.Categories, username string) (database.Categories, error) {
	user, err := database.GetUserByUsername(username)
	if err != nil {
		if errors.Is(err, database.ErrNoResult) {
			return database.Categories{}, ErrForbidden
		}
	}
	category.Belongs_to = user.Id
	category.Id = database.AddCategory(category)
	return category, nil
}

func (t *taskService) UpdateCategory(category database.Categories, username string) (database.Categories, error) {

	if !t.checkPermissionCategory(category.Id, username) {
		return database.Categories{}, ErrForbidden
	}

	return database.UpdateCategory(category), nil
}

func (t *taskService) DeleteCategory(category database.Categories, username string) error {
	if !t.checkPermissionCategory(category.Id, username) {
		return ErrForbidden
	}
	database.DeleteCategory(category)
	return nil
}

func (t *taskService) RelocateCategory(category database.Categories, username string) ([]database.Categories, error) {
	if !t.checkPermissionCategory(category.Id, username) {
		return nil, ErrForbidden
	}

	categories := database.ChangeCategoryOrder(category.Id, category.Order)
	return categories, nil
}

func (t *taskService) GetAllTasksAndCategories(username string) ([]database.Categories, []database.Task) {
	return database.GetCategoriesByUsername(username), database.GetTasksByUsername(username)
}

// These two functions check if the username encoded in the jwt is related to the ids of the entities that are changed to prevent a user from somehow modifying foreign entities
func (t *taskService) checkPermissionTask(belongs_to int64, task_id int64, username string) bool {
	dbUsername, err := database.GetUsernameByCategoryId(belongs_to)
	if err != nil {
		if errors.Is(err, database.ErrNoResult) {
			log.Printf("A permission to modify an entity was denied: %v (db) != %v (request)\n", dbUsername, username)
			return false
		}
		log.Fatalf("Something went wrong while checking permissions of task manipulation: %v", err)
		return false
	}
	if dbUsername == username {
		dbUsername, err = database.GetUsernameByTaskId(task_id)
		if err != nil {
			if errors.Is(err, database.ErrNoResult) {
				return false
			}
			log.Fatalf("Something went wrong while checking permissions of task manipulation: %v", err)
			return false
		}
		if dbUsername == username {
			return true
		} else {
			log.Printf("A permission to modify an entity was denied: %v (db) != %v (request)\n", dbUsername, username)
			return false
		}
	} else {
		return false
	}
}

func (t *taskService) checkPermissionCategory(category_id int64, username string) bool {
	dbUsername, err := database.GetUsernameByCategoryId(category_id)
	if err != nil {
		if errors.Is(err, database.ErrNoResult) {
			log.Printf("A permission to modify an entity was denied: %v (db) != %v (request)\n", dbUsername, username)
			return false
		}
		log.Fatalf("Something went wrong while checking permissions of task manipulation: %v", err)
		return false
	}
	if dbUsername == username {
		return true
	} else {
		log.Printf("A permission to modify an entity was denied: %v (db) != %v (request)\n", dbUsername, username)
		return false
	}
}
