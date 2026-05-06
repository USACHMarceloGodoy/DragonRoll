package controllers

import (
	"dragonroll/backend/forms"
	"dragonroll/backend/models"
	"dragonroll/backend/services"
	"dragonroll/backend/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func toObjectIDList(ids []string) []primitive.ObjectID {
	var out []primitive.ObjectID
	for _, s := range ids {
		if oid, err := primitive.ObjectIDFromHex(s); err == nil {
			out = append(out, oid)
		}
	}
	return out
}

func formToCharacter(f forms.CharacterForm) *models.Character {
	c := &models.Character{
		Name:        f.Name,
		Pronouns:    f.Pronouns,
		Background:  f.Background,
		Alignment:   f.Alignment,
		Level:       f.Level,
		XP:          f.XP,
		NextXP:      f.NextXP,
		Pic:         f.Pic,
		Color:       f.Color,
		HP:          f.HP,
		HPMax:       f.HPMax,
		TempHP:      f.TempHP,
		AC:          f.AC,
		Speed:       f.Speed,
		Initiative:  f.Initiative,
		Inspiration: f.Inspiration,
		HitDice:     f.HitDice,
		DeathSaves:  f.DeathSaves,
		Stats:       f.Stats,
		SlotsUsed:   f.SlotsUsed,
		Resources:   f.Resources,
		Attacks:     f.Attacks,
		Conditions:  f.Conditions,
		Inventory:   f.Inventory,
		Personality: f.Personality,
		Ideals:      f.Ideals,
		Bonds:       f.Bonds,
		Flaws:       f.Flaws,
	}
	if oid, err := primitive.ObjectIDFromHex(f.ClassID); err == nil {
		c.ClassID = oid
	}
	if oid, err := primitive.ObjectIDFromHex(f.SubclassID); err == nil {
		c.SubclassID = oid
	}
	if oid, err := primitive.ObjectIDFromHex(f.RaceID); err == nil {
		c.RaceID = oid
	}
	if oid, err := primitive.ObjectIDFromHex(f.CampaignID); err == nil {
		c.CampaignID = oid
	}
	c.SkillProficiencyIDs = toObjectIDList(f.SkillProficiencyIDs)
	c.SkillExpertiseIDs = toObjectIDList(f.SkillExpertiseIDs)
	c.CantripIDs = toObjectIDList(f.CantripIDs)
	c.SpellsKnownIDs = toObjectIDList(f.SpellsKnownIDs)
	c.SpellsPreparedIDs = toObjectIDList(f.SpellsPreparedIDs)
	return c
}

func CreateCharacterController(c *gin.Context) {
	ac, ok := MustAuthContext(c)
	if !ok {
		return
	}
	var f forms.CharacterForm
	if err := c.BindJSON(&f); err != nil {
		utils.JsonResponse(c, 400, err.Error(), nil)
		return
	}
	character := formToCharacter(f)
	created, err := services.CreateCharacterService(ac, character)
	if err != nil {
		utils.JsonResponse(c, 400, err.Error(), nil)
		return
	}
	utils.JsonResponse(c, 201, "Personaje creado.", created)
}

func ListCharactersController(c *gin.Context) {
	ac, ok := MustAuthContext(c)
	if !ok {
		return
	}
	chars, err := services.ListAccessibleCharacters(ac)
	if err != nil {
		utils.JsonResponse(c, 500, err.Error(), nil)
		return
	}
	utils.JsonResponse(c, 200, "ok", chars)
}

// GetCharacterController por defecto trae expanded=true (con catálogo poblado).
// El cliente puede pedir ?expand=false para sólo el documento crudo.
func GetCharacterController(c *gin.Context) {
	ac, ok := MustAuthContext(c)
	if !ok {
		return
	}
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.JsonResponse(c, 400, "ID inválido.", nil)
		return
	}
	if c.Query("expand") == "false" {
		ch, err := services.GetCharacterByIDAuthorized(ac, id)
		if err != nil {
			utils.JsonResponse(c, 404, err.Error(), nil)
			return
		}
		utils.JsonResponse(c, 200, "ok", ch)
		return
	}
	exp, err := services.GetCharacterExpanded(ac, id)
	if err != nil {
		utils.JsonResponse(c, 404, err.Error(), nil)
		return
	}
	utils.JsonResponse(c, 200, "ok", exp)
}

func UpdateCharacterController(c *gin.Context) {
	ac, ok := MustAuthContext(c)
	if !ok {
		return
	}
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.JsonResponse(c, 400, "ID inválido.", nil)
		return
	}
	var f forms.CharacterForm
	if err := c.BindJSON(&f); err != nil {
		utils.JsonResponse(c, 400, err.Error(), nil)
		return
	}
	updated := formToCharacter(f)
	saved, err := services.UpdateCharacterService(ac, id, updated)
	if err != nil {
		utils.JsonResponse(c, 400, err.Error(), nil)
		return
	}
	utils.JsonResponse(c, 200, "Personaje actualizado.", saved)
}

// PatchCombatController: cambios pequeños desde la pestaña de combate.
func PatchCombatController(c *gin.Context) {
	ac, ok := MustAuthContext(c)
	if !ok {
		return
	}
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.JsonResponse(c, 400, "ID inválido.", nil)
		return
	}
	var f forms.CombatPatchForm
	if err := c.BindJSON(&f); err != nil {
		utils.JsonResponse(c, 400, err.Error(), nil)
		return
	}
	fields := bson.M{}
	if f.HP != nil {
		fields["hp"] = *f.HP
	}
	if f.TempHP != nil {
		fields["temp_hp"] = *f.TempHP
	}
	if f.Conditions != nil {
		fields["conditions"] = f.Conditions
	}
	if f.SlotsUsed != nil {
		fields["slots_used"] = *f.SlotsUsed
	}
	if len(fields) == 0 {
		utils.JsonResponse(c, 400, "Sin campos a actualizar.", nil)
		return
	}
	if err := services.PatchCombat(ac, id, fields); err != nil {
		utils.JsonResponse(c, 400, err.Error(), nil)
		return
	}
	utils.JsonResponse(c, 200, "ok", fields)
}

func DeleteCharacterController(c *gin.Context) {
	ac, ok := MustAuthContext(c)
	if !ok {
		return
	}
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.JsonResponse(c, 400, "ID inválido.", nil)
		return
	}
	if err := services.DeleteCharacterService(ac, id); err != nil {
		utils.JsonResponse(c, 400, err.Error(), nil)
		return
	}
	utils.JsonResponse(c, 200, "ok", nil)
}

// AssignCampaignController asocia un PJ a una campaña (o lo desasocia si
// se envía campaign_id vacío). Solo el dueño puede.
func AssignCampaignController(c *gin.Context) {
	ac, ok := MustAuthContext(c)
	if !ok {
		return
	}
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.JsonResponse(c, 400, "ID inválido.", nil)
		return
	}
	var f forms.AssignCampaignForm
	if err := c.BindJSON(&f); err != nil {
		utils.JsonResponse(c, 400, err.Error(), nil)
		return
	}
	var campID primitive.ObjectID
	if f.CampaignID != "" {
		campID, err = primitive.ObjectIDFromHex(f.CampaignID)
		if err != nil {
			utils.JsonResponse(c, 400, "campaign_id inválido", nil)
			return
		}
	}
	if err := services.AssignCharacterToCampaign(ac, id, campID); err != nil {
		utils.JsonResponse(c, 400, err.Error(), nil)
		return
	}
	utils.JsonResponse(c, 200, "Asociación actualizada.", nil)
}

// ExportCharacterController devuelve el JSON portable del PJ.
func ExportCharacterController(c *gin.Context) {
	ac, ok := MustAuthContext(c)
	if !ok {
		return
	}
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.JsonResponse(c, 400, "ID inválido.", nil)
		return
	}
	portable, err := services.ExportCharacterPortable(ac, id)
	if err != nil {
		utils.JsonResponse(c, 404, err.Error(), nil)
		return
	}
	utils.JsonResponse(c, 200, "ok", portable)
}

// ImportCharacterController acepta un JSON portable y crea un PJ nuevo
// para el usuario logueado.
func ImportCharacterController(c *gin.Context) {
	ac, ok := MustAuthContext(c)
	if !ok {
		return
	}
	var portable services.CharacterPortable
	if err := c.BindJSON(&portable); err != nil {
		utils.JsonResponse(c, 400, err.Error(), nil)
		return
	}
	created, err := services.ImportCharacterPortable(ac, &portable)
	if err != nil {
		utils.JsonResponse(c, 400, err.Error(), nil)
		return
	}
	utils.JsonResponse(c, 201, "Personaje importado.", created)
}
