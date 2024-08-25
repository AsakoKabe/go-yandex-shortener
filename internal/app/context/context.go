package context

import "context"

type UserID struct {
	value string
}

func SetUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, UserID{}, userID)
}

func GetUserID(ctx context.Context) string {
	userID, ok := ctx.Value(UserID{}).(string)
	if !ok {
		return ""
	}

	return userID
}
