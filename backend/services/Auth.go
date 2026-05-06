package services

import (
	"dragonroll/backend/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AuthContext es el snapshot de quién está pidiendo. Lo construye el
// controlador a partir de los claims del JWT y se pasa hacia abajo.
type AuthContext struct {
	UserID    primitive.ObjectID
	ActiveRol models.Rol
	Roles     []models.Rol
}

func (ac AuthContext) IsAdmin() bool {
	for _, r := range ac.Roles {
		if r == models.ADMIN {
			return true
		}
	}
	return false
}

// CanEditCharacter resuelve la política de edición:
//   - admin edita cualquiera
//   - el dueño edita su propio PJ
//   - el DM de la campaña del PJ edita ese PJ ("controlador total")
func CanEditCharacter(ac AuthContext, character *models.Character) (bool, error) {
	if ac.IsAdmin() {
		return true, nil
	}
	if character.OwnerID == ac.UserID {
		return true, nil
	}
	if !character.CampaignID.IsZero() {
		camp, err := campaignRepo.FindByID(character.CampaignID)
		if err != nil {
			return false, err
		}
		if camp != nil && camp.DMID == ac.UserID {
			return true, nil
		}
	}
	return false, nil
}

// CanViewCharacter es más laxo:
//   - admin
//   - dueño
//   - DM o jugador de la campaña del PJ
func CanViewCharacter(ac AuthContext, character *models.Character) (bool, error) {
	if ac.IsAdmin() || character.OwnerID == ac.UserID {
		return true, nil
	}
	if !character.CampaignID.IsZero() {
		camp, err := campaignRepo.FindByID(character.CampaignID)
		if err != nil {
			return false, err
		}
		if camp != nil {
			if camp.DMID == ac.UserID {
				return true, nil
			}
			for _, p := range camp.PlayerIDs {
				if p == ac.UserID {
					return true, nil
				}
			}
		}
	}
	return false, nil
}

// MemberCampaigns devuelve los IDs de campañas donde el usuario es DM o jugador.
// Se usa para construir el VisibilityFilter del catálogo.
func MemberCampaigns(userID primitive.ObjectID) ([]primitive.ObjectID, []primitive.ObjectID, error) {
	camps, err := campaignRepo.FindByMember(userID)
	if err != nil {
		return nil, nil, err
	}
	var dm, player []primitive.ObjectID
	for _, c := range camps {
		if c.DMID == userID {
			dm = append(dm, c.ID)
			continue
		}
		for _, p := range c.PlayerIDs {
			if p == userID {
				player = append(player, c.ID)
				break
			}
		}
	}
	return dm, player, nil
}

// CanCreateCatalogScope decide si AuthContext puede crear un objeto de
// catálogo con el scope solicitado.
//   - global: sólo admin
//   - campaign: el DM de esa campaña, o admin
//   - user: cualquiera (sandbox personal)
func CanCreateCatalogScope(ac AuthContext, scope models.Scope, campaignID primitive.ObjectID) (bool, error) {
	switch scope {
	case models.ScopeGlobal:
		return ac.IsAdmin(), nil
	case models.ScopeCampaign:
		if ac.IsAdmin() {
			return true, nil
		}
		if campaignID.IsZero() {
			return false, nil
		}
		camp, err := campaignRepo.FindByID(campaignID)
		if err != nil {
			return false, err
		}
		return camp != nil && camp.DMID == ac.UserID, nil
	case models.ScopeUser:
		return true, nil
	}
	return false, nil
}
