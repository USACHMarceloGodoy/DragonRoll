package main

import (
	"dragonroll/backend/config"
	"dragonroll/backend/middleware"
	"dragonroll/backend/repositories"
	"dragonroll/backend/routes"
	"dragonroll/backend/services"
	"dragonroll/backend/utils"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadConfig()

	utils.InitLogger(config.Cfg.Logger.Filepath, config.Cfg.Logger.Filename, config.Cfg.Logger.Level, config.Cfg.Logger.Tz)

	repositories.InitConnections()
	services.InitRepositories()
	if err := services.InitIndexes(); err != nil {
		utils.Fatal("InitIndexes:", err)
	}

	// Bootstrap data: skills universales y admin inicial.
	// Se ejecuta sólo si la base está vacía. Idempotente.
	if err := services.BootstrapSeed(); err != nil {
		utils.Fatal("BootstrapSeed:", err)
	}

	utils.Info("server up en :" + config.Cfg.Server.Port)

	app := gin.Default()
	// CORS DEBE ir antes de las rutas para que aplique a todo (incluido OPTIONS).
	app.Use(middleware.CorsMiddleware(config.Cfg.Cors.CorsUrl))
	app.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "dragonroll-backend"})
	})

	routes.InitRoutes(app)

	if err := http.ListenAndServe(":"+config.Cfg.Server.Port, app); err != nil {
		utils.Fatal("listen:", err)
		os.Exit(1)
	}
}
