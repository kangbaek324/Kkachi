package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	db "github.com/kangbaek324/kkachi/db/sqlc"
	"golang.org/x/crypto/bcrypt"
)

var ErrUsernameAlreadyExists = errors.New("username already exists")

type Service interface {
	Register(ctx context.Context, req CreateUserRequest) (CreateUserResponse, error)
}

type userService struct {
	q db.Querier
}

func NewService(q db.Querier) Service {
	return &userService{q: q}
}

func (s *userService) Register(ctx context.Context, req CreateUserRequest) (CreateUserResponse, error) {
	// Password Hahsing
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return CreateUserResponse{}, fmt.Errorf("hash password: %w", err)
	}

	// Create User
	user, err := s.q.CreateUser(ctx, db.CreateUserParams{
		Username: req.Username,
		Password: string(hashed),
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return CreateUserResponse{}, ErrUsernameAlreadyExists
		}
		return CreateUserResponse{}, err
	}

	return CreateUserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Time,
	}, nil
}
