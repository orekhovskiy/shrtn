package middleware

import (
	"context"
	"github.com/orekhovskiy/shrtn/internal/logger"
	"go.uber.org/zap"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"

	"github.com/orekhovskiy/shrtn/config"
	"github.com/orekhovskiy/shrtn/internal/service/authservice"
)

const CookieName = "auth_token"

func AuthMiddleware(options config.Config, logger logger.Logger, toRequireAuth bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(CookieName)
			// Handle missing or empty cookie
			if cookie == nil || cookie.Value == "" || err != nil {
				if toRequireAuth {
					logger.Info("no auth cookie provided on protected route", zap.String("uri", r.RequestURI))
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}
				// Cookie is empty or not existing, generate a new one
				userID := issueNewToken(w, options.JWTSecretKey)
				ctx := context.WithValue(r.Context(), authservice.UserIDKey, userID)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			// Parse JWT from cookie
			claims := jwt.MapClaims{}
			token, err := jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(options.JWTSecretKey), nil
			})

			if err != nil || !token.Valid {
				if toRequireAuth {
					logger.Info("invalid auth cookie provided on protected route", zap.String("uri", r.RequestURI))
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}
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
