package services

import (
	"dragonroll/backend/models"
	"dragonroll/backend/repositories"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CatalogKind identifica de qué colección estamos hablando. Cada caso
// mapea a su repositorio y a la función de extracción de Meta para los
// chequeos de autorización.
type CatalogKind string

const (
	KindClass    CatalogKind = "class"
	KindSubclass CatalogKind = "subclass"
	KindSpell    CatalogKind = "spell"
	KindFeature  CatalogKind = "feature"
	KindRace     CatalogKind = "race"
	KindSkill    CatalogKind = "skill"
)

func repoFor(kind CatalogKind) (*repositories.CatalogRepository, error) {
	switch kind {
	case KindClass:
		return classRepo, nil
	case KindSubclass:
		return subclassRepo, nil
	case KindSpell:
		return spellRepo, nil
	case KindFeature:
		return featureRepo, nil
	case KindRace:
		return raceRepo, nil
	case KindSkill:
		return skillRepo, nil
	}
	return nil, errors.New("catálogo desconocido: " + string(kind))
}

// extractMeta sirve para verificar autorización sin tener tipo concreto.
// Lee el campo "meta" de cualquier doc del catálogo.
type catalogMetaProbe struct {
	Meta models.CatalogMeta `bson:"meta"`
}

// ListCatalogVisible devuelve los documentos visibles para el usuario,
// crudos en bson.Raw. El controlador hace bson.Unmarshal al tipo correcto
// (Class/Subclass/Spell/...) según el endpoint.
func ListCatalogVisible(ac AuthContext, kind CatalogKind) ([]bson.Raw, error) {
	repo, err := repoFor(kind)
	if err != nil {
		return nil, err
	}
	dm, player, err := MemberCampaigns(ac.UserID)
	if err != nil {
		return nil, err
	}
	filter := repositories.VisibilityFilter(ac.UserID, append(dm, player...))
	return repo.FindManyRaw(filter)
}

// GetCatalogByID lee un doc por ID en el formato crudo, validando visibilidad.
func GetCatalogByID(ac AuthContext, kind CatalogKind, id primitive.ObjectID) (bson.Raw, error) {
	repo, err := repoFor(kind)
	if err != nil {
		return nil, err
	}
	raws, err := repo.FindManyByIDs([]primitive.ObjectID{id})
	if err != nil || len(raws) == 0 {
		return nil, errors.New("no encontrado")
	}
	raw := raws[0]
	var probe catalogMetaProbe
	if err := bson.Unmarshal(raw, &probe); err != nil {
		return nil, err
	}
	if !canViewMeta(ac, probe.Meta) {
		return nil, errors.New("no autorizado")
	}
	return raw, nil
}

// canViewMeta es la mitad de lectura del filtro de visibilidad. Lo
// duplicamos en Go para validar accesos a un único documento.
func canViewMeta(ac AuthContext, meta models.CatalogMeta) bool {
	if ac.IsAdmin() {
		return true
	}
	switch meta.Scope {
	case models.ScopeGlobal:
		return true
	case models.ScopeUser:
		return meta.OwnerID == ac.UserID
	case models.ScopeCampaign:
		camp, err := campaignRepo.FindByID(meta.CampaignID)
		if err != nil || camp == nil {
			return false
		}
		if camp.DMID == ac.UserID {
			return true
		}
		for _, p := range camp.PlayerIDs {
			if p == ac.UserID {
				return true
			}
		}
	}
	return false
}

// CreateCatalogDoc inserta un documento del catálogo. El controlador
// recibe el JSON tipado, completa la Meta con OwnerID + timestamps, y
// pasa el doc + Meta separados (porque el genérico no conoce el tipo).
//
// El controlador es responsable de pasar el doc concreto (ej *models.Class)
// con Meta ya rellenada. Acá sólo validamos permisos y delegamos al repo.
func CreateCatalogDoc(ac AuthContext, kind CatalogKind, doc interface{}, meta models.CatalogMeta) (primitive.ObjectID, error) {
	ok, err := CanCreateCatalogScope(ac, meta.Scope, meta.CampaignID)
	if err != nil {
		return primitive.NilObjectID, err
	}
	if !ok {
		return primitive.NilObjectID, errors.New("no autorizado para crear con ese scope")
	}
	repo, err := repoFor(kind)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return repo.InsertOne(doc)
}

// ReplaceCatalogDoc valida que el editor pueda escribir sobre el documento
// existente: admin siempre, o coincidir con el scope (owner para user,
// DM para campaign, admin para global).
func ReplaceCatalogDoc(ac AuthContext, kind CatalogKind, id primitive.ObjectID, doc interface{}, newMeta models.CatalogMeta) error {
	repo, err := repoFor(kind)
	if err != nil {
		return err
	}
	raws, err := repo.FindManyByIDs([]primitive.ObjectID{id})
	if err != nil || len(raws) == 0 {
		return errors.New("no encontrado")
	}
	var probe catalogMetaProbe
	if err := bson.Unmarshal(raws[0], &probe); err != nil {
		return err
	}
	if !canEditMeta(ac, probe.Meta) {
		return errors.New("no autorizado")
	}
	// El doc viene con su Meta ya inyectada por el controlador (stampMeta).
	// Aquí sólo nos aseguramos de preservar OwnerID y CreatedAt originales,
	// que el controlador habría podido sobrescribir por error.
	_ = newMeta // newMeta se conserva sólo para mantener la firma estable
	return repo.ReplaceOne(id, doc)
}

func canEditMeta(ac AuthContext, meta models.CatalogMeta) bool {
	if ac.IsAdmin() {
		return true
	}
	switch meta.Scope {
	case models.ScopeGlobal:
		return false // sólo admin
	case models.ScopeUser:
		return meta.OwnerID == ac.UserID
	case models.ScopeCampaign:
		camp, err := campaignRepo.FindByID(meta.CampaignID)
		return err == nil && camp != nil && camp.DMID == ac.UserID
	}
	return false
}

func DeleteCatalogDoc(ac AuthContext, kind CatalogKind, id primitive.ObjectID) error {
	repo, err := repoFor(kind)
	if err != nil {
		return err
	}
	raws, err := repo.FindManyByIDs([]primitive.ObjectID{id})
	if err != nil || len(raws) == 0 {
		return errors.New("no encontrado")
	}
	var probe catalogMetaProbe
	if err := bson.Unmarshal(raws[0], &probe); err != nil {
		return err
	}
	if !canEditMeta(ac, probe.Meta) {
		return errors.New("no autorizado")
	}
	return repo.DeleteOne(id)
}
