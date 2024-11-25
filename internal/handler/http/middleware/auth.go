package middleware

import (
	"context"
	"github.com/orekhovskiy/shrtn/internal/logger"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"

	"github.com/orekhovskiy/shrtn/config"
	"github.com/orekhovskiy/shrtn/internal/service/authservice"
)

const CookieName = "auth_token"

func AuthMiddleware(options config.Config, logger logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(CookieName)
			if err != nil || cookie.Value == "" {
				// Cookie is empty or not existing, generate a new one
				userID := issueNewToken(w, options.JWTSecretKey)
				ctx := context.WithValue(r.Context(), authservice.UserIDKey, userID)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			// Check for JWT from Cookie
			claims := jwt.MapClaims{}
			token, err := jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(options.JWTSecretKey), nil
			})

			if err != nil || !token.Valid {
				// Cookie is not valid, generate a new one
				userID := issueNewToken(w, options.JWTSecretKey)
				ctx := context.WithValue(r.Context(), authservice.UserIDKey, userID)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			userID, ok := claims["user_id"].(string)
			if !ok {
				logger.Info("no user ID provided, rejecting")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), authservice.UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func issueNewToken(w http.ResponseWriter, secretKey string) string {
	// Generate new userID
	userID := uuid.New().String()

	// generate JWT
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return ""
	}

	// Set JWT inside a Cookie
	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    tokenString,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Path:     "/",
	})

	return userID
}
