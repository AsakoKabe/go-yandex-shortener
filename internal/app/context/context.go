package context

import "context"

// UserID хранение id пользователя из jwt
type UserID struct {
	value string
}

// SetUserID установить userID в контекст
func SetUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, UserID{}, userID)
}

// GetUserID получить userID из контекста
func GetUserID(ctx context.Context) string {
	userID, ok := ctx.Value(UserID{}).(string)
	if !ok {
		return ""
	}

	return userID
}
