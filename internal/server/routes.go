package server

import (
	"net/http"
	"todolist/internal/auth"
	"todolist/internal/controller"
	"todolist/internal/service"

	"github.com/gin-gonic/gin"
)

var (
	userService    service.UserService       = service.NewUserService()
	userController controller.UserController = controller.NewUserController(userService)
	taskService    service.TaskService       = service.NewTaskService()
	taskController controller.TaskController = controller.NewTaskController(taskService)
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.LoadHTMLFiles("internal/frontend/login.html", "internal/frontend/index.html")
	r.Static("/static", "internal/frontend/static")
	r.Static("/tasks/static", "internal/frontend/static")

	r.GET("/login", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "login.html", gin.H{})
	})

	r.POST("/login", userController.Login)
	r.POST("/register", userController.Register)

	r.GET("/health", s.healthHandler)

	authorized := r.Group("/tasks")
	authorized.Use(auth.JwtTokenCheck)
	authorized.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", gin.H{})
	})
	authorized.GET("/get", taskController.GetAllTasksAndCategories)

	authorized.POST("/addTask", taskController.AddTask)
	authorized.POST("/deleteTask", taskController.DeleteTask)
	authorized.POST("/updateTask", taskController.UpdateTask)
	authorized.POST("/relocateTask", taskController.RelocateTask)

	authorized.POST("/addCategory", taskController.AddCategory)
	authorized.POST("/updateCategory", taskController.UpdateCategory)
	authorized.POST("/deleteCategory", taskController.DeleteCategory)
	authorized.POST("/relocateCategory", taskController.RelocateCategory)

	return r
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}
