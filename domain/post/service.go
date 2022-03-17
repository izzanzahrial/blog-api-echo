package post

import (
	"context"
	"database/sql"

	"github.com/go-playground/validator/v10"
	"github.com/izzanzahrial/blog-api-echo/entity"
)

type PostService interface {
	Create(ctx context.Context, post entity.Post) entity.Post
	Update(ctx context.Context, post entity.Post) entity.Post
	Delete(ctx context.Context, postID uint64)
	FindByID(ctx context.Context, postID uint64) entity.Post
	FindAll(ctx context.Context) []entity.Post
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

func (ps *postService) Create(ctx context.Context, post entity.Post) entity.Post {
	err := ps.Validate.Struct(post)
	if err != nil {

	}

	tx, err := ps.DB.Begin()
	if err != nil {

	}

	post, err = ps.PostRepository.Create(ctx, tx, post)
	if err != nil {

	}

	return post
}

func (ps *postService) Update(ctx context.Context, post entity.Post) entity.Post {
	err := ps.Validate.Struct(post)
	if err != nil {

	}

	tx, err := ps.DB.Begin()
	if err != nil {

	}

	newPost, err := ps.PostRepository.FindByID(ctx, tx, post.ID)
	if err != nil {

	}

	newPost.Title = post.Title
	newPost.Content = post.Content

	post, err = ps.PostRepository.Update(ctx, tx, newPost)
	if err != nil {

	}

	return post
}

func (ps *postService) Delete(ctx context.Context, postID uint64) {
	tx, err := ps.DB.Begin()
	if err != nil {

	}

	post, err := ps.PostRepository.FindByID(ctx, tx, postID)
	if err != nil {

	}

	ps.PostRepository.Delete(ctx, tx, post)
}
func (ps *postService) FindByID(ctx context.Context, postID uint64) entity.Post {
	tx, err := ps.DB.Begin()
	if err != nil {

	}

	post, err := ps.PostRepository.FindByID(ctx, tx, postID)
	if err != nil {

	}

	return post
}
func (ps *postService) FindAll(ctx context.Context) []entity.Post {
	tx, err := ps.DB.Begin()
	if err != nil {

	}

	posts, err := ps.PostRepository.FindAll(ctx, tx)
	if err != nil {

	}

	return posts
}
