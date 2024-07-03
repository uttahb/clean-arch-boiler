package services

import (
	"cleanarch/boiler/internal/user/domain"
	"cleanarch/boiler/internal/utils/logger"
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AccessTokenCustomClaims struct {
	UserID string `json:"user_id"`
	Type   string `json:"type"`
	jwt.RegisteredClaims
}
type RefreshTokenCustomClaims struct {
	UserID string `json:"user_id"`
	Type   string `json:"type"`
	jwt.RegisteredClaims
}
type JwtService struct {
	l             logger.Interface
	jwtRepository JwtRepository
}
type JwtRepository interface {
	CreateRefreshJwt(ctx context.Context, jwt *domain.Jwt, currentRefreshId string) error
}

func NewJwtService(l logger.Interface, jwtRepository JwtRepository) *JwtService {
	return &JwtService{
		l:             l,
		jwtRepository: jwtRepository,
	}
}

func (s *JwtService) ValidateAccessToken(ctx context.Context, tokenString string) (string, error) {

	token, err := jwt.ParseWithClaims(tokenString, &AccessTokenCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			s.l.Error("Unexpected signing method in auth token")
			return nil, errors.New("UNEXPECTED SIGNING METHOD IN AUTH TOKEN")
		}
		verifyBytes, err := os.ReadFile("./keys/auth-public.pem")
		if err != nil {
			s.l.Error("unable to read public key", "error", err)
			return nil, err
		}

		verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
		if err != nil {
			s.l.Error("unable to parse public key", "error", err)
			return nil, err
		}

		return verifyKey, nil
	})

	if err != nil {
		s.l.Error("unable to parse claims", "error", err)
		return "", err
	}

	claims, ok := token.Claims.(*AccessTokenCustomClaims)
	if !ok || !token.Valid || claims.UserID == "" || claims.Type != "access" {
		return "", errors.New("invalid token: authentication failed")
	}
	return claims.UserID, nil
}
func (s *JwtService) GenerateAccessToken(ctx context.Context, user *domain.UserResponse) (string, error) {
	userID := user.ID
	tokenType := "access"

	claims := AccessTokenCustomClaims{
		userID,
		tokenType,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(10))),
			Issuer:    "cleanarch.service",
			ID:        uuid.NewString(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	signBytes, err := os.ReadFile("./keys/auth-private.pem")
	if err != nil {
		s.l.Error("unable to read private key", "error", err)
		return "", errors.New("could not generate access token. please try again later")
	}

	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		s.l.Error("unable to read private key", "error", err)
		return "", errors.New("could not generate access token. please try again later")
	}

	return token.SignedString(signKey)
}

func (s *JwtService) GenerateRefreshToken(ctx context.Context, user *domain.UserResponse, currentRefreshTokenId string) (string, error) {
	tokenType := "refresh"
	tokenId := uuid.NewString()
	issuedAt := time.Now()
	expiresAt := issuedAt.Add(time.Minute * time.Duration(100))
	claims := RefreshTokenCustomClaims{
		user.ID,
		tokenType,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			Issuer:    "cleanarch.service",
			ID:        tokenId,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	signBytes, err := os.ReadFile("./keys/auth-refresh-private.pem")
	if err != nil {
		s.l.Error("unable to read private key", "error", err)
		return "", errors.New("could not generate access token. please try again later")
	}

	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		s.l.Error("unable to read private key", "error", err)
		return "", errors.New("could not generate access token. please try again later")
	}
	signedToken, err := token.SignedString(signKey)
	if err != nil {
		s.l.Error("unable to read private key", "error", err)
		return "", errors.New("could not generate access token. please try again later")
	}
	jwtToken := domain.Jwt{
		Type:         tokenType,
		ID:           tokenId,
		UserID:       user.ID,
		RefreshToken: signedToken,
		CreatedAt:    fmt.Sprint(issuedAt.Unix()),
	}
	err = s.jwtRepository.CreateRefreshJwt(ctx, &jwtToken, currentRefreshTokenId)
	if err != nil {
		s.l.Error("unable to read private key", "error", err)
		return "", errors.New("could not generate access token. please try again later")
	}
	return signedToken, nil
}

func (s *JwtService) RefreshTokenAccess(ctx context.Context, refreshToken string) (string, string, error) {

	oldClaims, err := s.parseRefreshTokenWithClaims(refreshToken)
	if err != nil {
		return "", "", err
	}
	if oldClaims.Type != "refresh" {
		return "", "", errors.New("INVALID TOKEN TYPE")
	}
	return oldClaims.UserID, oldClaims.ID, nil
	// tokenType := "refresh"
	// tokenId := uuid.NewString()
	// issuedAt := time.Now()
	// expiresAt := issuedAt.Add(time.Minute * time.Duration(100))
	// claims := RefreshTokenCustomClaims{
	// 	oldClaims.UserID,
	// 	tokenType,
	// 	jwt.RegisteredClaims{
	// 		ExpiresAt: jwt.NewNumericDate(expiresAt),
	// 		Issuer:    "cleanarch.service",
	// 		ID:        tokenId,
	// 	},
	// }

	// token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	// signBytes, err := os.ReadFile("./keys/auth-refresh-private.pem")
	// if err != nil {
	// 	s.l.Error("unable to read private key", "error", err)
	// 	return "", errors.New("could not generate access token. please try again later")
	// }

	// signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	// if err != nil {
	// 	s.l.Error("unable to read private key", "error", err)
	// 	return "", errors.New("could not generate access token. please try again later")
	// }
	// signedToken, err := token.SignedString(signKey)
	// if err != nil {
	// 	s.l.Error("unable to read private key", "error", err)
	// 	return "", errors.New("could not generate access token. please try again later")
	// }
	// jwtToken := domain.Jwt{
	// 	Type:         tokenType,
	// 	ID:           tokenId,
	// 	UserID:       oldClaims.UserID,
	// 	RefreshToken: signedToken,
	// 	CreatedAt:    fmt.Sprint(issuedAt.Unix()),
	// }
	// err = s.jwtRepository.CreateRefreshJwt(ctx, &jwtToken, oldClaims.ID)
	// if err != nil {
	// 	s.l.Error("unable to read private key", "error", err)
	// 	return "", errors.New("could not generate access token. please try again later")
	// }
	// return signedToken, nil
}
func (s *JwtService) parseRefreshTokenWithClaims(token string) (*RefreshTokenCustomClaims, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &RefreshTokenCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			s.l.Error("UNEXPECTED SIGNING METHOD IN AUTH TOKEN")
			return nil, errors.New("UNEXPECTED SIGNING METHOD IN AUTH TOKEN")
		}
		verifyBytes, err := os.ReadFile("./keys/auth-refresh-public.pem")
		if err != nil {
			s.l.Error("UNABLE TO READ PUBLIC KEY", "error", err)
			return nil, err
		}

		verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
		if err != nil {
			s.l.Error("UNABLE TO PARSE PUBLIC KEY", "error", err)
			return nil, err
		}

		return verifyKey, nil
	})
	if err != nil {
		s.l.Error("UNABLE TO PARSE TOKEN", "error", err)
		return nil, err
	}
	return parsedToken.Claims.(*RefreshTokenCustomClaims), nil

}
