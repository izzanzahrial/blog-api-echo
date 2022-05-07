package repository

import (
	"context"
	"database/sql"
)

type userPostgre struct {
}

func NewUserPostgreRepository() UserRepository {
	return &userPostgre{}
}

func (p *userPostgre) Create(ctx context.Context, tx *sql.Tx, u User) (User, error) {
	SQL := "INSERT INTO user(email, username, name, password) VALUES(?, ?, ?, ?)"
	result, err := tx.ExecContext(ctx, SQL, u.Email, u.Username, u.Name, u.Password)
	if err != nil {
		return User{}, ErrFailedToCreateUser
	}

	id, err := result.LastInsertId()
	if err != nil {
		return User{}, ErrUserNotFound
	}

	u.ID = id

	return u, nil
}

func (p *userPostgre) UpdateUser(ctx context.Context, tx *sql.Tx, u User) (User, error) {
	SQL := "UPDATE user SET username = ?, name = ? WHERE LOWER(email) = LOWER(?)"
	_, err := tx.ExecContext(ctx, SQL, u.Name, u.ID)
	if err != nil {
		return User{}, ErrFailedUpdateUser
	}

	return u, nil
}

func (p *userPostgre) UpdatePassword(ctx context.Context, tx *sql.Tx, u User) (User, error) {
	SQL := "UPDATE user SET password = ? WHERE LOWER(email) = LOWER(?)"
	_, err := tx.ExecContext(ctx, SQL, u.Password, u.Email)
	if err != nil {
		return u, ErrFailedUpdateUser
	}

	return u, nil
}

func (p *userPostgre) Delete(ctx context.Context, tx *sql.Tx, u User) error {
	SQL := "DELETE FROM user WHERE LOWER(email) = LOWER(?)"
	_, err := tx.ExecContext(ctx, SQL, u.Email)
	if err != nil {
		return ErrFailedToDeleteUser
	}

	return nil
}

func (p *userPostgre) LoginByEmail(ctx context.Context, tx *sql.Tx, email string, pass string) (User, error) {
	SQL := "SELECT id, email, username, name FROM user WHERE LOWER(email) = LOWER(?) AND password = ?"
	rows, err := tx.QueryContext(ctx, SQL, email, pass)
	if err != nil {
		return User{}, ErrUserNotFound
	}
	defer rows.Close()

	var user User
	if rows.Next() {
		if err := rows.Scan(&user.ID, &user.Email, &user.Username, &user.Name); err != nil {
			return user, ErrFailedToAssertUser
		}
		return user, nil
	} else {
		return user, ErrUserNotFound
	}
}

func (p *userPostgre) LoginByUsername(ctx context.Context, tx *sql.Tx, username string, pass string) (User, error) {
	SQL := "SELECT id, email, username, name FROM user WHERE LOWER(username) = LOWER(?) AND password = ?"
	rows, err := tx.QueryContext(ctx, SQL, username, pass)
	if err != nil {
		return User{}, ErrUserNotFound
	}
	defer rows.Close()

	user := User{}
	if rows.Next() {
		if err := rows.Scan(&user.ID, &user.Email, &user.Username, &user.Name); err != nil {
			return user, ErrFailedToAssertUser
		}
		return user, nil
	} else {
		return user, ErrUserNotFound
	}
}

func (p *userPostgre) FindByEmail(ctx context.Context, tx *sql.Tx, email string) (User, error) {
	SQL := "SELECT id, email, username, name FROM user WHERE LOWER(email) = LOWER(?)"
	rows, err := tx.QueryContext(ctx, SQL, email)
	if err != nil {
		return User{}, ErrUserNotFound
	}
	defer rows.Close()

	user := User{}
	if rows.Next() {
		if err := rows.Scan(&user.ID, &user.Email, &user.Username, &user.Name); err != nil {
			return user, ErrFailedToAssertUser
		}
		return user, nil
	} else {
		return user, ErrUserNotFound
	}
}

func (p *userPostgre) FindByUsername(ctx context.Context, tx *sql.Tx, username string) (User, error) {
	SQL := "SELECT id, email, username, name FROM user WHERE LOWER(username) = LOWER(?)"
	rows, err := tx.QueryContext(ctx, SQL, username)
	if err != nil {
		return User{}, ErrUserNotFound
	}
	defer rows.Close()

	user := User{}
	if rows.Next() {
		if err := rows.Scan(&user.ID, &user.Email, &user.Username, &user.Name); err != nil {
			return user, ErrFailedToAssertUser
		}
		return user, nil
	} else {
		return user, ErrUserNotFound
	}
}
