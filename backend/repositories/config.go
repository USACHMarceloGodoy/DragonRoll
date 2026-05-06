package repositories

import "dragonroll/backend/config/database"

const (
	CollectionUsers      = "users"
	CollectionCharacters = "characters"
	CollectionCampaigns  = "campaigns"

	// Catálogo data-driven
	CollectionClasses    = "classes"
	CollectionSubclasses = "subclasses"
	CollectionSpells     = "spells"
	CollectionFeatures   = "features"
	CollectionRaces      = "races"
	CollectionSkills     = "skills"
)

var (
	DB *database.MongoConnection
)

func InitConnections() {
	DB = database.NewMongoConnection()
}
