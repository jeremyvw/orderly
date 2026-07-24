package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"orderly/internal/entity"
	"orderly/internal/pkg/password"
)

const (
	minPasswordLength = 8
	maxPasswordLength = 72 // bcrypt ignores bytes beyond this
)

// dummyHash is compared against when no user is found, so login consumes
// comparable CPU whether or not the email is registered.
var dummyHash = mustHash("dummy-password-for-constant-time-compare")

type UserRepository interface {
	Create(ctx context.Context, email, passwordHash string) (entity.User, error)
	FindByEmail(ctx context.Context, email string) (entity.User, error)
}

type TokenIssuer func(userID int64) (string, error)

type Usecase struct {
	userRepo   UserRepository
	issueToken TokenIssuer
}

func New(userRepo UserRepository, issueToken TokenIssuer) *Usecase {
	return &Usecase{
		userRepo:   userRepo,
		issueToken: issueToken,
	}
}

func (u *Usecase) Signup(ctx context.Context, email, plainPassword string) (string, error) {
	email = normalizeEmail(email)

	if err := validate(email, plainPassword); err != nil {
		return "", err
	}

	hashed, err := password.Hash(plainPassword)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}

	// Uniqueness is enforced by the index; the repo translates 23505
	// into entity.ErrEmailTaken.
	user, err := u.userRepo.Create(ctx, email, hashed)
	if err != nil {
		return "", err
	}

	return u.issue(user.ID)
}

func (u *Usecase) Login(ctx context.Context, email, plainPassword string) (string, error) {
	email = normalizeEmail(email)

	user, err := u.userRepo.FindByEmail(ctx, email)
	if errors.Is(err, entity.ErrUserNotFound) {
		_ = password.Compare(plainPassword, dummyHash)
		return "", entity.ErrInvalidCredentials
	}
	if err != nil {
		return "", fmt.Errorf("find user by email: %w", err)
	}

	if err := password.Compare(plainPassword, user.PasswordHash); err != nil {
		return "", entity.ErrInvalidCredentials
	}

	return u.issue(user.ID)
}

func (u *Usecase) issue(userID int64) (string, error) {
	tok, err := u.issueToken(userID)
	if err != nil {
		return "", fmt.Errorf("issue token: %w", err)
	}
	return tok, nil
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func validate(email, plainPassword string) error {
	if email == "" || !strings.Contains(email, "@") {
		return fmt.Errorf("%w: email is not valid", entity.ErrInvalidInput)
	}
	if len(plainPassword) < minPasswordLength {
		return fmt.Errorf("%w: password must be at least %d characters",
			entity.ErrInvalidInput, minPasswordLength)
	}
	if len(plainPassword) > maxPasswordLength {
		return fmt.Errorf("%w: password must be at most %d characters",
			entity.ErrInvalidInput, maxPasswordLength)
	}
	return nil
}

func mustHash(plain string) string {
	h, err := password.Hash(plain)
	if err != nil {
		panic(fmt.Sprintf("precompute dummy hash: %v", err))
	}
	return h
}
