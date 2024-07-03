package mongo

import (
	"context"
	"errors"

	"cleanarch/boiler/internal/user/domain"
	"cleanarch/boiler/internal/utils/logger"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Email    string             `bson:"email"`
	Password string             `bson:"password"`
}

type UserRepository struct {
	l  logger.Interface
	db *mongo.Database
}

func NewUserRepository(l logger.Interface, db *mongo.Database) *UserRepository {
	return &UserRepository{
		l:  l,
		db: db,
	}
}

// CreateUser creates a new user in the database. It first hashes the provided password
// using the bcrypt algorithm, then inserts a new user document into the "users" collection
// with the hashed password. If the password hashing or database insert operation fails,
// an error is returned.
func (r UserRepository) CreateUser(ctx context.Context, user *domain.AddUserRequest) error {
	hash, err := hashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = string(hash)
	_, error := r.db.Collection("users").InsertOne(ctx, user)
	if error != nil {
		return err
	}
	return nil
}

// GetUser retrieves a user from the database by their username, and verifies that the provided
// password matches the stored hashed password for that user. If the user is found and the
// password is valid, a User model is returned. If the user is not found or the password is
// invalid, an error is returned.
func (r UserRepository) GetUser(ctx context.Context, username, password string) (*domain.User, error) {
	user := new(User)
	err := r.db.Collection("users").FindOne(ctx, bson.M{
		"username": username,
	}).Decode(user)

	if err != nil {
		return nil, err
	}
	isAuthorized := doPasswordsMatch(user.Password, password)
	if !isAuthorized {
		return nil, domain.ErrUserNotFound
	}
	return toModel(user), nil
}

// GetUserByID retrieves a user from the database by their unique identifier (userId).
// It first converts the userId string to a MongoDB ObjectID, then uses that to find the
// corresponding user document in the "users" collection. The password field is excluded
// from the returned user data.
// If the user is found, a UserResponse is returned containing the user's ID and email.
// If the user is not found, an error is returned.
func (r UserRepository) GetUserByID(ctx context.Context, userId string) (*domain.UserResponse, error) {
	user := new(User)
	objID, _ := primitive.ObjectIDFromHex(userId)

	err := r.db.Collection("users").FindOne(ctx, bson.M{
		"_id": objID,
	}, options.FindOne().SetProjection(bson.M{"password": 0})).Decode(user)

	if err != nil {
		return nil, err
	}
	return toResponse(user), nil
}

// Hash password
func hashPassword(password string) (string, error) {
	// Convert password string to byte slice
	var passwordBytes = []byte(password)

	// Hash password with Bcrypt's min cost
	hashedPasswordBytes, err := bcrypt.
		GenerateFromPassword(passwordBytes, 10)

	return string(hashedPasswordBytes), err
}

// Check if two passwords match using Bcrypt's CompareHashAndPassword
// which return nil on success and an error on failure.
func doPasswordsMatch(hashedPassword, currPassword string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword), []byte(currPassword))
	return err == nil
}

// SignUp creates a new user in the database with the provided user information.
// It first hashes the user's password using the hashPassword function, then
// inserts the new user into the "users" collection in the database.
// If any errors occur during the process, they are returned.
func (r UserRepository) SignUp(ctx context.Context, user *domain.AddUserRequest, tenantId string) error {
	hash, err := hashPassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = string(hash)
	_, error := r.db.Collection("users").InsertOne(ctx, bson.M{
		"email":    user.Email,
		"password": user.Password,
		"tenantId": tenantId,
	})

	if error != nil {
		if IsDup(error) {
			return domain.ErrUserAlreadyExists
		}
		return error
	}
	return nil
}

// GetAuthenticatedUser retrieves a user from the database based on the provided email and password.
// If the user is found and the password matches, it returns a UserResponse containing the user's ID and email.
// If the user is not found or the password does not match, it returns an error.
func (r UserRepository) GetAuthenticatedUser(ctx context.Context, user *domain.AddUserRequest) (*domain.UserResponse, error) {
	dbUser := new(User)

	err := r.db.Collection("users").FindOne(ctx, bson.M{
		"email": user.Email,
	}).Decode(dbUser)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	if !doPasswordsMatch(dbUser.Password, user.Password) {
		return nil, domain.ErrUserNotFound
	}
	return toResponse(dbUser), nil
}
func toModel(u *User) *domain.User {
	return &domain.User{
		ID:       u.ID.Hex(),
		Email:    u.Email,
		Password: u.Password,
	}
}

// / toResponse converts a User model to a UserResponse model.
// / It extracts the ID and Email fields from the User and returns a new UserResponse.
func toResponse(u *User) *domain.UserResponse {
	return &domain.UserResponse{
		ID:    u.ID.Hex(),
		Email: u.Email,
	}
}

// / IsDup checks if the provided error is a MongoDB duplicate key error.
// / It returns true if the error is a duplicate key error, false otherwise.
func IsDup(err error) bool {
	var e mongo.WriteException
	if errors.As(err, &e) {
		for _, we := range e.WriteErrors {
			if we.Code == 11000 {
				return true
			}
		}
	}
	return false
}
