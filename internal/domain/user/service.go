package user

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/kangbaek324/kkachi/db/sqlc"
	"golang.org/x/crypto/bcrypt"
)

var ErrUsernameAlreadyExists = errors.New("username already exists")
var ErrInvalidCredentials = errors.New("invalid credentials")

var accessTokenExpireDate = time.Hour
var refershTokenExpireDate = 7 * 24 * time.Hour

type Service interface {
	Register(ctx context.Context, req CreateUserRequest) (CreateUserResponse, error)
	Login(ctx context.Context, req LoginRequest) (LoginResponse, error)
}

type userService struct {
	q         db.Querier
	jwtSecret string
}

type tokenPair struct {
	access  string
	refresh string
}

func NewService(q db.Querier, jwtSecret string) Service {
	return &userService{q: q, jwtSecret: jwtSecret}
}

func (s *userService) Register(ctx context.Context, req CreateUserRequest) (CreateUserResponse, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return CreateUserResponse{}, fmt.Errorf("hash password: %w", err)
	}

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

func (s *userService) Login(ctx context.Context, req LoginRequest) (LoginResponse, error) {
	user, err := s.q.GetUser(ctx, req.Username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return LoginResponse{}, ErrInvalidCredentials
		}
		return LoginResponse{}, fmt.Errorf("get user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return LoginResponse{}, ErrInvalidCredentials
	}

	tok, err := s.generateTokenPair(user.ID)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("generate token: %w", err)
	}

	_, err = s.q.CreateRefreshToken(ctx, db.CreateRefreshTokenParams{
		UserID: user.ID,
		Token:  tok.refresh,
		ExpiresAt: pgtype.Timestamptz{
			Time:  time.Now().Add(7 * 24 * time.Hour),
			Valid: true,
		},
	})
	if err != nil {
		return LoginResponse{}, fmt.Errorf("CreateRefreshToken: %w", err)
	}

	return LoginResponse{
		AccessToken:  tok.access,
		RefreshToken: tok.refresh,
	}, nil
}

func (s *userService) generateTokenPair(userID int64) (tokenPair, error) {
	access, err := signJWT(userID, s.jwtSecret, accessTokenExpireDate)
	if err != nil {
		return tokenPair{}, err
	}
	refresh, err := signJWT(userID, s.jwtSecret, refershTokenExpireDate)
	if err != nil {
		return tokenPair{}, err
	}
	return tokenPair{access: access, refresh: refresh}, nil
}

func signJWT(userID int64, secret string, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"sub": strconv.FormatInt(userID, 10),
		"exp": time.Now().Add(duration).Unix(),
		"iat": time.Now().Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(secret))
}
