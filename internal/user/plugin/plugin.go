package plugin

import (
	"cleanarch/boiler/internal/user/adapters/handlers/http"
	repositories "cleanarch/boiler/internal/user/adapters/repositories/mongo"
	"cleanarch/boiler/internal/user/services"
	"cleanarch/boiler/internal/user/usecases"
	"cleanarch/boiler/internal/utils/logger"
	"context"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserPlugin struct {
	r  chi.Router
	db *mongo.Database
	l  logger.Interface
}

func NewUserPlugin(r chi.Router, db *mongo.Database, logger logger.Interface) *UserPlugin {
	return &UserPlugin{
		db: db,
		r:  r,
		l:  logger,
	}
}
func (p *UserPlugin) Register() {
	userRepository := repositories.NewUserRepository(p.l, p.db)
	tenantRepository := repositories.NewTeanantRepository(p.db) //tenantrepository
	jwtRepository := repositories.NewJwtRepository(p.l, p.db)
	userService := services.NewUserService(p.l, userRepository)
	authService := services.NewAuthService(userRepository)
	jwtService := services.NewJwtService(p.l, jwtRepository)
	tenantService := services.NewTeanantService(tenantRepository) //tenantservice create
	authUsecase := usecases.NewAuthUseCases(p.l, authService, jwtService, userService)
	tenantUsecase := usecases.NewTenantUseCases(tenantService)
	userUsecase := usecases.NewUserUsecases(p.l, userService)
	accountHandler := http.NewHandler(p.l, authUsecase, userUsecase, tenantUsecase)
	http.RegisterAuthHTTPEndpoints(p.r, accountHandler)
	createDbIndices(p.db)
}
func createDbIndices(db *mongo.Database) error {
	// Create index on email
	_, err := db.Collection("users").Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys:    bson.M{"email": 1},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return err
	}
	return nil
}
