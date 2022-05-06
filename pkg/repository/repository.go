package repository

import (
	"context"
	"database/sql"
	"errors"

	user "github.com/izzanzahrial/blog-api-echo/pkg/user"
)

var (
	ErrPostNotFound       = errors.New("the Post was not found in the repository")
	ErrFailedToCreatePost = errors.New("failed to create the Post to the repository")
	ErrFailedUpdatePost   = errors.New("failed to update the Post in the repository")
	ErrFailedToDeletePost = errors.New("failed to delete the Post in the repository")
	ErrFailedToScanPost   = errors.New("failed to scan the Post")
	ErrUserNotFound       = errors.New("the user was not found in the repository")
	ErrFailedToCreateUser = errors.New("failed to create the user to the repository")
	ErrFailedUpdateUser   = errors.New("failed to update the user in the repository")
	ErrFailedToDeleteUser = errors.New("failed to delete the user in the repository")
	ErrFailedToAssertUser = errors.New("failed to assert the user")
)

type Post interface {
	Create(ctx context.Context, tx *sql.Tx, pd PostData) (PostData, error)
	Update(ctx context.Context, tx *sql.Tx, pd PostData) error
	Delete(ctx context.Context, tx *sql.Tx, pd PostData) error
	FindByID(ctx context.Context, tx *sql.Tx, id int64) (PostData, error)
	FindByTitleContent(ctx context.Context, tx *sql.Tx, query string, from int, size int) ([]PostData, error)
	FindRecent(ctx context.Context, tx *sql.Tx, from int, size int) ([]PostData, error)
}

type UserRepository interface {
	Create(ctx context.Context, tx *sql.Tx, user user.User) (user.User, error)
	UpdateUser(ctx context.Context, tx *sql.Tx, user user.User) (user.User, error)
	UpdatePassword(ctx context.Context, tx *sql.Tx, user user.User) (user.User, error)
	Delete(ctx context.Context, tx *sql.Tx, user user.User) error
	Login(ctx context.Context, tx *sql.Tx, ID uint64, pass string) (user.User, error)
	FindByID(ctx context.Context, tx *sql.Tx, ID uint64) (user.User, error)
}
