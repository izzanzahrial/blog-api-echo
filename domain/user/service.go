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
	UpdateUser(ctx context.Context, user entity.User) (entity.User, error)
	UpdatePassword(ctx context.Context, user entity.User) (entity.User, error)
	Delete(ctx context.Context, id uint64, pass string) error
	Login(ctx context.Context, id uint64, pass string) (entity.User, error)
}

type userService struct {
	UserRepository UserRepository
	DB             *sql.DB
	Validate       *validator.Validate
}

func NewUserService(ur UserRepository, db *sql.DB, val *validator.Validate) UserService {
	return &userService{
		UserRepository: ur,
		DB:             db,
		Validate:       val,
	}
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

func (us *userService) UpdateUser(ctx context.Context, user entity.User) (entity.User, error) {
	if err := us.Validate.Struct(user); err != nil {
		return user, ErrUserIsntValidate
	}

	tx, err := us.DB.Begin()
	if err != nil {
		return user, ErrFailedToBeginTransaction
	}
	defer tx.Rollback()

	oldUser, err := us.UserRepository.FindByID(ctx, tx, user.ID)
	if err != nil {
		return user, err
	}

	oldUser.Name = user.Name

	user, err = us.UserRepository.UpdateUser(ctx, tx, oldUser)
	if err != nil {
		return user, err
	}

	if err := tx.Commit(); err != nil {
		return user, ErrFailedToCommitTransaction
	}

	return user, nil

}

func (us *userService) UpdatePassword(ctx context.Context, user entity.User) (entity.User, error) {
	if err := us.Validate.Struct(user); err != nil {
		return user, ErrUserIsntValidate
	}

	tx, err := us.DB.Begin()
	if err != nil {
		return user, ErrFailedToBeginTransaction
	}
	defer tx.Rollback()

	oldUser, err := us.UserRepository.Find(ctx, tx, user.ID, user.Password)
	if err != nil {
		return user, err
	}

	oldUser.Password = user.Password

	user, err = us.UserRepository.UpdatePassword(ctx, tx, oldUser)
	if err != nil {
		return user, err
	}

	if err := tx.Commit(); err != nil {
		return user, ErrFailedToCommitTransaction
	}

	return user, nil
}

func (us *userService) Delete(ctx context.Context, id uint64, pass string) error {
	tx, err := us.DB.Begin()
	if err != nil {
		return ErrFailedToBeginTransaction
	}
	defer tx.Rollback()

	user, err := us.UserRepository.Find(ctx, tx, id, pass)
	if err != nil {
		return err
	}

	if err := us.UserRepository.Delete(ctx, tx, user); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return ErrFailedToCommitTransaction
	}

	return nil
}

func (us *userService) Login(ctx context.Context, id uint64, pass string) (entity.User, error) {
	tx, err := us.DB.Begin()
	if err != nil {
		return entity.User{}, ErrFailedToBeginTransaction
	}
	defer tx.Rollback()

	user, err := us.UserRepository.Find(ctx, tx, id, pass)
	if err != nil {
		return user, err
	}

	if err := tx.Commit(); err != nil {
		return user, ErrFailedToCommitTransaction
	}

	return user, nil
}
