package core

import (
	"context"
	"ips-lacpass-backend/internal/wallet/client"
)

type WalletService struct {
	Repository *client.WalletClient
}

func NewService(r *client.WalletClient) WalletService {
	return WalletService{
		Repository: r,
	}
}

func (ws *WalletService) GenerateWalletLink(ctx context.Context, claims map[string]interface{}, credentialType client.CredentialType) (*client.GenerateWalletLinkResponse, error) {
	walletResponse, err := ws.Repository.GenerateWalletLink(ctx, claims, credentialType)
	if err != nil {
		return nil, err
	}
	return walletResponse, nil
}
