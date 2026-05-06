package routes

import (
	"dragonroll/backend/controllers"
	"dragonroll/backend/middleware"
	"dragonroll/backend/models"

	"github.com/gin-gonic/gin"
)

func InitUserRoutes(r *gin.RouterGroup) {
	g := r.Group("/user")
	auth := middleware.LoadJWTAuth().MiddlewareFunc()
	{
		g.POST("/", controllers.CreateUserController) // registro público
		g.GET("/me", auth, controllers.MeController)
		g.PUT("/me", auth, controllers.UpdateProfileController)
		g.PUT("/me/password", auth, controllers.ChangePasswordController)
		g.PUT("/me/active-role", auth, controllers.SetActiveRoleController)
		g.GET("/", auth, middleware.SetRoles(models.ADMIN), controllers.GetAllUsersController)
	}
}
