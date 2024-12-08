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
			if isInvalidToken(cookie, err) {
				handleInvalidToken(w, r, next, options, logger, toRequireAuth)
				return
			}

			claims := jwt.MapClaims{}
			token, err := jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(options.JWTSecretKey), nil
			})

			if err != nil || !token.Valid {
				handleInvalidToken(w, r, next, options, logger, toRequireAuth)
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

func isInvalidToken(cookie *http.Cookie, err error) bool {
	return cookie == nil || cookie.Value == "" || err != nil
}

func handleInvalidToken(w http.ResponseWriter, r *http.Request, next http.Handler, options config.Config, logger logger.Logger, toRequireAuth bool) {
	if toRequireAuth {
		logger.Info("invalid or missing auth cookie", zap.String("uri", r.RequestURI))
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, err := issueNewToken(w, options.JWTSecretKey)
	if err != nil {
		logger.Error("failed to issue new token", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	ctx := context.WithValue(r.Context(), authservice.UserIDKey, userID)
	next.ServeHTTP(w, r.WithContext(ctx))
}

func issueNewToken(w http.ResponseWriter, secretKey string) (string, error) {
	userID := uuid.New().String()

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    tokenString,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Path:     "/",
	})

	return userID, nil
}
