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

type Post interface {
	Create(ctx context.Context, tx *sql.Tx, pd PostData) (PostData, error)
	Update(ctx context.Context, tx *sql.Tx, pd PostData) error
	Delete(ctx context.Context, tx *sql.Tx, pd PostData) error
	FindByID(ctx context.Context, tx *sql.Tx, id int64) (PostData, error)
	FindByTitleContent(ctx context.Context, tx *sql.Tx, query string, from int, size int) ([]PostData, error)
	FindRecent(ctx context.Context, tx *sql.Tx, from int, size int) ([]PostData, error)
}
