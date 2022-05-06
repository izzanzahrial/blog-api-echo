package posting

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/izzanzahrial/blog-api-echo/pkg/elastic"
	caching "github.com/izzanzahrial/blog-api-echo/pkg/redis"
	"github.com/izzanzahrial/blog-api-echo/pkg/repository"
	"github.com/stretchr/testify/mock"
)

var (
	ErrPostIsntValidate          = errors.New("post data from handler isn't validate")
	ErrFailedToBeginTransaction  = errors.New("failed to begin transaction to the repository")
	ErrFailedToCommitTransaction = errors.New("failed to commit transaction to the repository")
	ErrFailedToCachePost         = errors.New("failed to cache post")
)

type Service interface {
	Create(ctx context.Context, post PostData) (repository.PostData, error)
	Update(ctx context.Context, post repository.PostData) error
	Delete(ctx context.Context, id int64) error
	FindByID(ctx context.Context, id int64) (repository.PostData, error)
	FindByTitleContent(ctx context.Context, query string, from int, size int) ([]repository.PostData, error)
	FindRecent(ctx context.Context, from int, size int) ([]repository.PostData, error)
}

type MockService struct {
	mock.Mock
}

func (m *MockService) Create(ctx context.Context, post PostData) (repository.PostData, error) {
	args := m.Called(ctx, post)
	return args.Get(0).(repository.PostData), args.Error(1)
}

func (m *MockService) Update(ctx context.Context, post repository.PostData) error {
	args := m.Called(ctx, post)
	return args.Error(0)
}

func (m *MockService) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockService) FindByID(ctx context.Context, id int64) (repository.PostData, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(repository.PostData), args.Error(1)
}

func (m *MockService) FindRecent(ctx context.Context, from int, size int) ([]repository.PostData, error) {
	args := m.Called(ctx, from, size)
	return args.Get(0).([]repository.PostData), args.Error(1)
}

func (m *MockService) FindByTitleContent(ctx context.Context, query string, from int, size int) ([]repository.PostData, error) {
	args := m.Called(ctx, query, from, size)
	return args.Get(0).([]repository.PostData), args.Error(1)
}

// type transaction interface {
// 	Rollback() error
// 	Commit() error
// }

// naming things is hard
type DBtx interface {
	Begin() (*sql.Tx, error)
}

type MockDBtx struct {
	mock.Mock
}

func (m *MockDBtx) Begin() (*sql.Tx, error) {
	args := m.Called()
	return args.Get(0).(*sql.Tx), args.Error(1)
}

type service struct {
	Repository repository.Post
	DB         DBtx
	Validate   *validator.Validate
	Cache      caching.Cache
	Es         elastic.ElasticDB
}

func NewService(rp repository.Post, db DBtx, val *validator.Validate, cache caching.Cache, es elastic.ElasticDB) Service {
	return &service{
		Repository: rp,
		DB:         db,
		Validate:   val,
		Cache:      cache,
		Es:         es,
	}
}

func (ps *service) Create(ctx context.Context, post PostData) (repository.PostData, error) {
	err := ps.Validate.Struct(post)
	if err != nil {
		return repository.PostData{}, fmt.Errorf("failed to validate: %v because %w", post, err)
	}

	tx, err := ps.DB.Begin()
	if err != nil {
		return repository.PostData{}, fmt.Errorf("failed to begin transaction for: %v because %w", post, err)
	}
	defer tx.Rollback()

	postData := repository.PostData{
		Title:     post.Title,
		ShortDesc: post.ShortDesc,
		Content:   post.Content,
		CreatedAt: time.Now(),
	}

	createdPost, err := ps.Repository.Create(ctx, tx, postData)
	if err != nil {
		return repository.PostData{}, err
	}

	if err := tx.Commit(); err != nil {
		return repository.PostData{}, fmt.Errorf("failed to commit transaction: %v because %w", createdPost, err)
	}

	if err = ps.Es.Insert(ctx, createdPost); err != nil {
		return repository.PostData{}, fmt.Errorf("failed to insert data: %v to elasticsearch because %w", createdPost, err)
	}

	ttl := time.Duration(3600) * time.Second
	strID := strconv.Itoa(int(createdPost.ID))

	str := strings.Builder{}
	str.WriteString("post")
	str.WriteString(strID)

	op1 := ps.Cache.Set(context.Background(), str.String(), createdPost, ttl)
	if err := op1.Err(); err != nil {
		return createdPost, fmt.Errorf("failed to cache post: %v because %w", createdPost, err)
	}

	return createdPost, nil
}

func (ps *service) Update(ctx context.Context, post repository.PostData) error {
	err := ps.Validate.Struct(post)
	if err != nil {
		return fmt.Errorf("failed to validate: %v because %w", post, err)
	}

	tx, err := ps.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction for: %v because %w", post, err)
	}
	defer tx.Rollback()

	foundPost, err := ps.Repository.FindByID(ctx, tx, post.ID)
	if err != nil {
		return err
	}

	if err := ps.Repository.Update(ctx, tx, post); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transcation: %v because %w", post, err)
	}

	if err = ps.Es.Update(ctx, post); err != nil {
		return fmt.Errorf("failed to update data: %v from elasticsearch because %w", post, err)
	}

	ttl := time.Duration(3600) * time.Second
	strID := strconv.Itoa(int(foundPost.ID))

	str := strings.Builder{}
	str.WriteString("post")
	str.WriteString(strID)

	op1 := ps.Cache.Set(context.Background(), str.String(), post, ttl)
	if err := op1.Err(); err != nil {
		return fmt.Errorf("failed to cache post: %v because %w", post, err)
	}

	return nil
}

func (ps *service) Delete(ctx context.Context, id int64) error {
	tx, err := ps.DB.Begin()
	if err != nil {
		return ErrFailedToBeginTransaction
	}
	defer tx.Rollback()

	foundPost, err := ps.Repository.FindByID(ctx, tx, id)
	if err != nil {
		return err
	}

	if err := ps.Repository.Delete(ctx, tx, foundPost); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %d because %w", id, err)
	}

	strID := strconv.Itoa(int(id))

	if err = ps.Es.Delete(ctx, strID); err != nil {
		return fmt.Errorf("failed to delete data: %d from elasticsearch because %w", id, err)
	}

	str := strings.Builder{}
	str.WriteString("post")
	str.WriteString(strID)

	op1 := ps.Cache.Del(context.Background(), str.String())
	if err := op1.Err(); err != nil {
		return err
	}

	return nil
}
func (ps *service) FindByID(ctx context.Context, id int64) (repository.PostData, error) {
	strID := strconv.Itoa(int(id))

	str := strings.Builder{}
	str.WriteString("post")
	str.WriteString(strID)

	val, err := ps.Cache.Get(context.Background(), str.String()).Result()
	if err == nil {
		var post repository.PostData
		json.Unmarshal([]byte(val), &post)
		return post, nil
	}

	foundPost, err := ps.Es.FindByID(ctx, strID)
	if err == nil {
		return foundPost, nil
	}

	tx, err := ps.DB.Begin()
	if err != nil {
		return repository.PostData{}, fmt.Errorf("failed to begin transaction: %d because %w", id, err)
	}
	defer tx.Rollback()

	foundPost, err = ps.Repository.FindByID(ctx, tx, id)
	if err != nil {
		return repository.PostData{}, err
	}

	if err := tx.Commit(); err != nil {
		return repository.PostData{}, fmt.Errorf("failed to commit transaction: %d because %w", id, err)
	}

	return foundPost, nil
}

func (ps *service) FindByTitleContent(ctx context.Context, query string, from int, size int) ([]repository.PostData, error) {
	var posts []repository.PostData

	val, err := ps.Cache.Get(ctx, query).Result()
	if err == nil {
		json.Unmarshal([]byte(val), &posts)
		return posts, nil
	}

	foundPosts, err := ps.Es.FindByTitleContent(ctx, query, from, size)
	if err == nil {
		for _, doc := range foundPosts.Hits {
			var post repository.PostData
			post.ID = int64(doc.ID)
			post.Title = doc.Title
			post.Content = doc.Content

			posts = append(posts, post)
		}

		ttl := time.Duration(3600) * time.Second
		op1 := ps.Cache.Set(ctx, query, posts, ttl)
		if err := op1.Err(); err != nil {
			return posts, fmt.Errorf("failed to cache: %v because %w", posts, err)
		}

		return posts, nil
	}

	tx, err := ps.DB.Begin()
	if err != nil {
		return posts, fmt.Errorf("failed to begin transaction for query: %v because %w", query, err)
	}
	defer tx.Rollback()

	posts, err = ps.Repository.FindByTitleContent(ctx, tx, query, from, size)
	if err != nil {
		return posts, err
	}

	if err := tx.Commit(); err != nil {
		return []repository.PostData{}, fmt.Errorf("failed to commit transaction for query: %v because %w", query, err)
	}

	return posts, nil
}

func (ps *service) FindRecent(ctx context.Context, from int, size int) ([]repository.PostData, error) {
	val, err := ps.Cache.Keys(context.Background(), "post*").Result()
	if err == nil {
		var posts []repository.PostData
		var post repository.PostData
		for _, n := range val {
			json.Unmarshal([]byte(n), &post)
			posts = append(posts, post)
		}
		return posts, nil
	}

	tx, err := ps.DB.Begin()
	if err != nil {
		return []repository.PostData{}, fmt.Errorf("failed to begin transaction for finding recent post because: %w", err)
	}
	defer tx.Rollback()

	posts, err := ps.Repository.FindRecent(ctx, tx, from, size)
	if err != nil {
		return []repository.PostData{}, err
	}

	if err := tx.Commit(); err != nil {
		return []repository.PostData{}, fmt.Errorf("failed to commit transaction for finding recent post because: %w", err)
	}

	return posts, nil
}
