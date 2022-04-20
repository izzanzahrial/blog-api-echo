package posting

import (
	"context"
	"database/sql"
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/izzanzahrial/blog-api-echo/pkg/repository"
	"github.com/stretchr/testify/mock"
)

var (
	ErrPostIsntValidate          = errors.New("post data from handler isn't validate")
	ErrFailedToBeginTransaction  = errors.New("failed to begin transaction to the repository")
	ErrFailedToCommitTransaction = errors.New("failed to commit transaction to the repository")
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
}

func NewService(pr repository.PostDatabase, db txDB, val *validator.Validate) Service {
	return &service{
		Repository: pr,
		DB:         db,
		Validate:   val,
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
		return createdPost, err
	}

	if err := tx.Commit(); err != nil {
		return createdPost, ErrFailedToCommitTransaction
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
		return foundPost, err
	}

	updatedPost, err := ps.Repository.Update(ctx, tx, foundPost)
	if err != nil {
		return updatedPost, err
	}

	if err := tx.Commit(); err != nil {
		return updatedPost, ErrFailedToCommitTransaction
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

	return nil
}
func (ps *service) FindByID(ctx context.Context, postID uint64) (repository.Post, error) {
	tx, err := ps.DB.Begin()
	if err != nil {
		return repository.Post{}, ErrFailedToBeginTransaction
	}
	defer tx.Rollback()

	foundPost, err := ps.Repository.FindByID(ctx, tx, postID)
	if err != nil {
		return foundPost, err
	}

	if err := tx.Commit(); err != nil {
		return foundPost, ErrFailedToCommitTransaction
	}

	return foundPost, nil
}
func (ps *service) FindAll(ctx context.Context) ([]repository.Post, error) {
	tx, err := ps.DB.Begin()
	if err != nil {
		return []repository.Post{}, ErrFailedToBeginTransaction
	}
	defer tx.Rollback()

	posts, err := ps.Repository.FindAll(ctx, tx)
	if err != nil {
		return posts, err
	}

	if err := tx.Commit(); err != nil {
		return posts, ErrFailedToCommitTransaction
	}

	return posts, nil
}
