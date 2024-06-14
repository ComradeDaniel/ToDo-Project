package controller

import (
	"errors"
	"fmt"
	"net/http"
	"todolist/internal/database"
	"todolist/internal/service"

	"github.com/gin-gonic/gin"
)

type UserController interface {
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
}

type userController struct {
	service service.UserService
}

func NewUserController(service service.UserService) UserController {
	return userController{
		service: service,
	}
}

func (c userController) Register(ctx *gin.Context) {
	var user database.User
	var token string
	err := ctx.BindJSON(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	token, err = c.service.RegisterUser(user)

	if err != nil {
		if errors.Is(err, service.ErrUserAlreadyExists) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "A user with that username already exists",
			})
			return
		}
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"jwt": token,
	})
}

func (c userController) Login(ctx *gin.Context) {
	var user database.User
	err := ctx.BindJSON(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	var token string
	token, err = c.service.LoginUser(user)
	if err != nil {
		if errors.Is(err, service.ErrNoSuchUser) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		if errors.Is(err, service.ErrWrongPassword) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"jwt": token,
	})
}
