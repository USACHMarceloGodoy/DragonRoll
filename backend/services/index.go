package services

import (
	"dragonroll/backend/repositories"
	"errors"
)

var (
	userRepo      *repositories.UserRepository
	characterRepo *repositories.CharacterRepository
	campaignRepo  *repositories.CampaignRepository

	// Catálogos: todos usan el repo genérico, parametrizado por colección.
	classRepo    *repositories.CatalogRepository
	subclassRepo *repositories.CatalogRepository
	spellRepo    *repositories.CatalogRepository
	featureRepo  *repositories.CatalogRepository
	raceRepo     *repositories.CatalogRepository
	skillRepo    *repositories.CatalogRepository
)

func InitRepositories() {
	userRepo = repositories.NewUserRepository()
	characterRepo = repositories.NewCharacterRepository()
	campaignRepo = repositories.NewCampaignRepository()

	classRepo = repositories.NewCatalogRepository(repositories.CollectionClasses)
	subclassRepo = repositories.NewCatalogRepository(repositories.CollectionSubclasses)
	spellRepo = repositories.NewCatalogRepository(repositories.CollectionSpells)
	featureRepo = repositories.NewCatalogRepository(repositories.CollectionFeatures)
	raceRepo = repositories.NewCatalogRepository(repositories.CollectionRaces)
	skillRepo = repositories.NewCatalogRepository(repositories.CollectionSkills)
}

func InitIndexes() error {
	if err := userRepo.CreateIndexes(); err != nil {
		return errors.New("usuarios: " + err.Error())
	}
	if err := characterRepo.CreateIndexes(); err != nil {
		return errors.New("personajes: " + err.Error())
	}
	if err := campaignRepo.CreateIndexes(); err != nil {
		return errors.New("campañas: " + err.Error())
	}
	if err := skillRepo.CreateNameIndex("code"); err != nil {
		return errors.New("skills idx: " + err.Error())
	}
	if err := classRepo.CreateNameIndex("name"); err != nil {
		return errors.New("classes idx: " + err.Error())
	}
	return nil
}
