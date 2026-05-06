package controllers

import (
	"dragonroll/backend/forms"
	"dragonroll/backend/models"
	"dragonroll/backend/services"
	"dragonroll/backend/utils"

	"github.com/gin-gonic/gin"
)

func CreateUserController(c *gin.Context) {
	var form forms.UserForm
	if err := c.BindJSON(&form); err != nil {
		utils.JsonResponse(c, 400, err.Error(), nil)
		return
	}
	user := &models.User{
		Username: form.Username,
		Email:    form.Email,
		Password: form.Password,
		Roles:    form.Roles,
		Avatar:   form.Avatar,
	}
	created, err := services.CreateUserService(user)
	if err != nil {
		utils.Info("Registro fallido: " + err.Error())
		utils.JsonResponse(c, 400, err.Error(), nil)
		return
	}
	utils.Info("Usuario registrado: " + created.Username)
	utils.JsonResponse(c, 201, "Usuario creado.", created)
}

func GetAllUsersController(c *gin.Context) {
	users, err := services.GetAllUsersService()
	if err != nil {
		utils.JsonResponse(c, 500, err.Error(), nil)
		return
	}
	utils.JsonResponse(c, 200, "ok", users)
}

func MeController(c *gin.Context) {
	ac, ok := MustAuthContext(c)
	if !ok {
		return
	}
	user, err := services.GetUserByID(ac.UserID)
	if err != nil {
		utils.JsonResponse(c, 404, err.Error(), nil)
		return
	}
	utils.JsonResponse(c, 200, "ok", user)
}

func SetActiveRoleController(c *gin.Context) {
	ac, ok := MustAuthContext(c)
	if !ok {
		return
	}
	var form forms.SetActiveRoleForm
	if err := c.BindJSON(&form); err != nil {
		utils.JsonResponse(c, 400, err.Error(), nil)
		return
	}
	if err := services.SetActiveRoleService(ac.UserID, form.ActiveRol); err != nil {
		utils.JsonResponse(c, 400, err.Error(), nil)
		return
	}
	utils.JsonResponse(c, 200, "Rol activo actualizado.", form.ActiveRol)
}

// UpdateProfileController permite editar username/email/avatar del usuario logueado.
func UpdateProfileController(c *gin.Context) {
	ac, ok := MustAuthContext(c)
	if !ok {
		return
	}
	var form forms.UpdateProfileForm
	if err := c.BindJSON(&form); err != nil {
		utils.JsonResponse(c, 400, err.Error(), nil)
		return
	}
	user, err := services.UpdateProfileService(ac.UserID, form.Username, form.Email, form.Avatar)
	if err != nil {
		utils.JsonResponse(c, 400, err.Error(), nil)
		return
	}
	utils.JsonResponse(c, 200, "Perfil actualizado.", user)
}

// ChangePasswordController exige la contraseña actual antes de cambiarla.
func ChangePasswordController(c *gin.Context) {
	ac, ok := MustAuthContext(c)
	if !ok {
		return
	}
	var form forms.ChangePasswordForm
	if err := c.BindJSON(&form); err != nil {
		utils.JsonResponse(c, 400, err.Error(), nil)
		return
	}
	if err := services.ChangePasswordService(ac.UserID, form.CurrentPassword, form.NewPassword); err != nil {
		utils.JsonResponse(c, 400, err.Error(), nil)
		return
	}
	utils.JsonResponse(c, 200, "Contraseña actualizada.", nil)
}
