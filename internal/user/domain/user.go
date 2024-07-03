package domain

import "github.com/go-playground/validator/v10"

type UserIDKey struct{}
type UserKey struct{}

type User struct {
	ID           string
	FirstName    string
	MiddleName   string
	LastName     string
	Email        string
	Password     string
	PasswordSalt string
	CreatedAt    string
	UpdatedAt    string
	Status       UserStatus
}
type UserResponse struct {
	ID         string
	FirstName  string
	MiddleName string
	LastName   string
	Email      string
	Status     UserStatus
}
type UserTokens struct {
	AccessToken  string
	RefreshToken string
}

const AccessTokenKey = "access_token"
const RefreshTokenKey = "refresh_token"

type AddUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=20"`
}

func (u *AddUserRequest) ValidateAddUserRequest() error {
	return validator.New().Struct(u)
}
func (u *User) FullName() string {
	return u.FirstName + " " + u.MiddleName + " " + u.LastName
}

type UserStatus int

const (
	Invited        UserStatus = iota // EnumIndex = 0
	InviteAccepted                   // EnumIndex = 1
	Active                           // EnumIndex = 2
	Inactive                         // EnumIndex = 3
	Suspended                        // EnumIndex = 4
)

// String - Creating common behavior - give the type a String function
func (w UserStatus) String() string {
	return [...]string{"Invited", "InviteAccepted", "Active", "Inactive", "Suspended"}[w]
}

// EnumIndex - Creating common behavior - give the type a EnumIndex function
func (w UserStatus) EnumIndex() int {
	return int(w)
}
