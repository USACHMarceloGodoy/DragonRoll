package routes

import (
	"dragonroll/backend/controllers"
	"dragonroll/backend/middleware"

	"github.com/gin-gonic/gin"
)

func InitAuthRoutes(r *gin.RouterGroup) {
	g := r.Group("/auth")
	{
		g.POST("/login", controllers.LoginFunc)
		g.POST("/refresh", middleware.LoadJWTAuth().MiddlewareFunc(), controllers.RefreshToken)
		g.POST("/logout", middleware.LoadJWTAuth().MiddlewareFunc(), controllers.Logout)
	}
}
