package routes

import (
	"dragonroll/backend/controllers"
	"dragonroll/backend/middleware"

	"github.com/gin-gonic/gin"
)

func InitCharacterRoutes(r *gin.RouterGroup) {
	g := r.Group("/characters")
	g.Use(middleware.LoadJWTAuth().MiddlewareFunc())
	{
		g.POST("/", controllers.CreateCharacterController)
		g.POST("/import", controllers.ImportCharacterController)
		g.GET("/", controllers.ListCharactersController)
		g.GET("/:id", controllers.GetCharacterController)
		g.GET("/:id/export", controllers.ExportCharacterController)
		g.PUT("/:id", controllers.UpdateCharacterController)
		g.PUT("/:id/campaign", controllers.AssignCampaignController)
		g.PATCH("/:id/combat", controllers.PatchCombatController)
		g.DELETE("/:id", controllers.DeleteCharacterController)
	}
}
