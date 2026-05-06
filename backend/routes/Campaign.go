package routes

import (
	"dragonroll/backend/controllers"
	"dragonroll/backend/middleware"

	"github.com/gin-gonic/gin"
)

func InitCampaignRoutes(r *gin.RouterGroup) {
	g := r.Group("/campaigns")
	g.Use(middleware.LoadJWTAuth().MiddlewareFunc())
	{
		g.POST("/", controllers.CreateCampaignController)
		g.GET("/", controllers.ListCampaignsController)
		g.GET("/:id", controllers.GetCampaignController)
		g.PUT("/:id", controllers.UpdateCampaignController)
		g.POST("/:id/players", controllers.AddPlayerToCampaignController)
		g.DELETE("/:id", controllers.DeleteCampaignController)
	}
}
