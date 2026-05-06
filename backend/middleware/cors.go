package middleware

import (
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CorsMiddleware: requerimos AllowCredentials=true porque la auth viaja
// en cookie HttpOnly emitida por gin-jwt. Eso obliga a que AllowOrigins
// sea explícito (sin "*"), por eso leemos CORS_URLS del entorno.
func CorsMiddleware(corsUrl string) gin.HandlerFunc {
	configCors := cors.DefaultConfig()
	configCors.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	configCors.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	configCors.AllowOrigins = strings.Split(corsUrl, ",")
	configCors.AllowCredentials = true
	return cors.New(configCors)
}
