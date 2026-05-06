package controllers

import (
	"dragonroll/backend/models"
	"dragonroll/backend/services"
	"dragonroll/backend/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MustAuthContext lee el AuthContext desde los claims del JWT que el
// middleware ya dejó en el contexto de Gin (`user_id`, `active_rol`).
// Si algo falla, responde 401 y devuelve `false` para que el handler
// haga return temprano.
func MustAuthContext(c *gin.Context) (services.AuthContext, bool) {
	idAny, ok := c.Get("user_id")
	if !ok {
		utils.JsonResponse(c, 401, "No autorizado.", nil)
		return services.AuthContext{}, false
	}
	id, ok := idAny.(primitive.ObjectID)
	if !ok {
		utils.JsonResponse(c, 500, "Claims malformados.", nil)
		return services.AuthContext{}, false
	}

	user, err := services.GetUserByID(id)
	if err != nil {
		utils.JsonResponse(c, 401, "Usuario no encontrado.", nil)
		return services.AuthContext{}, false
	}

	activeRolStr, _ := c.Get("active_rol")
	activeRol, _ := activeRolStr.(string)

	return services.AuthContext{
		UserID:    user.ID,
		ActiveRol: models.Rol(activeRol),
		Roles:     user.Roles,
	}, true
}
