package repository

import (
	"context"
	"database/sql"
	"errors"
)

var (
	ErrPostNotFound       = errors.New("the Post was not found in the repository")
	ErrFailedToCreatePost = errors.New("failed to create the Post to the repository")
	ErrFailedUpdatePost   = errors.New("failed to update the Post in the repository")
	ErrFailedToDeletePost = errors.New("failed to delete the Post in the repository")
	ErrFailedToScanPost   = errors.New("failed to scan the Post")
)

type PostDatabase interface {
	Create(ctx context.Context, tx *sql.Tx, ps Post) (Post, error)
	Update(ctx context.Context, tx *sql.Tx, ps Post) (Post, error)
	Delete(ctx context.Context, tx *sql.Tx, ps Post) error
	FindByID(ctx context.Context, tx *sql.Tx, ID uint64) (Post, error)
	FindByTitleContent(ctx context.Context, tx *sql.Tx, query string, from int, size int) ([]Post, error)
	FindAll(ctx context.Context, tx *sql.Tx) ([]Post, error)
}
