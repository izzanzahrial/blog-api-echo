package repository

import (
	"context"
	"database/sql"

	user "github.com/izzanzahrial/blog-api-echo/pkg/user"
)

type userPostgre struct {
}

func NewUserPostgreRepository() UserRepository {
	return &userPostgre{}
}

func (p *userPostgre) Create(ctx context.Context, tx *sql.Tx, user user.User) (user.User, error) {
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

func (p *userPostgre) UpdateUser(ctx context.Context, tx *sql.Tx, user user.User) (user.User, error) {
	SQL := "UPDATE user SET name = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, SQL, user.Name, user.ID)
	if err != nil {
		return user, ErrFailedUpdateUser
	}

	return user, nil
}

func (p *userPostgre) UpdatePassword(ctx context.Context, tx *sql.Tx, user user.User) (user.User, error) {
	SQL := "UPDATE user SET password = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, SQL, user.Password, user.ID)
	if err != nil {
		return user, ErrFailedUpdateUser
	}

	return user, nil
}

func (p *userPostgre) Delete(ctx context.Context, tx *sql.Tx, user user.User) error {
	SQL := "DELETE FROM user WHERE id = ?"
	_, err := tx.ExecContext(ctx, SQL, user.ID)
	if err != nil {
		return ErrFailedToDeleteUser
	}

	return nil
}

func (p *userPostgre) Login(ctx context.Context, tx *sql.Tx, ID uint64, pass string) (user.User, error) {
	SQL := "SELECT name FROM user WHERE id = ? AND password = ?"
	rows, err := tx.QueryContext(ctx, SQL, ID, pass)
	if err != nil {
		return user.User{}, ErrUserNotFound
	}
	defer rows.Close()

	user := user.User{}
	if rows.Next() {
		if err := rows.Scan(&user.Name); err != nil {
			return user, ErrFailedToAssertUser
		}
		return user, nil
	} else {
		return user, ErrUserNotFound
	}
}

func (p *userPostgre) FindByID(ctx context.Context, tx *sql.Tx, ID uint64) (user.User, error) {
	SQL := "SELECT name FROM user WHERE id = ?"
	rows, err := tx.QueryContext(ctx, SQL, ID)
	if err != nil {
		return user.User{}, ErrUserNotFound
	}
	defer rows.Close()

	user := user.User{}
	if rows.Next() {
		if err := rows.Scan(&user.Name); err != nil {
			return user, ErrFailedToAssertUser
		}
		return user, nil
	} else {
		return user, ErrUserNotFound
	}
}
