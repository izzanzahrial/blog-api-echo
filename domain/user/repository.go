package user

import (
	"context"
	"database/sql"
	"errors"

	"github.com/izzanzahrial/blog-api-echo/entity"
)

var (
	ErrUserNotFound       = errors.New("the user was not found in the repository")
	ErrFailedToCreateUser = errors.New("failed to create the user to the repository")
	ErrFailedUpdateUser   = errors.New("failed to update the user in the repository")
	ErrFailedToDeleteUser = errors.New("failed to delete the user in the repository")
	ErrFailedToAssertUser = errors.New("failed to assert the user")
)

type UserRepository interface {
	Create(ctx context.Context, tx *sql.Tx, user entity.User) (entity.User, error)
	Update(ctx context.Context, tx *sql.Tx, user entity.User) (entity.User, error)
	Delete(ctx context.Context, tx *sql.Tx, user entity.User) error
	Find(ctx context.Context, tx *sql.Tx, ID uint64, password string) (entity.User, error)
}
