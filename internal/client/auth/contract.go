package auth

import "github.com/itksb/go-secret-keeper/pkg/contract"

// IClientAuthService - client auth service
type IClientAuthService interface {
	contract.IAuthService
}
