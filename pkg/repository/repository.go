package repository

import (
	"context"
	"database/sql"
	"errors"
)

var (
	ErrPostNotFound       = errors.New("the Post was not found in the repository")
	ErrFailedToAddPost    = errors.New("failed to add the Post to the repository")
	ErrFailedUpdatePost   = errors.New("failed to update the Post in the repository")
	ErrFailedToDeletePost = errors.New("failed to delete the Post in the repository")
	ErrFailedToAssertPost = errors.New("failed to assert the Post")
)

type PostDatabase interface {
	Create(ctx context.Context, tx *sql.Tx, ps Post) (Post, error)
	Update(ctx context.Context, tx *sql.Tx, ps Post) (Post, error)
	Delete(ctx context.Context, tx *sql.Tx, ps Post) error
	FindByID(ctx context.Context, tx *sql.Tx, ID uint64) (Post, error)
	FindAll(ctx context.Context, tx *sql.Tx) ([]Post, error)
}
