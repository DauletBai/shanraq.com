package session

import (
	"context"

	"shanraq.com/internal/auth"
)

// contextKey is an unexported type to prevent collisions in context.
type contextKey string

const identityContextKey contextKey = "shanraq.com/auth/identity"

// WithIdentity attaches the identity to the provided context.
func WithIdentity(ctx context.Context, identity auth.Identity) context.Context {
	return context.WithValue(ctx, identityContextKey, identity)
}

// IdentityFromContext retrieves an identity previously stored in the context, if present.
func IdentityFromContext(ctx context.Context) (auth.Identity, bool) {
	value := ctx.Value(identityContextKey)
	if value == nil {
		return auth.Identity{}, false
	}
	identity, ok := value.(auth.Identity)
	return identity, ok
}
