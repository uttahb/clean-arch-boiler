package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type TenantRepository struct {
	db *mongo.Database
}

func NewTeanantRepository(db *mongo.Database) *TenantRepository {
	return &TenantRepository{
		db: db,
	}
}

func (r *TenantRepository) Create(ctx context.Context, tenantId string) error {
	_, error := r.db.Collection("tenants").InsertOne(ctx, bson.M{
		"_id": tenantId,
	})
	if error != nil {
		return error
	}
	return nil
}
