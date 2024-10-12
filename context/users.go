package context

import (
	"context"

	"github.com/SuWh1/WebDevGo/models"
)

type key string

const (
	userKey key = "user" // unique to avoid key collisions with other values stored in the context by other packages or parts of the code
)

func WithUser(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, userKey, user)
} // stores the User object in the context.
// returns a new context that carries the User object.

func GetUser(ctx context.Context) *models.User {
	val := ctx.Value(userKey)
	user, ok := val.(*models.User)
	if !ok {
		// The most likely case is that nothing was ever stored in the context,
		// so it does not have a type of *models.User. It is also possible that
		// other code in this package wrote an invalid value using the user key.
		return nil
	}
	return user
} // retrieving user from the context
