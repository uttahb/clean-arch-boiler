package http

import (
	"cleanarch/boiler/internal/user/domain"
	"cleanarch/boiler/internal/utils/logger"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type Handler struct {
	l             logger.Interface
	authUseCase   AuthUseCases
	tenantUsecase TenantUsecases
	userUseCase   UserUseCases
}

type AuthUseCases interface {
	SignUp(ctx context.Context, user *domain.AddUserRequest, tenantId string) error
	Login(ctx context.Context, user *domain.AddUserRequest) (*domain.UserTokens, error)
	ValidateAccessToken(ctx context.Context, token string) (string, error)
	RefreshTokenAccess(ctx context.Context, refreshToken string) (*domain.UserTokens, error)
}

type TenantUsecases interface {
	Create(ctx context.Context, tenantId string) error
}

type UserUseCases interface {
	GetUserByID(ctx context.Context, id string) (*domain.UserResponse, error)
}

func NewHandler(l logger.Interface, authUseCase AuthUseCases, userUseCase UserUseCases, tenant TenantUsecases) *Handler {
	return &Handler{
		l:             l,
		authUseCase:   authUseCase,
		userUseCase:   userUseCase,
		tenantUsecase: tenant,
	}
}

/**
 * Handlers start here
 */
func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", http.DetectContentType([]byte("pong")))
	w.Write([]byte("pong"))
}
func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(domain.UserKey{})
	SuccessResponse(user, "success").Send(w, r, http.StatusOK)
}

// swagger:route POST /signup signup signupRequest
// User registration can be done with this endpoint.
// responses:
//
//	200: signupSuccessResponse
//	409: signupErrorResponse
//	401: signupErrorResponse
//	500: signupErrorResponse
func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	tenantId := uuid.NewString()

	if err := h.tenantUsecase.Create(ctx, tenantId); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	addUserRequest := new(domain.AddUserRequest)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&addUserRequest)
	if err != nil {
		fmt.Println(err)
		ErrorResponse(err.Error()).Send(w, r, http.StatusBadRequest)
		return
	}

	err = h.authUseCase.SignUp(ctx, addUserRequest, tenantId)
	h.l.Debug("signup", "error", err)
	if err != nil {
		ErrorResponse(err.Error()).Send(w, r, http.StatusBadRequest)
		return
	}
	SuccessResponse("success", "signed up successfully").Send(w, r, http.StatusOK)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	addUserRequest := new(domain.AddUserRequest)

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&addUserRequest)
	if err != nil {
		ErrorResponse(err.Error()).Send(w, r, http.StatusBadRequest)
	}
	ctx := r.Context()
	tokens, err := h.authUseCase.Login(ctx, addUserRequest)
	if err != nil {
		if err == domain.ErrUserNotFound {
			ErrorResponse(err.Error()).Send(w, r, http.StatusUnauthorized)
			return
		}

		ErrorResponse(err.Error()).Send(w, r, http.StatusInternalServerError)
		return
	}
	h.setCookieValues(w, tokens)
	SuccessResponse(tokens, "Login successful").Send(w, r, http.StatusOK)
}

func (h *Handler) setCookieValues(w http.ResponseWriter, tokens *domain.UserTokens) {
	cookie := http.Cookie{
		Name:     domain.AccessTokenKey,
		Value:    tokens.AccessToken,
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		// Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	// Use the http.SetCookie() function to send the cookie to the client.
	// Behind the scenes this adds a `Set-Cookie` header to the response
	// containing the necessary cookie data.
	http.SetCookie(w, &cookie)

	cookie = http.Cookie{
		Name:     domain.RefreshTokenKey,
		Value:    tokens.RefreshToken,
		Path:     "/auth/refresh-access",
		MaxAge:   36000,
		HttpOnly: true,
		// Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &cookie)
}

/**
 * Middleware
 * This is the middleware that validates the access token.
 */
func (h *Handler) MiddlewareValidateAccessToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token, err := h.extractToken(r)
		if err != nil {
			h.l.Error("Token not provided or malformed")
			ErrorResponse(err.Error()).Send(w, r, http.StatusBadRequest)
			return
		}

		userId, err := h.authUseCase.ValidateAccessToken(r.Context(), token)
		if err != nil {
			h.l.Error("token validation failed", "error", err)
			ErrorResponse(err.Error()).Send(w, r, http.StatusBadRequest)
			return
		}
		ctx := context.WithValue(r.Context(), domain.UserIDKey{}, userId)
		r = r.WithContext(ctx)
		user, error := h.userUseCase.GetUserByID(ctx, userId)
		if error != nil {
			h.l.Error("unable to get user", "error", error)
			ErrorResponse(error.Error()).Send(w, r, http.StatusInternalServerError)
			return
		}
		ctx = context.WithValue(r.Context(), domain.UserKey{}, user)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
func (h *Handler) RefreshAccess(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := getTokenFromCookie(r, domain.RefreshTokenKey)
	if err != nil {
		ErrorResponse(err.Error()).Send(w, r, http.StatusBadRequest)
		return
	}
	tokens, err := h.authUseCase.RefreshTokenAccess(r.Context(), refreshToken)
	if err != nil {
		ErrorResponse(err.Error()).Send(w, r, http.StatusBadRequest)
		return
	}
	h.setCookieValues(w, tokens)
	SuccessResponse(tokens, "Refresh access token successful").Send(w, r, http.StatusOK)
}
func (h *Handler) extractToken(r *http.Request) (string, error) {
	token, err := getTokenFromCookie(r, domain.AccessTokenKey)
	if err != nil {
		h.l.Error("unable to get token from cookie", "error", err)
	} else {
		return token, nil
	}

	authHeader := r.Header.Get("Authorization")
	authHeaderContent := strings.Split(authHeader, " ")
	if len(authHeaderContent) != 2 {
		return "", errors.New("TOKEN NOT PROVIDED OR MALFORMED")
	}
	return authHeaderContent[1], nil
}
func getTokenFromCookie(r *http.Request, tokenType string) (string, error) {
	// Retrieve the cookie from the request using its name (which in our case is
	// "exampleCookie"). If no matching cookie is found, this will return a
	// http.ErrNoCookie error. We check for this, and return a 400 Bad Request
	// response to the client.
	cookie, err := r.Cookie(tokenType)
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			return "", http.ErrNoCookie
		default:
			log.Println(err)
			return "", err
		}

	}
	return cookie.Value, nil
	// Echo out the cookie value in the response body.
}
