package routes

import "github.com/gin-gonic/gin"

func InitRoutes(r *gin.Engine) {
	g := r.Group("/api/v1")
	InitAuthRoutes(g)
	InitUserRoutes(g)
	InitCharacterRoutes(g)
	InitCampaignRoutes(g)
	InitCatalogRoutes(g)
}
