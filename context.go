package main

import "context"

type userContextKey struct{}

func WithUser(ctx context.Context, user string) context.Context {
	return context.WithValue(ctx, userContextKey{}, user)
}

func GetUserFromContext(ctx context.Context) (string, bool) {
	user, found := ctx.Value(userContextKey{}).(string)
	return user, found
}