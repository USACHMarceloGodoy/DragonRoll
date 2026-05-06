package repositories

import (
	"context"
	"dragonroll/backend/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// CatalogRepository es el CRUD genérico de las colecciones de catálogo
// (Class, Subclass, Spell, Feature, Race, Skill). Usa bson.M crudo para
// filtros y devuelve documentos sin tipo, porque cada servicio sabe a
// qué struct decodificar. Esto evita 6 archivos casi idénticos.
//
// La parte interesante es VisibilityFilter: combina los criterios de
// scope para devolver SOLO lo que el usuario puede ver:
//   - todo lo global
//   - + lo que él creó (scope=user)
//   - + lo de las campañas en las que es miembro (scope=campaign)
type CatalogRepository struct {
	Collection *mongo.Collection
}

func NewCatalogRepository(name string) *CatalogRepository {
	return &CatalogRepository{Collection: DB.GetCollection(name)}
}

// VisibilityFilter arma el filtro $or para listar documentos visibles.
func VisibilityFilter(userID primitive.ObjectID, memberCampaigns []primitive.ObjectID) bson.M {
	or := []bson.M{
		{"meta.scope": models.ScopeGlobal},
		{"meta.scope": models.ScopeUser, "meta.owner_id": userID},
	}
	if len(memberCampaigns) > 0 {
		or = append(or, bson.M{
			"meta.scope":       models.ScopeCampaign,
			"meta.campaign_id": bson.M{"$in": memberCampaigns},
		})
	}
	return bson.M{"$or": or}
}

func (r *CatalogRepository) InsertOne(doc interface{}) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := r.Collection.InsertOne(ctx, doc)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return res.InsertedID.(primitive.ObjectID), nil
}

// FindOneRaw decodifica al destino indicado.
func (r *CatalogRepository) FindOneRaw(filter bson.M, dst interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := r.Collection.FindOne(ctx, filter).Decode(dst)
	if err == mongo.ErrNoDocuments {
		return nil // no encontrado: dst sin tocar, sin error
	}
	return err
}

func (r *CatalogRepository) FindByID(id primitive.ObjectID, dst interface{}) error {
	return r.FindOneRaw(bson.M{"_id": id}, dst)
}

// FindManyRaw devuelve docs en bson.Raw para que el servicio decodifique
// al tipo correcto en su lado.
func (r *CatalogRepository) FindManyRaw(filter bson.M) ([]bson.Raw, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cur, err := r.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var out []bson.Raw
	for cur.Next(ctx) {
		out = append(out, append(bson.Raw{}, cur.Current...))
	}
	return out, cur.Err()
}

// FindManyByIDs es un atajo común al hacer populate.
func (r *CatalogRepository) FindManyByIDs(ids []primitive.ObjectID) ([]bson.Raw, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	return r.FindManyRaw(bson.M{"_id": bson.M{"$in": ids}})
}

func (r *CatalogRepository) ReplaceOne(id primitive.ObjectID, doc interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := r.Collection.ReplaceOne(ctx, bson.M{"_id": id}, doc)
	return err
}

func (r *CatalogRepository) DeleteOne(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := r.Collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// CountByFilter útil para validar unicidad de seeds.
func (r *CatalogRepository) CountByFilter(filter bson.M) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return r.Collection.CountDocuments(ctx, filter)
}

// CreateNameIndex es opcional — útil para Skills donde queremos unique
// (code, scope) — el servicio decide si llamarlo.
func (r *CatalogRepository) CreateNameIndex(field string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := r.Collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: field, Value: 1}, {Key: "meta.scope", Value: 1}},
	})
	return err
}
