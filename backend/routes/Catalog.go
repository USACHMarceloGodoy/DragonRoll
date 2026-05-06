package routes

import (
	"dragonroll/backend/controllers"
	"dragonroll/backend/middleware"

	"github.com/gin-gonic/gin"
)

// InitCatalogRoutes monta los 6 prefijos del catálogo. Los handlers
// detectan a cuál corresponde mediante c.FullPath().
func InitCatalogRoutes(r *gin.RouterGroup) {
	prefixes := []string{"/classes", "/subclasses", "/spells", "/features", "/races", "/skills"}
	auth := middleware.LoadJWTAuth().MiddlewareFunc()
	for _, p := range prefixes {
		g := r.Group(p)
		g.Use(auth)
		g.GET("", controllers.ListCatalogController)
		g.POST("", controllers.CreateCatalogController)
		g.GET("/:id", controllers.GetCatalogController)
		g.PUT("/:id", controllers.UpdateCatalogController)
		g.DELETE("/:id", controllers.DeleteCatalogController)
	}
}
