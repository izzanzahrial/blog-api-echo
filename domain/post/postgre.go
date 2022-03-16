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

func (p *postgreRepository) Create(ctx context.Context, tx *sql.Tx, post entity.Post) entity.Post {
	SQL := "INSERT INTO post(title, content) VALUES (?, ?)"
	result, err := tx.ExecContext(ctx, SQL, post.Title, post.Content)
	if err != nil {

	}

	id, err := result.LastInsertId()
	if err != nil {

	}

	post.ID = uint64(id)
	return post
}

func (p *postgreRepository) Update(ctx context.Context, tx *sql.Tx, post entity.Post) entity.Post {
	SQL := "UPDATE post SET title = ?, content = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, SQL, post.Title, post.Content, post.ID)
	if err != nil {

	}

	return post
}

func (p *postgreRepository) Delete(ctx context.Context, tx *sql.Tx, post entity.Post) {
	SQL := "DELETE FROM post WHERE id = ?"
	_, err := tx.ExecContext(ctx, SQL, post.ID)
	if err != nil {

	}
}

func (p *postgreRepository) FindByID(ctx context.Context, tx *sql.Tx, ID uint64) (*entity.Post, error) {
	SQL := "SELECT title, content FROM post WHERE id = ?"
	rows, err := tx.QueryContext(ctx, SQL, ID)
	if err != nil {

	}
	defer rows.Close()

	post := entity.Post{}
	if rows.Next() {
		err := rows.Scan(&post.ID, &post.Title, &post.Content)
		if err != nil {

		}
		return &post, nil
	} else {
		return nil, err
	}
}

func (p *postgreRepository) FindAll(ctx context.Context, tx *sql.Tx) []entity.Post {
	SQL := "SELECT title, content FROM post"
	rows, err := tx.QueryContext(ctx, SQL)
	if err != nil {

	}
	defer rows.Close()

	var posts []entity.Post
	for rows.Next() {
		post := entity.Post{}
		err := rows.Scan(&post.ID, &post.Title, &post.Content)
		if err != nil {

		}
		posts = append(posts, post)
	}

	return posts
}
