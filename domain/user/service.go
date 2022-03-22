package user

import (
	"context"
	"database/sql"
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/izzanzahrial/blog-api-echo/entity"
)

var (
	ErrUserIsntValidate          = errors.New("user data from handle isn't validate")
	ErrFailedToBeginTransaction  = errors.New("failed to begin transaction to the repository")
	ErrFailedToCommitTransaction = errors.New("failed to commit transaction to the repository")
)

type UserService interface {
	Create(ctx context.Context, user entity.User) (entity.User, error)
	Update(ctx context.Context, user entity.User) (entity.User, error)
	Delete(ctx context.Context, id uint64) error
	Find(ctx context.Context, id uint64, password string) (entity.User, error)
}

type userService struct {
	UserRepository UserRepository
	DB             *sql.DB
	Validate       *validator.Validate
}

func (us *userService) Create(ctx context.Context, user entity.User) (entity.User, error) {
	err := us.Validate.Struct(user)
	if err != nil {
		return user, ErrUserIsntValidate
	}

	tx, err := us.DB.Begin()
	if err != nil {
		return user, ErrFailedToBeginTransaction
	}
	defer tx.Rollback()

	user, err = us.UserRepository.Create(ctx, tx, user)
	if err != nil {
		return user, err
	}

	if err := tx.Commit(); err != nil {
		return user, ErrFailedToCommitTransaction
	}

	return user, nil
}

func (us *userService) Update(ctx context.Context, user entity.User) (entity.User, error) {
}

func (us *userService) Delete(ctx context.Context, id uint64) error {

}

func (us *userService) Find(ctx context.Context, id uint64, password string) (entity.User, error) {
}
