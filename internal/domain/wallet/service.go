package wallet

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"

	"github.com/jackc/pgx/v5/pgconn"
	db "github.com/kangbaek324/kkachi/db/sqlc"
)

type Service interface {
	CreateWallet(ctx context.Context, req CreateWalletRequest, userId int64) (CreateWalletResponse, error)
}

type walletService struct {
	q db.Querier
}

func NewService(q db.Querier) Service {
	return &walletService{q: q}
}

func (s *walletService) CreateWallet(ctx context.Context, req CreateWalletRequest, userId int64) (CreateWalletResponse, error) {
	for range 5 {
		walletNumber := fmt.Sprintf("KK-%06d", rand.IntN(900000)+100000)
		wallet, err := s.q.CreateWallet(ctx, db.CreateWalletParams{
			UserID:       userId,
			WalletNumber: walletNumber,
			Nickname:     req.Nickname,
		})
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				continue
			}
			return CreateWalletResponse{}, fmt.Errorf("createWallet: %w", err)
		}

		return CreateWalletResponse{
			WalletNumber: wallet.WalletNumber,
			Nickname:     wallet.Nickname,
		}, nil
	}

	return CreateWalletResponse{}, fmt.Errorf("createWallet: failed to generate unique wallet number")
}
