package posting

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
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
	Create(ctx context.Context, post repository.Post) (repository.Post, error)
	Update(ctx context.Context, post repository.Post) (repository.Post, error)
	Delete(ctx context.Context, postID uint64) error
	FindByID(ctx context.Context, postID uint64) (repository.Post, error)
	FindAll(ctx context.Context) ([]repository.Post, error)
}

type MockService struct {
	mock.Mock
}

func (m *MockService) Create(ctx context.Context, post repository.Post) (repository.Post, error) {
	args := m.Called(ctx, post)
	return args.Get(0).(repository.Post), args.Error(1)
}

func (m *MockService) Update(ctx context.Context, post repository.Post) (repository.Post, error) {
	args := m.Called(ctx, post)
	return args.Get(0).(repository.Post), args.Error(1)
}

func (m *MockService) Delete(ctx context.Context, postID uint64) error {
	args := m.Called(ctx, postID)
	return args.Error(0)
}

func (m *MockService) FindByID(ctx context.Context, postID uint64) (repository.Post, error) {
	args := m.Called(ctx, postID)
	return args.Get(0).(repository.Post), args.Error(1)
}

func (m *MockService) FindAll(ctx context.Context) ([]repository.Post, error) {
	args := m.Called(ctx)
	return args.Get(0).([]repository.Post), args.Error(1)
}

// naming things is hard
type txDB interface {
	Begin() (*sql.Tx, error)
}

type MockDB struct {
	mock.Mock
}

func (m *MockDB) Begin() (*sql.Tx, error) {
	args := m.Called()
	return args.Get(0).(*sql.Tx), args.Error(1)
}

type service struct {
	Repository repository.PostDatabase
	DB         txDB
	Validate   *validator.Validate
	Rdb        *redis.Client
}

func NewService(pr repository.PostDatabase, db txDB, val *validator.Validate, rdb *redis.Client) Service {
	return &service{
		Repository: pr,
		DB:         db,
		Validate:   val,
		Rdb:        rdb,
	}
}

func (ps *service) Create(ctx context.Context, post repository.Post) (repository.Post, error) {
	err := ps.Validate.Struct(post)
	if err != nil {
		return post, ErrPostIsntValidate
	}

	tx, err := ps.DB.Begin()
	if err != nil {
		return post, ErrFailedToBeginTransaction
	}
	defer tx.Rollback()

	createdPost, err := ps.Repository.Create(ctx, tx, post)
	if err != nil {
		return repository.Post{}, err
	}

	if err := tx.Commit(); err != nil {
		return repository.Post{}, ErrFailedToCommitTransaction
	}

	ttl := time.Duration(3600) * time.Second
	id := strconv.FormatUint(createdPost.ID, 10)
	str := strings.Builder{}
	str.WriteString("post")
	str.WriteString(id)
	op1 := ps.Rdb.Set(context.Background(), str.String(), createdPost, ttl)
	if err := op1.Err(); err != nil {
		return createdPost, ErrFailedToCachePost
	}

	return createdPost, nil
}

func (ps *service) Update(ctx context.Context, post repository.Post) (repository.Post, error) {
	err := ps.Validate.Struct(post)
	if err != nil {
		return post, ErrPostIsntValidate
	}

	tx, err := ps.DB.Begin()
	if err != nil {
		return post, ErrFailedToBeginTransaction
	}
	defer tx.Rollback()

	foundPost, err := ps.Repository.FindByID(ctx, tx, post.ID)
	if err != nil {
		return repository.Post{}, err
	}

	updatedPost, err := ps.Repository.Update(ctx, tx, foundPost)
	if err != nil {
		return repository.Post{}, err
	}

	if err := tx.Commit(); err != nil {
		return repository.Post{}, ErrFailedToCommitTransaction
	}

	ttl := time.Duration(3600) * time.Second
	id := strconv.FormatUint(updatedPost.ID, 10)
	str := strings.Builder{}
	str.WriteString("post")
	str.WriteString(id)
	op1 := ps.Rdb.Set(context.Background(), str.String(), updatedPost, ttl)
	if err := op1.Err(); err != nil {
		return updatedPost, ErrFailedToCachePost
	}

	return updatedPost, nil
}

func (ps *service) Delete(ctx context.Context, postID uint64) error {
	tx, err := ps.DB.Begin()
	if err != nil {
		return ErrFailedToBeginTransaction
	}
	defer tx.Rollback()

	foundPost, err := ps.Repository.FindByID(ctx, tx, postID)
	if err != nil {
		return err
	}

	if err := ps.Repository.Delete(ctx, tx, foundPost); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return ErrFailedToCommitTransaction
	}

	id := strconv.FormatUint(postID, 10)
	str := strings.Builder{}
	str.WriteString("post")
	str.WriteString(id)
	op1 := ps.Rdb.Del(context.Background(), str.String())
	if err := op1.Err(); err != nil {
		return err
	}

	return nil
}
func (ps *service) FindByID(ctx context.Context, postID uint64) (repository.Post, error) {
	id := strconv.FormatUint(postID, 10)
	str := strings.Builder{}
	str.WriteString("post")
	str.WriteString(id)
	val, err := ps.Rdb.Get(context.Background(), str.String()).Result()
	if err == nil {
		post := repository.Post{}
		json.Unmarshal([]byte(val), &post)
		return post, nil
	}

	tx, err := ps.DB.Begin()
	if err != nil {
		return repository.Post{}, ErrFailedToBeginTransaction
	}
	defer tx.Rollback()

	foundPost, err := ps.Repository.FindByID(ctx, tx, postID)
	if err != nil {
		return repository.Post{}, err
	}

	if err := tx.Commit(); err != nil {
		return repository.Post{}, ErrFailedToCommitTransaction
	}

	return foundPost, nil
}
func (ps *service) FindAll(ctx context.Context) ([]repository.Post, error) {
	val, err := ps.Rdb.Keys(context.Background(), "post*").Result()
	if err == nil {
		posts := []repository.Post{}
		post := repository.Post{}
		for _, n := range val {
			json.Unmarshal([]byte(n), &post)
			posts = append(posts, post)
		}
		return posts, nil
	}

	tx, err := ps.DB.Begin()
	if err != nil {
		return []repository.Post{}, ErrFailedToBeginTransaction
	}
	defer tx.Rollback()

	posts, err := ps.Repository.FindAll(ctx, tx)
	if err != nil {
		return []repository.Post{}, err
	}

	if err := tx.Commit(); err != nil {
		return []repository.Post{}, ErrFailedToCommitTransaction
	}

	return posts, nil
}
