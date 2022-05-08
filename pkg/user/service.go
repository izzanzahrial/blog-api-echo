package user

import (
	"context"
	"database/sql"
	"errors"
	"net/mail"
	"os"
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
	Create(ctx context.Context, u User) (repository.User, error)
	UpdateUser(ctx context.Context, u repository.User) (repository.User, error)
	UpdatePassword(ctx context.Context, u repository.User, newPass string) (repository.User, error)
	Delete(ctx context.Context, u repository.User) error
	Login(ctx context.Context, emailOrUname string, pass string) (repository.User, string, error)
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

func (us *userService) Create(ctx context.Context, u User) (repository.User, error) {
	err := us.Validate.Struct(u)
	if err != nil {
		return repository.User{}, ErrUserIsntValidate
	}

	tx, err := us.DB.Begin()
	if err != nil {
		return repository.User{}, ErrFailedToBeginTransaction
	}
	defer tx.Rollback()

	user := repository.User{
		Email:    u.Email,
		Username: u.Username,
		Name:     u.Name,
		Password: u.Password,
	}

	user, err = us.UserRepository.Create(ctx, tx, user)
	if err != nil {
		return repository.User{}, err
	}

	if err := tx.Commit(); err != nil {
		return repository.User{}, ErrFailedToCommitTransaction
	}

	return user, nil
}

func (us *userService) UpdateUser(ctx context.Context, u repository.User) (repository.User, error) {
	if err := us.Validate.Struct(u); err != nil {
		return repository.User{}, ErrUserIsntValidate
	}

	tx, err := us.DB.Begin()
	if err != nil {
		return repository.User{}, ErrFailedToBeginTransaction
	}
	defer tx.Rollback()

	user, err := us.UserRepository.FindByEmail(ctx, tx, u.Email)
	if err != nil {
		return repository.User{}, err
	}

	user.Email = u.Email
	user.Username = u.Username
	user.Name = u.Name

	newUser, err := us.UserRepository.UpdateUser(ctx, tx, user)
	if err != nil {
		return repository.User{}, err
	}

	if err := tx.Commit(); err != nil {
		return repository.User{}, ErrFailedToCommitTransaction
	}

	return newUser, nil
}

func (us *userService) UpdatePassword(ctx context.Context, u repository.User, newPass string) (repository.User, error) {
	if err := us.Validate.Struct(u); err != nil {
		return repository.User{}, ErrUserIsntValidate
	}

	tx, err := us.DB.Begin()
	if err != nil {
		return repository.User{}, ErrFailedToBeginTransaction
	}
	defer tx.Rollback()

	user, err := us.UserRepository.FindByEmail(ctx, tx, u.Email)
	if err != nil {
		return repository.User{}, err
	}

	if ok := CheckPasswordHash(u.Password, user.Password); !ok {
		return repository.User{}, bcrypt.ErrMismatchedHashAndPassword
	}

	hashPass, err := hashPassword(newPass)
	if err != nil {
		return repository.User{}, bcrypt.ErrHashTooShort
	}

	user.Password = hashPass

	newUser, err := us.UserRepository.UpdatePassword(ctx, tx, user)
	if err != nil {
		return repository.User{}, err
	}

	if err := tx.Commit(); err != nil {
		return repository.User{}, ErrFailedToCommitTransaction
	}

	return newUser, nil
}

func (us *userService) Delete(ctx context.Context, u repository.User) error {
	tx, err := us.DB.Begin()
	if err != nil {
		return ErrFailedToBeginTransaction
	}
	defer tx.Rollback()

	user, err := us.UserRepository.LoginByEmail(ctx, tx, u.Email, u.Password)
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

func (us *userService) Login(ctx context.Context, emailOrUname string, pass string) (repository.User, string, error) {
	tx, err := us.DB.Begin()
	if err != nil {
		return repository.User{}, "", ErrFailedToBeginTransaction
	}
	defer tx.Rollback()

	hashPass, err := hashPassword(pass)
	if err != nil {
		return repository.User{}, "", err
	}

	var user repository.User
	if ok := validateEmail(emailOrUname); !ok {
		user, err = us.UserRepository.LoginByUsername(ctx, tx, emailOrUname, hashPass)
		if err != nil {
			return repository.User{}, "", err
		}
	} else {
		user, err = us.UserRepository.LoginByEmail(ctx, tx, emailOrUname, hashPass)
		if err != nil {
			return repository.User{}, "", err
		}
	}

	if ok := CheckPasswordHash(hashPass, user.Password); !ok {
		return repository.User{}, "", ErrUnauthorizedUser
	}

	if err := tx.Commit(); err != nil {
		return repository.User{}, "", ErrFailedToCommitTransaction
	}

	token, err := createJWTToken(false, user.ID, user.Email, user.Username, user.Name)
	if err != nil {
		return repository.User{}, token, err
	}

	return user, token, nil
}

type JWTClaims struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Admin    bool   `json:"admin"`
	jwt.StandardClaims
}

func createJWTToken(admin bool, id int64, email, username, name string) (string, error) {
	claims := JWTClaims{
		id,
		email,
		username,
		name,
		admin,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 3600).Unix(),
		},
	}

	rawToken := jwt.NewWithClaims(jwt.SigningMethodES512, claims)

	token, err := rawToken.SignedString([]byte(os.Getenv("tokenSign")))
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

func validateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return false
	}
	return true
}
