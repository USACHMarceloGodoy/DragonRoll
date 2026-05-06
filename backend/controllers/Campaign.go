package controllers

import (
	"dragonroll/backend/forms"
	"dragonroll/backend/models"
	"dragonroll/backend/services"
	"dragonroll/backend/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateCampaignController(c *gin.Context) {
	ac, ok := MustAuthContext(c)
	if !ok {
		return
	}
	var f forms.CampaignForm
	if err := c.BindJSON(&f); err != nil {
		utils.JsonResponse(c, 400, err.Error(), nil)
		return
	}
	camp := &models.Campaign{
		Name:        f.Name,
		Description: f.Description,
		System:      f.System,
	}
	created, err := services.CreateCampaignService(ac, camp)
	if err != nil {
		utils.JsonResponse(c, 400, err.Error(), nil)
		return
	}
	utils.JsonResponse(c, 201, "Campaña creada.", created)
}

func ListCampaignsController(c *gin.Context) {
	ac, ok := MustAuthContext(c)
	if !ok {
		return
	}
	camps, err := services.ListMyCampaigns(ac)
	if err != nil {
		utils.JsonResponse(c, 500, err.Error(), nil)
		return
	}
	utils.JsonResponse(c, 200, "ok", camps)
}

func GetCampaignController(c *gin.Context) {
	ac, ok := MustAuthContext(c)
	if !ok {
		return
	}
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.JsonResponse(c, 400, "ID inválido.", nil)
		return
	}
	camp, err := services.GetCampaignAuthorized(ac, id)
	if err != nil {
		utils.JsonResponse(c, 404, err.Error(), nil)
		return
	}
	utils.JsonResponse(c, 200, "ok", camp)
}

func UpdateCampaignController(c *gin.Context) {
	ac, ok := MustAuthContext(c)
	if !ok {
		return
	}
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.JsonResponse(c, 400, "ID inválido.", nil)
		return
	}
	var camp models.Campaign
	if err := c.BindJSON(&camp); err != nil {
		utils.JsonResponse(c, 400, err.Error(), nil)
		return
	}
	saved, err := services.UpdateCampaignService(ac, id, &camp)
	if err != nil {
		utils.JsonResponse(c, 400, err.Error(), nil)
		return
	}
	utils.JsonResponse(c, 200, "ok", saved)
}

func AddPlayerToCampaignController(c *gin.Context) {
	ac, ok := MustAuthContext(c)
	if !ok {
		return
	}
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.JsonResponse(c, 400, "ID inválido.", nil)
		return
	}
	var f forms.AddPlayerForm
	if err := c.BindJSON(&f); err != nil {
		utils.JsonResponse(c, 400, err.Error(), nil)
		return
	}
	player, err := services.GetUserByLogin(f.UsernameOrEmail)
	if err != nil {
		utils.JsonResponse(c, 404, "Usuario a invitar no encontrado.", nil)
		return
	}
	if err := services.AddPlayerToCampaign(ac, id, player.ID); err != nil {
		utils.JsonResponse(c, 400, err.Error(), nil)
		return
	}
	utils.JsonResponse(c, 200, "Jugador agregado.", player.Username)
}

func DeleteCampaignController(c *gin.Context) {
	ac, ok := MustAuthContext(c)
	if !ok {
		return
	}
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.JsonResponse(c, 400, "ID inválido.", nil)
		return
	}
	if err := services.DeleteCampaignService(ac, id); err != nil {
		utils.JsonResponse(c, 400, err.Error(), nil)
		return
	}
	utils.JsonResponse(c, 200, "ok", nil)
}
