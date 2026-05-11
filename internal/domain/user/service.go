package user

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/kangbaek324/kkachi/db/sqlc"
	"github.com/kangbaek324/kkachi/internal/common"
	"golang.org/x/crypto/bcrypt"
)

var ErrUsernameAlreadyExists = common.NewAppError(http.StatusConflict, "username already exists")
var ErrInvalidCredentials = common.NewAppError(http.StatusUnauthorized, "invalid credentials")
var ErrRefreshTokenExpired = common.NewAppError(http.StatusUnauthorized, "refreshToken expired")

var accessTokenExpireDate = time.Hour
var refreshTokenExpireDate = time.Hour / 60 / 60

type Service interface {
	Register(ctx context.Context, req CreateUserRequest) (CreateUserResponse, error)
	Login(ctx context.Context, req LoginRequest) (LoginResponse, error)
	RefreshAccessToken(ctx context.Context, req RefreshAccessTokenRequest) (RefreshAccessTokenResponse, error)
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
			return CreateUserResponse{}, fmt.Errorf("register: %w", ErrUsernameAlreadyExists)
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
			return LoginResponse{}, fmt.Errorf("login: %w", ErrInvalidCredentials)
		}
		return LoginResponse{}, fmt.Errorf("login: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return LoginResponse{}, fmt.Errorf("login: %w", ErrInvalidCredentials)
	}

	tok, err := s.generateTokenPair(user.ID)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("generate token: %w", err)
	}

	_, err = s.q.CreateRefreshToken(ctx, db.CreateRefreshTokenParams{
		UserID: user.ID,
		Token:  hashToken(tok.refresh),
		ExpiresAt: pgtype.Timestamptz{
			Time:  time.Now().Add(refreshTokenExpireDate),
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

func (s *userService) RefreshAccessToken(ctx context.Context, req RefreshAccessTokenRequest) (RefreshAccessTokenResponse, error) {
	refreshToken, err := s.q.GetRefreshToken(ctx, hashToken(req.RefreshToken))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return RefreshAccessTokenResponse{}, fmt.Errorf("refreshAccessToken: %w", ErrInvalidCredentials)
		}
		return RefreshAccessTokenResponse{}, fmt.Errorf("refreshAccessToken: %w", err)
	}

	if refreshToken.ExpiresAt.Time.Before(time.Now()) {
		err := s.q.DeleteRefreshToken(ctx, refreshToken.Token)
		if err != nil {
			return RefreshAccessTokenResponse{}, fmt.Errorf("deleteRefershToken: %w", err)
		}

		return RefreshAccessTokenResponse{}, fmt.Errorf("refreshAccessToken: %w", ErrRefreshTokenExpired)
	}

	accessToken, err := signJWT(refreshToken.UserID, s.jwtSecret, accessTokenExpireDate)
	if err != nil {
		return RefreshAccessTokenResponse{}, err
	}
	return RefreshAccessTokenResponse{
		AccessToken: accessToken,
	}, nil
}

func (s *userService) generateTokenPair(userID int64) (tokenPair, error) {
	access, err := signJWT(userID, s.jwtSecret, accessTokenExpireDate)
	if err != nil {
		return tokenPair{}, err
	}
	refresh, err := signJWT(userID, s.jwtSecret, refreshTokenExpireDate)
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

func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return fmt.Sprintf("%x", h)
}
