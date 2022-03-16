package post

import (
	"context"
	"database/sql"

	"github.com/izzanzahrial/blog-api-echo/entity"
)

var ()

type PostRepository interface {
	Create(ctx context.Context, tx *sql.Tx, post entity.Post) entity.Post
	Update(ctx context.Context, tx *sql.Tx, post entity.Post) entity.Post
	Delete(ctx context.Context, tx *sql.Tx, post entity.Post)
	FindByID(ctx context.Context, tx *sql.Tx, ID uint64) (*entity.Post, error)
	FindAll(ctx context.Context, tx *sql.Tx) []entity.Post
}
