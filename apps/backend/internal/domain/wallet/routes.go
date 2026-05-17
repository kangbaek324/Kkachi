package wallet

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/kangbaek324/kkachi/apps/backend/db/sqlc"
	"github.com/kangbaek324/kkachi/apps/backend/internal/middleware"
)

func Register(rg *gin.RouterGroup, pool *pgxpool.Pool, jwtSecret string) {
	svc := NewService(db.New(pool), pool)
	h := NewHandler(svc)

	wallets := rg.Group("/wallets")
	wallets.Use(middleware.Auth(jwtSecret))

	// Wallet
	wallets.POST("/", h.createWallet)
	wallets.GET("/", h.getWallets)
	wallets.GET("/:wallet_number/balances", h.getWalletBalances)
	wallets.PATCH("/", h.editWalletNickname)

	// Transfer
	wallets.POST("/:wallet_number/transfers", h.transfer)
	wallets.GET("/:wallet_number/transfers", h.getTransferLogs)

	// Exchange
	wallets.POST("/:wallet_number/exchanges", h.exchange)
	wallets.GET("/:wallet_number/exchanges", h.getExchangeLogs)
}
