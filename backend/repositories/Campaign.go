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

type CampaignRepository struct {
	Collection *mongo.Collection
}

func NewCampaignRepository() *CampaignRepository {
	return &CampaignRepository{Collection: DB.GetCollection(CollectionCampaigns)}
}

func (r *CampaignRepository) InsertOne(c *models.Campaign) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := r.Collection.InsertOne(ctx, c)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return res.InsertedID.(primitive.ObjectID), nil
}

func (r *CampaignRepository) FindByID(id primitive.ObjectID) (*models.Campaign, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var c *models.Campaign
	err := r.Collection.FindOne(ctx, bson.M{"_id": id}).Decode(&c)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return c, nil
}

// FindByDM lista las campañas donde el usuario es DM.
func (r *CampaignRepository) FindByDM(dmID primitive.ObjectID) ([]*models.Campaign, error) {
	return r.findMany(bson.M{"dm_id": dmID})
}

// FindByPlayer lista las campañas donde el usuario es jugador.
func (r *CampaignRepository) FindByPlayer(playerID primitive.ObjectID) ([]*models.Campaign, error) {
	return r.findMany(bson.M{"player_ids": playerID})
}

// FindByMember devuelve campañas donde el usuario es DM o jugador.
func (r *CampaignRepository) FindByMember(userID primitive.ObjectID) ([]*models.Campaign, error) {
	return r.findMany(bson.M{
		"$or": []bson.M{
			{"dm_id": userID},
			{"player_ids": userID},
		},
	})
}

func (r *CampaignRepository) findMany(filter bson.M) ([]*models.Campaign, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cur, err := r.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var out []*models.Campaign
	if err := cur.All(ctx, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *CampaignRepository) UpdateOne(id primitive.ObjectID, c *models.Campaign) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := r.Collection.ReplaceOne(ctx, bson.M{"_id": id}, c)
	return err
}

func (r *CampaignRepository) PatchFields(id primitive.ObjectID, fields bson.M) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := r.Collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": fields})
	return err
}

func (r *CampaignRepository) DeleteOne(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := r.Collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *CampaignRepository) CreateIndexes() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := r.Collection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "dm_id", Value: 1}}, Options: options.Index()},
		{Keys: bson.D{{Key: "player_ids", Value: 1}}, Options: options.Index()},
	})
	return err
}
