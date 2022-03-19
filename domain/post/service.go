package post

import (
	"context"
	"database/sql"
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/izzanzahrial/blog-api-echo/entity"
)

var (
	ErrPostIsntValidate          = errors.New("post data from handler isn't validate")
	ErrFailedToBeginTransaction  = errors.New("failed to begin transaction to the repository")
	ErrFailedToCommitTransaction = errors.New("failed to commit transaction to the repository")
)

type PostService interface {
	Create(ctx context.Context, post entity.Post) (entity.Post, error)
	Update(ctx context.Context, post entity.Post) (entity.Post, error)
	Delete(ctx context.Context, postID uint64) error
	FindByID(ctx context.Context, postID uint64) (entity.Post, error)
	FindAll(ctx context.Context) ([]entity.Post, error)
}

type postService struct {
	PostRepository PostRepository
	DB             *sql.DB
	Validate       *validator.Validate
}

func NewPostService(pr PostRepository, db *sql.DB, val *validator.Validate) PostService {
	return &postService{
		PostRepository: pr,
		DB:             db,
		Validate:       val,
	}
}

func (ps *postService) Create(ctx context.Context, post entity.Post) (entity.Post, error) {
	err := ps.Validate.Struct(post)
	if err != nil {
		return post, ErrPostIsntValidate
	}

	tx, err := ps.DB.Begin()
	if err != nil {
		return post, ErrFailedToBeginTransaction
	}
	defer tx.Rollback()

	post, err = ps.PostRepository.Create(ctx, tx, post)
	if err != nil {
		return post, err
	}

	if err := tx.Commit(); err != nil {
		return post, ErrFailedToCommitTransaction
	}

	return post, nil
}

func (ps *postService) Update(ctx context.Context, post entity.Post) (entity.Post, error) {
	err := ps.Validate.Struct(post)
	if err != nil {
		return post, ErrPostIsntValidate
	}

	tx, err := ps.DB.Begin()
	if err != nil {
		return post, ErrFailedToBeginTransaction
	}
	defer tx.Rollback()

	newPost, err := ps.PostRepository.FindByID(ctx, tx, post.ID)
	if err != nil {
		return post, err
	}

	newPost.Title = post.Title
	newPost.Content = post.Content

	post, err = ps.PostRepository.Update(ctx, tx, newPost)
	if err != nil {
		return post, err
	}

	if err := tx.Commit(); err != nil {
		return post, ErrFailedToCommitTransaction
	}

	return post, nil
}

func (ps *postService) Delete(ctx context.Context, postID uint64) error {
	tx, err := ps.DB.Begin()
	if err != nil {
		return ErrFailedToBeginTransaction
	}
	defer tx.Rollback()

	post, err := ps.PostRepository.FindByID(ctx, tx, postID)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return ErrFailedToCommitTransaction
	}

	ps.PostRepository.Delete(ctx, tx, post)

	return nil
}
func (ps *postService) FindByID(ctx context.Context, postID uint64) (entity.Post, error) {
	tx, err := ps.DB.Begin()
	if err != nil {
		return entity.Post{}, ErrFailedToBeginTransaction
	}
	defer tx.Rollback()

	post, err := ps.PostRepository.FindByID(ctx, tx, postID)
	if err != nil {
		return post, err
	}

	if err := tx.Commit(); err != nil {
		return post, ErrFailedToCommitTransaction
	}

	return post, nil
}
func (ps *postService) FindAll(ctx context.Context) ([]entity.Post, error) {
	tx, err := ps.DB.Begin()
	if err != nil {
		return []entity.Post{}, ErrFailedToBeginTransaction
	}
	defer tx.Rollback()

	posts, err := ps.PostRepository.FindAll(ctx, tx)
	if err != nil {
		return posts, err
	}

	if err := tx.Commit(); err != nil {
		return posts, ErrFailedToCommitTransaction
	}

	return posts, nil
}
