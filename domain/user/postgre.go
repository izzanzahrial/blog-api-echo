package user

import (
	"context"
	"database/sql"

	"github.com/izzanzahrial/blog-api-echo/entity"
)

type postgreRepository struct {
}

func NewPostgreRepository() UserRepository {
	return &postgreRepository{}
}

func (p *postgreRepository) Create(ctx context.Context, tx *sql.Tx, user entity.User) (entity.User, error) {
	SQL := "INSERT INTO user(name, pass) VALUES(?, ?)"
	result, err := tx.ExecContext(ctx, SQL, user.Name, user.Password)
	if err != nil {
		return user, ErrFailedToCreateUser
	}

	id, err := result.LastInsertId()
	if err != nil {
		return user, ErrUserNotFound
	}

	user.ID = uint64(id)

	return user, nil
}

func (p *postgreRepository) UpdateUser(ctx context.Context, tx *sql.Tx, user entity.User) (entity.User, error) {
	SQL := "UPDATE user SET name = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, SQL, user.Name, user.ID)
	if err != nil {
		return user, ErrFailedUpdateUser
	}

	return user, nil
}

func (p *postgreRepository) UpdatePassword(ctx context.Context, tx *sql.Tx, user entity.User) (entity.User, error) {
	SQL := "UPDATE user SET password = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, SQL, user.Password, user.ID)
	if err != nil {
		return user, ErrFailedUpdateUser
	}

	return user, nil
}

func (p *postgreRepository) Delete(ctx context.Context, tx *sql.Tx, user entity.User) error {
	SQL := "DELETE FROM user WHERE id = ?"
	_, err := tx.ExecContext(ctx, SQL, user.ID)
	if err != nil {
		return ErrFailedToDeleteUser
	}

	return nil
}

func (p *postgreRepository) Login(ctx context.Context, tx *sql.Tx, ID uint64, pass string) (entity.User, error) {
	SQL := "SELECT name FROM user WHERE id = ? AND password = ?"
	rows, err := tx.QueryContext(ctx, SQL, ID, pass)
	if err != nil {
		return entity.User{}, ErrUserNotFound
	}
	defer rows.Close()

	user := entity.User{}
	if rows.Next() {
		if err := rows.Scan(&user.Name); err != nil {
			return user, ErrFailedToAssertUser
		}
		return user, nil
	} else {
		return user, ErrUserNotFound
	}
}

func (p *postgreRepository) FindByID(ctx context.Context, tx *sql.Tx, ID uint64) (entity.User, error) {
	SQL := "SELECT name FROM user WHERE id = ?"
	rows, err := tx.QueryContext(ctx, SQL, ID)
	if err != nil {
		return entity.User{}, ErrUserNotFound
	}
	defer rows.Close()

	user := entity.User{}
	if rows.Next() {
		if err := rows.Scan(&user.Name); err != nil {
			return user, ErrFailedToAssertUser
		}
		return user, nil
	} else {
		return user, ErrUserNotFound
	}
}
