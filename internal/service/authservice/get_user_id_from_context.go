package authservice

import "context"

type contextKey string

const UserIDKey contextKey = "userID"

func (s *AuthService) GetUserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDKey).(string)
	return userID, ok
}
