package repositories

import (
	"context"
	"dragonroll/backend/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CharacterRepository struct {
	Collection *mongo.Collection
}

func NewCharacterRepository() *CharacterRepository {
	return &CharacterRepository{Collection: DB.GetCollection(CollectionCharacters)}
}

func (r *CharacterRepository) InsertOne(c *models.Character) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := r.Collection.InsertOne(ctx, c)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return res.InsertedID.(primitive.ObjectID), nil
}

func (r *CharacterRepository) FindByID(id primitive.ObjectID) (*models.Character, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var c *models.Character
	err := r.Collection.FindOne(ctx, bson.M{"_id": id}).Decode(&c)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return c, nil
}

func (r *CharacterRepository) FindByOwner(ownerID primitive.ObjectID) ([]*models.Character, error) {
	return r.findMany(bson.M{"owner_id": ownerID})
}

func (r *CharacterRepository) FindByCampaign(campaignID primitive.ObjectID) ([]*models.Character, error) {
	return r.findMany(bson.M{"campaign_id": campaignID})
}

// FindAccessible devuelve los PJ que el usuario puede ver:
//   - sus propios personajes
//   - + los personajes de cualquier campaña donde es DM (campaignDMIDs)
//   - + los de campañas donde es jugador (campaignPlayerIDs) — excepto privados
func (r *CharacterRepository) FindAccessible(userID primitive.ObjectID, dmCampaigns, playerCampaigns []primitive.ObjectID) ([]*models.Character, error) {
	or := []bson.M{{"owner_id": userID}}
	if len(dmCampaigns) > 0 {
		or = append(or, bson.M{"campaign_id": bson.M{"$in": dmCampaigns}})
	}
	if len(playerCampaigns) > 0 {
		or = append(or, bson.M{"campaign_id": bson.M{"$in": playerCampaigns}})
	}
	return r.findMany(bson.M{"$or": or})
}

func (r *CharacterRepository) findMany(filter bson.M) ([]*models.Character, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cur, err := r.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var out []*models.Character
	if err := cur.All(ctx, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *CharacterRepository) UpdateOne(id primitive.ObjectID, c *models.Character) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := r.Collection.ReplaceOne(ctx, bson.M{"_id": id}, c)
	return err
}

func (r *CharacterRepository) PatchFields(id primitive.ObjectID, fields bson.M) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := r.Collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": fields})
	return err
}

func (r *CharacterRepository) DeleteOne(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := r.Collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *CharacterRepository) CreateIndexes() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := r.Collection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "owner_id", Value: 1}}, Options: options.Index()},
		{Keys: bson.D{{Key: "campaign_id", Value: 1}}, Options: options.Index()},
	})
	return err
}
