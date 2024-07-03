package mongo

import (
	"cleanarch/boiler/internal/user/domain"
	"cleanarch/boiler/internal/utils/logger"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type JwtRepository struct {
	l  logger.Interface
	db *mongo.Database
}

func NewJwtRepository(l logger.Interface, db *mongo.Database) *JwtRepository {
	return &JwtRepository{
		l:  l,
		db: db,
	}
}

func (r JwtRepository) CreateRefreshJwt(ctx context.Context, jwt *domain.Jwt, currentRefreshId string) error {
	//TODO : put this in a transaction
	// TODO : Handle this error well with condition check for ErrNoDocuments
	r.l.Info("creating refresh token", currentRefreshId)
	if currentRefreshId != "" {

		result := r.db.Collection("jwt").FindOneAndDelete(ctx, bson.M{
			"id": currentRefreshId,
		})
		if result.Err() != nil {
			r.l.Error("unable to delete refresh token", "error", result.Err())
		}
	}
	_, error := r.db.Collection("jwt").InsertOne(ctx, jwt)
	return error
}
