package keeper

import "github.com/itksb/go-secret-keeper/pkg/contract"

// IClientKeeper - client keeper interface
type IClientKeeper interface {
	contract.IKeeper
}

// GetPrivateKeyFunc - get private key function
type GetPrivateKeyFunc func() ([]byte, error)
