package post

import (
	"context"
	"database/sql"

	"github.com/izzanzahrial/blog-api-echo/entity"
)

type postgreRepository struct {
}

func NewPostgreRepository() PostRepository {
	return &postgreRepository{}
}

func (p *postgreRepository) Create(ctx context.Context, tx *sql.Tx, post entity.Post) (entity.Post, error) {
	SQL := "INSERT INTO post(title, content) VALUES (?, ?)"
	result, err := tx.ExecContext(ctx, SQL, post.Title, post.Content)
	if err != nil {
		return post, ErrFailedToAddPost
	}

	id, err := result.LastInsertId()
	if err != nil {
		return post, ErrPostNotFound
	}

	post.ID = uint64(id)

	return post, nil
}

func (p *postgreRepository) Update(ctx context.Context, tx *sql.Tx, post entity.Post) (entity.Post, error) {
	SQL := "UPDATE post SET title = ?, content = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, SQL, post.Title, post.Content, post.ID)
	if err != nil {
		return post, ErrFailedUpdatePost
	}

	return post, nil
}

func (p *postgreRepository) Delete(ctx context.Context, tx *sql.Tx, post entity.Post) error {
	SQL := "DELETE FROM post WHERE id = ?"
	_, err := tx.ExecContext(ctx, SQL, post.ID)
	if err != nil {
		return ErrFailedToDeletePost
	}

	return nil
}

func (p *postgreRepository) FindByID(ctx context.Context, tx *sql.Tx, ID uint64) (entity.Post, error) {
	SQL := "SELECT title, content FROM post WHERE id = ?"
	rows, err := tx.QueryContext(ctx, SQL, ID)
	if err != nil {
		return entity.Post{}, ErrPostNotFound
	}
	defer rows.Close()

	post := entity.Post{}
	if rows.Next() {
		if err := rows.Scan(&post.ID, &post.Title, &post.Content); err != nil {
			return post, ErrFailedToAssertPost
		}
		return post, nil
	} else {
		return post, ErrPostNotFound
	}
}

func (p *postgreRepository) FindAll(ctx context.Context, tx *sql.Tx) ([]entity.Post, error) {
	SQL := "SELECT id, title, content FROM post"
	rows, err := tx.QueryContext(ctx, SQL)
	if err != nil {
		return nil, ErrPostNotFound
	}
	defer rows.Close()

	var posts []entity.Post
	for rows.Next() {
		post := entity.Post{}
		if err := rows.Scan(&post.ID, &post.Title, &post.Content); err != nil {
			return nil, ErrFailedToAssertPost
		}
		posts = append(posts, post)
	}

	return posts, nil
}
