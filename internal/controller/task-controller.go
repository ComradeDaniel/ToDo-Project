package controller

import (
	"log"
	"net/http"
	"todolist/internal/auth"
	"todolist/internal/database"
	"todolist/internal/service"

	"github.com/gin-gonic/gin"
)

type TaskController interface {
	AddTask(ctx *gin.Context)
	UpdateTask(ctx *gin.Context)
	DeleteTask(ctx *gin.Context)
	RelocateTask(ctx *gin.Context)
	AddCategory(ctx *gin.Context)
	UpdateCategory(ctx *gin.Context)
	DeleteCategory(ctx *gin.Context)
	RelocateCategory(ctx *gin.Context)
	GetAllTasksAndCategories(ctx *gin.Context)
}

type taskController struct {
	service service.TaskService
}

func NewTaskController(service service.TaskService) TaskController {
	return &taskController{
		service: service,
	}
}

func (c *taskController) AddTask(ctx *gin.Context) {
	var task database.Task
	err := ctx.BindJSON(&task)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	username, err := auth.GetUsernameFromCtx(ctx)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	task, err = c.service.AddTask(task, username)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, task)
}

func (c *taskController) UpdateTask(ctx *gin.Context) {
	var task database.Task
	err := ctx.BindJSON(&task)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	username, err := auth.GetUsernameFromCtx(ctx)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	task, err = c.service.UpdateTask(task, username)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, task)
}

func (c *taskController) DeleteTask(ctx *gin.Context) {
	var task database.Task
	err := ctx.BindJSON(&task)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	username, err := auth.GetUsernameFromCtx(ctx)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	err = c.service.DeleteTask(task, username)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.Status(http.StatusOK)
}

func (c *taskController) RelocateTask(ctx *gin.Context) {
	var task database.Task
	err := ctx.BindJSON(&task)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	username, err := auth.GetUsernameFromCtx(ctx)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	task, err = c.service.RelocateTask(task, username)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, task)
}

func (c *taskController) AddCategory(ctx *gin.Context) {
	var category database.Categories
	err := ctx.BindJSON(&category)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	username, err := auth.GetUsernameFromCtx(ctx)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	category, err = c.service.AddCategory(category, username)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, category)
}

func (c *taskController) UpdateCategory(ctx *gin.Context) {
	var category database.Categories
	err := ctx.BindJSON(&category)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	username, err := auth.GetUsernameFromCtx(ctx)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	category, err = c.service.UpdateCategory(category, username)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, category)
}

func (c *taskController) DeleteCategory(ctx *gin.Context) {
	var category database.Categories
	err := ctx.BindJSON(&category)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	username, err := auth.GetUsernameFromCtx(ctx)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	err = c.service.DeleteCategory(category, username)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.Status(http.StatusOK)
}

func (c *taskController) RelocateCategory(ctx *gin.Context) {
	var category database.Categories
	err := ctx.BindJSON(&category)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	username, err := auth.GetUsernameFromCtx(ctx)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	err = c.service.RelocateCategory(category, username)
	type getdata struct {
		Categories []database.Categories `json:"categories"`
		Tasks      []database.Task       `json:"tasks"`
	}
	var data getdata
	data.Categories, data.Tasks = c.service.GetAllTasksAndCategories(username)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, data)
}

func (c *taskController) GetAllTasksAndCategories(ctx *gin.Context) {
	username, err := auth.GetUsernameFromCtx(ctx)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	type getdata struct {
		Categories []database.Categories `json:"categories"`
		Tasks      []database.Task       `json:"tasks"`
	}
	var data getdata

	data.Categories, data.Tasks = c.service.GetAllTasksAndCategories(username)

	ctx.JSON(http.StatusOK, data)
}
