package post

import (
	"context"
	"database/sql"
	"errors"

	"github.com/izzanzahrial/blog-api-echo/entity"
)

var (
	ErrPostNotFound       = errors.New("the post was not found in the repository")
	ErrFailedToAddPost    = errors.New("failed to add the post to the repository")
	ErrFailedUpdatePost   = errors.New("failed to update the post in the repository")
	ErrFailedToDeletePost = errors.New("failed to delete the post in the repository")
	ErrFailedToAssertPost = errors.New("failed to assert the post")
)

type PostRepository interface {
	Create(ctx context.Context, tx *sql.Tx, post entity.Post) (entity.Post, error)
	Update(ctx context.Context, tx *sql.Tx, post entity.Post) (entity.Post, error)
	Delete(ctx context.Context, tx *sql.Tx, post entity.Post) error
	FindByID(ctx context.Context, tx *sql.Tx, ID uint64) (entity.Post, error)
	FindAll(ctx context.Context, tx *sql.Tx) ([]entity.Post, error)
}
