package wallet

type CreateWalletRequest struct {
	Nickname string `json:"nickname" binding:"required"`
}

type CreateWalletResponse struct {
	WalletNumber string `json:"walletNumber"`
	Nickname     string `json:"nickname"`
}

type WalletItem struct {
	WalletNumber string `json:"walletNumber"`
	Nickname     string `json:"nickname"`
}

type GetWalletsResponse struct {
	Wallets []WalletItem `json:"wallets"`
}
