package core

import (
	"context"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user User, password string) (*User, error)
	UpdateUser(ctx context.Context, userUUID string, ur UserUpdateRequest) (*User, error)
}
