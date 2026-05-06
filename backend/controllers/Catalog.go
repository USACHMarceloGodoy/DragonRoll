package controllers

import (
	"dragonroll/backend/models"
	"dragonroll/backend/services"
	"dragonroll/backend/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// pickKind toma el segmento de URL (/classes, /spells, ...) y resuelve
// el CatalogKind. Permite tener un solo set de handlers reutilizables
// montados con paths distintos.
func pickKind(path string) (services.CatalogKind, bool) {
	switch {
	case strings.HasPrefix(path, "/classes"):
		return services.KindClass, true
	case strings.HasPrefix(path, "/subclasses"):
		return services.KindSubclass, true
	case strings.HasPrefix(path, "/spells"):
		return services.KindSpell, true
	case strings.HasPrefix(path, "/features"):
		return services.KindFeature, true
	case strings.HasPrefix(path, "/races"):
		return services.KindRace, true
	case strings.HasPrefix(path, "/skills"):
		return services.KindSkill, true
	}
	return "", false
}

// ListCatalogController devuelve los documentos visibles del catálogo.
// El JSON sale tal cual lo guarda Mongo, así que el frontend ve los
// campos completos del tipo (Class incluye levels, etc.).
func ListCatalogController(c *gin.Context) {
	ac, ok := MustAuthContext(c)
	if !ok {
		return
	}
	kind, ok := pickKind(c.FullPath())
	if !ok {
		utils.JsonResponse(c, 400, "ruta de catálogo desconocida", nil)
		return
	}
	raws, err := services.ListCatalogVisible(ac, kind)
	if err != nil {
		utils.JsonResponse(c, 500, err.Error(), nil)
		return
	}
	out := decodeRaws(raws, kind)
	utils.JsonResponse(c, 200, "ok", out)
}

func GetCatalogController(c *gin.Context) {
	ac, ok := MustAuthContext(c)
	if !ok {
		return
	}
	kind, ok := pickKind(c.FullPath())
	if !ok {
		utils.JsonResponse(c, 400, "ruta de catálogo desconocida", nil)
		return
	}
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.JsonResponse(c, 400, "ID inválido.", nil)
		return
	}
	raw, err := services.GetCatalogByID(ac, kind, id)
	if err != nil {
		utils.JsonResponse(c, 404, err.Error(), nil)
		return
	}
	doc := decodeOne(raw, kind)
	utils.JsonResponse(c, 200, "ok", doc)
}

// CreateCatalogController acepta el JSON tipado del frontend (un Class,
// Spell, etc., con su Meta ya con scope/campaign_id/source_label) e inserta.
// Llenamos OwnerID + timestamps server-side.
func CreateCatalogController(c *gin.Context) {
	ac, ok := MustAuthContext(c)
	if !ok {
		return
	}
	kind, ok := pickKind(c.FullPath())
	if !ok {
		utils.JsonResponse(c, 400, "ruta de catálogo desconocida", nil)
		return
	}
	doc, meta, err := bindCatalogDoc(c, kind, ac.UserID, true)
	if err != nil {
		utils.JsonResponse(c, 400, err.Error(), nil)
		return
	}
	id, err := services.CreateCatalogDoc(ac, kind, doc, meta)
	if err != nil {
		utils.JsonResponse(c, 400, err.Error(), nil)
		return
	}
	utils.JsonResponse(c, 201, "creado", gin.H{"_id": id})
}

func UpdateCatalogController(c *gin.Context) {
	ac, ok := MustAuthContext(c)
	if !ok {
		return
	}
	kind, ok := pickKind(c.FullPath())
	if !ok {
		utils.JsonResponse(c, 400, "ruta de catálogo desconocida", nil)
		return
	}
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.JsonResponse(c, 400, "ID inválido.", nil)
		return
	}
	doc, meta, err := bindCatalogDoc(c, kind, ac.UserID, false)
	if err != nil {
		utils.JsonResponse(c, 400, err.Error(), nil)
		return
	}
	if err := services.ReplaceCatalogDoc(ac, kind, id, doc, meta); err != nil {
		utils.JsonResponse(c, 400, err.Error(), nil)
		return
	}
	utils.JsonResponse(c, 200, "ok", nil)
}

func DeleteCatalogController(c *gin.Context) {
	ac, ok := MustAuthContext(c)
	if !ok {
		return
	}
	kind, ok := pickKind(c.FullPath())
	if !ok {
		utils.JsonResponse(c, 400, "ruta de catálogo desconocida", nil)
		return
	}
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.JsonResponse(c, 400, "ID inválido.", nil)
		return
	}
	if err := services.DeleteCatalogDoc(ac, kind, id); err != nil {
		utils.JsonResponse(c, 400, err.Error(), nil)
		return
	}
	utils.JsonResponse(c, 200, "ok", nil)
}

// decodeRaws → []interface{} con structs concretos según kind.
func decodeRaws(raws []bson.Raw, kind services.CatalogKind) []interface{} {
	out := make([]interface{}, 0, len(raws))
	for _, r := range raws {
		out = append(out, decodeOne(r, kind))
	}
	return out
}

func decodeOne(raw bson.Raw, kind services.CatalogKind) interface{} {
	switch kind {
	case services.KindClass:
		var v models.Class
		_ = bson.Unmarshal(raw, &v)
		return &v
	case services.KindSubclass:
		var v models.Subclass
		_ = bson.Unmarshal(raw, &v)
		return &v
	case services.KindSpell:
		var v models.Spell
		_ = bson.Unmarshal(raw, &v)
		return &v
	case services.KindFeature:
		var v models.Feature
		_ = bson.Unmarshal(raw, &v)
		return &v
	case services.KindRace:
		var v models.Race
		_ = bson.Unmarshal(raw, &v)
		return &v
	case services.KindSkill:
		var v models.Skill
		_ = bson.Unmarshal(raw, &v)
		return &v
	}
	return nil
}

// bindCatalogDoc parsea el body al tipo correcto y devuelve el doc + Meta.
// Si isCreate, fuerza CreatedAt y OwnerID. Siempre setea UpdatedAt.
func bindCatalogDoc(c *gin.Context, kind services.CatalogKind, ownerID primitive.ObjectID, isCreate bool) (interface{}, models.CatalogMeta, error) {
	switch kind {
	case services.KindClass:
		var v models.Class
		if err := c.BindJSON(&v); err != nil {
			return nil, models.CatalogMeta{}, err
		}
		v.Meta = stampMeta(v.Meta, ownerID, isCreate)
		return &v, v.Meta, nil
	case services.KindSubclass:
		var v models.Subclass
		if err := c.BindJSON(&v); err != nil {
			return nil, models.CatalogMeta{}, err
		}
		v.Meta = stampMeta(v.Meta, ownerID, isCreate)
		return &v, v.Meta, nil
	case services.KindSpell:
		var v models.Spell
		if err := c.BindJSON(&v); err != nil {
			return nil, models.CatalogMeta{}, err
		}
		v.Meta = stampMeta(v.Meta, ownerID, isCreate)
		return &v, v.Meta, nil
	case services.KindFeature:
		var v models.Feature
		if err := c.BindJSON(&v); err != nil {
			return nil, models.CatalogMeta{}, err
		}
		v.Meta = stampMeta(v.Meta, ownerID, isCreate)
		return &v, v.Meta, nil
	case services.KindRace:
		var v models.Race
		if err := c.BindJSON(&v); err != nil {
			return nil, models.CatalogMeta{}, err
		}
		v.Meta = stampMeta(v.Meta, ownerID, isCreate)
		return &v, v.Meta, nil
	case services.KindSkill:
		var v models.Skill
		if err := c.BindJSON(&v); err != nil {
			return nil, models.CatalogMeta{}, err
		}
		v.Meta = stampMeta(v.Meta, ownerID, isCreate)
		return &v, v.Meta, nil
	}
	return nil, models.CatalogMeta{}, nil
}

func stampMeta(m models.CatalogMeta, ownerID primitive.ObjectID, isCreate bool) models.CatalogMeta {
	now := timeNow()
	if isCreate {
		m.OwnerID = ownerID
		m.CreatedAt = now
	}
	m.UpdatedAt = now
	if m.Scope == "" {
		m.Scope = models.ScopeUser
	}
	if !m.IsHomebrew {
		m.IsHomebrew = m.Scope != models.ScopeGlobal
	}
	return m
}
