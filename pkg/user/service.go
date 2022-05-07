package user

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"github.com/izzanzahrial/blog-api-echo/pkg/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserIsntValidate          = errors.New("user data from handle isn't validate")
	ErrFailedToBeginTransaction  = errors.New("failed to begin transaction to the repository")
	ErrFailedToCommitTransaction = errors.New("failed to commit transaction to the repository")
	ErrFailedToGeneratePassword  = errors.New("failed to generate password")
	ErrUnauthorizedUser          = errors.New("unathorized user")
)

type UserService interface {
	Create(ctx context.Context, u User) (User, error)
	UpdateUser(ctx context.Context, u User) (User, error)
	UpdatePassword(ctx context.Context, u User) (User, error)
	Delete(ctx context.Context, id int64, pass string) error
	Login(ctx context.Context, id int64, pass string) (User, string, error)
}

type userService struct {
	UserRepository repository.UserRepository
	DB             *sql.DB
	Validate       *validator.Validate
}

func NewUserService(ur repository.UserRepository, db *sql.DB, val *validator.Validate) UserService {
	return &userService{
		UserRepository: ur,
		DB:             db,
		Validate:       val,
	}
}

func (us *userService) Create(ctx context.Context, u User) (User, error) {
	err := us.Validate.Struct(u)
	if err != nil {
		return User{}, ErrUserIsntValidate
	}

	tx, err := us.DB.Begin()
	if err != nil {
		return User{}, ErrFailedToBeginTransaction
	}
	defer tx.Rollback()

	repoUser, err = us.UserRepository.Create(ctx, tx, repository.User(u))
	if err != nil {
		return User{}, err
	}

	if err := tx.Commit(); err != nil {
		return User{}, ErrFailedToCommitTransaction
	}

	return repoUser, nil
}

func (us *userService) UpdateUser(ctx context.Context, u User) (User, error) {
	if err := us.Validate.Struct(u); err != nil {
		return User{}, ErrUserIsntValidate
	}

	tx, err := us.DB.Begin()
	if err != nil {
		return User{}, ErrFailedToBeginTransaction
	}
	defer tx.Rollback()

	oldUser, err := us.UserRepository.FindByID(ctx, tx, u.ID)
	if err != nil {
		return User{}, err
	}

	oldUser.Name = u.Name

	newUser, err = us.UserRepository.UpdateUser(ctx, tx, oldUser)
	if err != nil {
		return User{}, err
	}

	if err := tx.Commit(); err != nil {
		return User{}, ErrFailedToCommitTransaction
	}

	return newUser, nil
}

func (us *userService) UpdatePassword(ctx context.Context, u User) (User, error) {
	if err := us.Validate.Struct(u); err != nil {
		return User{}, ErrUserIsntValidate
	}

	tx, err := us.DB.Begin()
	if err != nil {
		return User{}, ErrFailedToBeginTransaction
	}
	defer tx.Rollback()

	oldUser, err := us.UserRepository.Login(ctx, tx, u.ID, u.Password)
	if err != nil {
		return User{}, err
	}

	oldUser.Password = u.Password

	newUser, err = us.UserRepository.UpdatePassword(ctx, tx, oldUser)
	if err != nil {
		return User{}, err
	}

	if err := tx.Commit(); err != nil {
		return User{}, ErrFailedToCommitTransaction
	}

	return newUser, nil
}

func (us *userService) Delete(ctx context.Context, id int64, pass string) error {
	tx, err := us.DB.Begin()
	if err != nil {
		return ErrFailedToBeginTransaction
	}
	defer tx.Rollback()

	user, err := us.UserRepository.Login(ctx, tx, id, pass)
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

func (us *userService) Login(ctx context.Context, id int64, pass string) (User, string, error) {
	tx, err := us.DB.Begin()
	if err != nil {
		return User{}, "", ErrFailedToBeginTransaction
	}
	defer tx.Rollback()

	hashedPassword, err := hashPassword(pass)
	if err != nil {
		return User{}, "", err
	}

	user, err := us.UserRepository.Login(ctx, tx, id, hashedPassword)
	if err != nil {
		return User{}, "", err
	}

	if ok := CheckPasswordHash(hashedPassword, user.Password); !ok {
		return User{}, "", ErrUnauthorizedUser
	}

	if err := tx.Commit(); err != nil {
		return User{}, "", ErrFailedToCommitTransaction
	}

	id64 := uint64(user.ID)

	token, err := createJWTToken(id64, false)
	if err != nil {
		return User{}, token, err
	}

	return user, token, nil
}

type JWTClaims struct {
	UserID uint64
	Admin  bool
	jwt.StandardClaims
}

func createJWTToken(userID uint64, admin bool) (string, error) {
	claims := JWTClaims{
		userID,
		admin,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 3600).Unix(),
		},
	}

	rawToken := jwt.NewWithClaims(jwt.SigningMethodES512, claims)

	token, err := rawToken.SignedString([]byte("izzan"))
	if err != nil {
		return "", err
	}

	return token, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", ErrFailedToGeneratePassword
	}

	return string(bytes), nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
