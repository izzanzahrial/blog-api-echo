package repository

import (
	"context"
	"database/sql"

	"github.com/stretchr/testify/mock"
)

type MockPostgre struct {
	mock.Mock
}

func (m *MockPostgre) Create(ctx context.Context, tx *sql.Tx, ps Post) (Post, error) {
	args := m.Called(ctx, tx, ps)
	return args.Get(0).(Post), args.Error(1)
}
func (m *MockPostgre) Update(ctx context.Context, tx *sql.Tx, ps Post) (Post, error) {
	args := m.Called(ctx, tx, ps)
	return args.Get(0).(Post), args.Error(1)
}
func (m *MockPostgre) Delete(ctx context.Context, tx *sql.Tx, ps Post) error {
	args := m.Called(ctx, tx, ps)
	return args.Error(0)
}
func (m *MockPostgre) FindByID(ctx context.Context, tx *sql.Tx, ID uint64) (Post, error) {
	args := m.Called(ctx, tx, ID)
	return args.Get(0).(Post), args.Error(1)
}
func (m *MockPostgre) FindAll(ctx context.Context, tx *sql.Tx) ([]Post, error) {
	args := m.Called(ctx, tx)
	return args.Get(0).([]Post), args.Error(1)
}

type postgre struct {
}

func NewPostgre() PostDatabase {
	return &postgre{}
}

func (p *postgre) Create(ctx context.Context, tx *sql.Tx, ps Post) (Post, error) {
	SQL := "INSERT INTO post(title, content) VALUES (?, ?)"
	result, err := tx.ExecContext(ctx, SQL, ps.Title, ps.Content)
	if err != nil {
		return ps, ErrFailedToCreatePost
	}

	id, err := result.LastInsertId()
	if err != nil {
		return ps, ErrPostNotFound
	}

	ps.ID = uint64(id)

	return ps, nil
}

func (p *postgre) Update(ctx context.Context, tx *sql.Tx, ps Post) (Post, error) {
	SQL := "UPDATE post SET title = ?, content = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, SQL, ps.Title, ps.Content, ps.ID)
	if err != nil {
		return ps, ErrFailedUpdatePost
	}

	return ps, nil
}

func (p *postgre) Delete(ctx context.Context, tx *sql.Tx, ps Post) error {
	SQL := "DELETE FROM post WHERE id = ?"
	_, err := tx.ExecContext(ctx, SQL, ps.ID)
	if err != nil {
		return ErrFailedToDeletePost
	}

	return nil
}

func (p *postgre) FindByID(ctx context.Context, tx *sql.Tx, ID uint64) (Post, error) {
	SQL := "SELECT title, content FROM post WHERE id = ?"
	rows, err := tx.QueryContext(ctx, SQL, ID)
	if err != nil {
		return Post{}, ErrPostNotFound
	}
	defer rows.Close()

	post := Post{}
	if rows.Next() {
		if err := rows.Scan(&post.ID, &post.Title, &post.Content); err != nil {
			return post, ErrFailedToScanPost
		}
		return post, nil
	} else {
		return post, ErrPostNotFound
	}
}

func (p *postgre) FindAll(ctx context.Context, tx *sql.Tx) ([]Post, error) {
	SQL := "SELECT id, title, content FROM post"
	rows, err := tx.QueryContext(ctx, SQL)
	if err != nil {
		return nil, ErrPostNotFound
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		post := Post{}
		if err := rows.Scan(&post.ID, &post.Title, &post.Content); err != nil {
			return nil, ErrFailedToScanPost
		}
		posts = append(posts, post)
	}

	return posts, nil
}
