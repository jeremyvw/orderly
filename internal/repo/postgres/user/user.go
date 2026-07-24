package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"

	"orderly/internal/entity"
)

const uniqueViolation = "23505"

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, email, passwordHash string) (entity.User, error) {
	const query = `
		INSERT INTO users (email, password_hash)
		VALUES ($1, $2)
		RETURNING id, email, created_at`

	var user entity.User
	err := r.db.QueryRowContext(ctx, query, email, passwordHash).
		Scan(&user.ID, &user.Email, &user.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == uniqueViolation {
			return entity.User{}, entity.ErrEmailTaken
		}
		return entity.User{}, fmt.Errorf("insert user: %w", err)
	}

	return user, nil
}

func (r *UserRepo) FindByEmail(ctx context.Context, email string) (entity.User, error) {
	const query = `
		SELECT id, email, password_hash, created_at
		FROM users
		WHERE LOWER(email) = LOWER($1)`

	var user entity.User
	err := r.db.QueryRowContext(ctx, query, email).
		Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return entity.User{}, entity.ErrUserNotFound
	}
	if err != nil {
		return entity.User{}, fmt.Errorf("find user by email: %w", err)
	}

	return user, nil
}
