package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/stretchr/testify/mock"
)

type MockPostingPostgre struct {
	mock.Mock
}

func (m *MockPostingPostgre) Create(ctx context.Context, tx *sql.Tx, pd PostData) (PostData, error) {
	args := m.Called(ctx, tx, pd)
	return args.Get(0).(PostData), args.Error(1)
}

func (m *MockPostingPostgre) Update(ctx context.Context, tx *sql.Tx, pd PostData) error {
	args := m.Called(ctx, tx, pd)
	return args.Error(0)
}

func (m *MockPostingPostgre) Delete(ctx context.Context, tx *sql.Tx, pd PostData) error {
	args := m.Called(ctx, tx, pd)
	return args.Error(0)
}

func (m *MockPostingPostgre) FindByID(ctx context.Context, tx *sql.Tx, id int64) (PostData, error) {
	args := m.Called(ctx, tx, id)
	return args.Get(0).(PostData), args.Error(1)
}

func (m *MockPostingPostgre) FindByTitleContent(ctx context.Context, tx *sql.Tx, query string, from int, size int) ([]PostData, error) {
	args := m.Called(ctx, tx, query, from, size)
	return args.Get(0).([]PostData), args.Error(1)
}

func (m *MockPostingPostgre) FindRecent(ctx context.Context, tx *sql.Tx, from int, size int) ([]PostData, error) {
	args := m.Called(ctx, tx, from, size)
	return args.Get(0).([]PostData), args.Error(1)
}

type postingPostgre struct {
}

func NewPostgre() Post {
	return &postingPostgre{}
}

func (p *postingPostgre) Create(ctx context.Context, tx *sql.Tx, pd PostData) (PostData, error) {
	SQL := "INSERT INTO post(title, short_desc, content) VALUES (?, ?, ?)"
	result, err := tx.ExecContext(ctx, SQL, pd.Title, pd.ShortDesc, pd.Content, pd.CreatedAt)
	if err != nil {
		return pd, fmt.Errorf("failed to created post: %v, because %w", pd, err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return pd, fmt.Errorf("failed post cannot be found: %v, because %w", pd, err)
	}

	pd.ID = id

	return pd, nil
}

func (p *postingPostgre) Update(ctx context.Context, tx *sql.Tx, pd PostData) error {
	SQL := "UPDATE post SET title = ?, short_desc = ?, content = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, SQL, pd.Title, pd.ShortDesc, pd.Content, pd.ID)
	if err != nil {
		return fmt.Errorf("failed to update post: %v, because %w", pd, err)
	}

	return nil
}

func (p *postingPostgre) Delete(ctx context.Context, tx *sql.Tx, pd PostData) error {
	SQL := "DELETE FROM post WHERE id = ?"
	_, err := tx.ExecContext(ctx, SQL, pd.ID)
	if err != nil {
		return fmt.Errorf("failed to delete post: %v because %w", pd, err)
	}

	return nil
}

func (p *postingPostgre) FindByID(ctx context.Context, tx *sql.Tx, id int64) (PostData, error) {
	SQL := "SELECT title, short_desc, content, created_at FROM post WHERE id = ?"
	rows, err := tx.QueryContext(ctx, SQL, id)
	if err != nil {
		return PostData{}, fmt.Errorf("failed to find post with id: %d because %w", id, err)
	}
	defer rows.Close()

	var post PostData
	if rows.Next() {
		if err := rows.Scan(&post.ID, &post.Title, &post.ShortDesc, &post.Content, &post.CreatedAt); err != nil {
			return post, fmt.Errorf("failed to scan post with id: %d because %w", id, err)
		}
		return post, nil
	} else {
		return post, fmt.Errorf("failed to find post with id: %d because %w", id, err)
	}
}

func (p *postingPostgre) FindByTitleContent(ctx context.Context, tx *sql.Tx, query string, from int, size int) ([]PostData, error) {
	// full text search postgres https://blog.crunchydata.com/blog/postgres-full-text-search-a-search-engine-in-a-database
	selectFrom := "SELECT id, title, short_desc, content, created_at FROM posts"
	condition := "WHERE ts_title_content @@ to_tsquery('english', '?')"
	orderBy := "ORDER BY ts_rank(ts_title_content, to_tsquery('english', '?')) LIMIT ? OFFSET ? DESC"
	SQL := selectFrom + " " + condition + " " + orderBy
	rows, err := tx.QueryContext(ctx, SQL, query, query, size, from)
	if err != nil {
		return []PostData{}, fmt.Errorf("failed to find post with keywords: %s because %w", query, err)
	}
	defer rows.Close()

	var result []PostData
	for rows.Next() {
		var post PostData
		if err := rows.Scan(&post.ID, &post.Title, &post.ShortDesc, &post.Content, &post.CreatedAt); err != nil {
			return []PostData{}, fmt.Errorf("failed to find scan post with keywords: %s because %w", query, err)
		}
		result = append(result, post)
	}

	return result, nil
}

func (p *postingPostgre) FindRecent(ctx context.Context, tx *sql.Tx, from int, size int) ([]PostData, error) {
	SQL := "SELECT id, title, short_desc, created_at content FROM post LIMIT ? OFFSET ?"
	rows, err := tx.QueryContext(ctx, SQL, size, from)
	if err != nil {
		return []PostData{}, fmt.Errorf("failed to find posts because %w", err)
	}
	defer rows.Close()

	var posts []PostData
	for rows.Next() {
		var post PostData
		if err := rows.Scan(&post.ID, &post.Title, &post.ShortDesc, &post.Content, &post.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan post becasue %w", err)
		}
		posts = append(posts, post)
	}

	return posts, nil
}
