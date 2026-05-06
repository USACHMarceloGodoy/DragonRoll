package services

import (
	"dragonroll/backend/models"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateCampaignService(ac AuthContext, c *models.Campaign) (*models.Campaign, error) {
	if c.Name == "" {
		return nil, errors.New("nombre obligatorio")
	}
	c.DMID = ac.UserID // el creador queda como DM por definición
	c.Active = true
	now := time.Now()
	c.CreatedAt = now
	c.UpdatedAt = now
	if c.System == "" {
		c.System = "D&D 5e"
	}
	id, err := campaignRepo.InsertOne(c)
	if err != nil {
		return nil, err
	}
	c.ID = id
	return c, nil
}

func ListMyCampaigns(ac AuthContext) ([]*models.Campaign, error) {
	return campaignRepo.FindByMember(ac.UserID)
}

func GetCampaignAuthorized(ac AuthContext, id primitive.ObjectID) (*models.Campaign, error) {
	c, err := campaignRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, errors.New("no encontrado")
	}
	if ac.IsAdmin() || c.DMID == ac.UserID {
		return c, nil
	}
	for _, p := range c.PlayerIDs {
		if p == ac.UserID {
			return c, nil
		}
	}
	return nil, errors.New("no autorizado")
}

func UpdateCampaignService(ac AuthContext, id primitive.ObjectID, updated *models.Campaign) (*models.Campaign, error) {
	current, err := campaignRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if current == nil {
		return nil, errors.New("no encontrado")
	}
	if !ac.IsAdmin() && current.DMID != ac.UserID {
		return nil, errors.New("solo el DM o admin pueden editar la campaña")
	}
	updated.ID = current.ID
	updated.DMID = current.DMID
	updated.CreatedAt = current.CreatedAt
	updated.UpdatedAt = time.Now()
	if err := campaignRepo.UpdateOne(id, updated); err != nil {
		return nil, err
	}
	return updated, nil
}

// AddPlayerToCampaign: el DM mete a un jugador (por user ID).
func AddPlayerToCampaign(ac AuthContext, campID, playerID primitive.ObjectID) error {
	c, err := campaignRepo.FindByID(campID)
	if err != nil || c == nil {
		return errors.New("campaña no encontrada")
	}
	if !ac.IsAdmin() && c.DMID != ac.UserID {
		return errors.New("no autorizado")
	}
	for _, p := range c.PlayerIDs {
		if p == playerID {
			return nil // ya está
		}
	}
	c.PlayerIDs = append(c.PlayerIDs, playerID)
	c.UpdatedAt = time.Now()
	return campaignRepo.UpdateOne(campID, c)
}

func DeleteCampaignService(ac AuthContext, id primitive.ObjectID) error {
	c, err := campaignRepo.FindByID(id)
	if err != nil || c == nil {
		return errors.New("campaña no encontrada")
	}
	if !ac.IsAdmin() && c.DMID != ac.UserID {
		return errors.New("no autorizado")
	}
	return campaignRepo.DeleteOne(id)
}
