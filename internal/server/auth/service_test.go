package auth

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/itksb/go-secret-keeper/pkg/contract"
	mock_contract "github.com/itksb/go-secret-keeper/pkg/contract/mock"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestAuthService_SignIn(t *testing.T) {
	ctrl := gomock.NewController(t)
	tok := mock_contract.NewMockITokenProvider(ctrl)
	repo := mock_contract.NewMockIAuthRepository(ctrl)
	hasher := mock_contract.NewMockIPassHasher(ctrl)
	authService := NewAuthService(
		tok,
		repo,
		hasher,
	)
	login := "testUser"
	password := "password123"

	hasher.EXPECT().HashPassword(gomock.Any()).Return(password, nil)
	repo.EXPECT().Find(gomock.Any(), login, password).Return(&Account{
		ID:           "testID",
		Login:        login,
		PasswordHash: "hash",
	}, nil)

	account, err := authService.SignIn(
		context.Background(),
		login,
		password,
	)
	assert.NoError(t, err, "SignIn should not return an error")
	assert.NotNil(t, account, "SignIn should return an account")
	assert.Equal(
		t,
		login,
		account.GetLogin(),
		"SignIn should find an account with the specified login, passwd",
	)
}

func TestAuthService_SignUp(t *testing.T) {
	ctrl := gomock.NewController(t)
	tok := mock_contract.NewMockITokenProvider(ctrl)
	repo := mock_contract.NewMockIAuthRepository(ctrl)
	hasher := mock_contract.NewMockIPassHasher(ctrl)
	authService := NewAuthService(
		tok,
		repo,
		hasher,
	)
	login := "testUser"
	password := "password123"

	hasher.EXPECT().HashPassword(gomock.Any()).Return(password, nil)
	repo.EXPECT().Create(gomock.Any(), login, password).Return(&Account{
		ID:           "testID",
		Login:        login,
		PasswordHash: "hash",
	}, nil)

	account, err := authService.SignUp(
		context.Background(),
		login,
		password,
	)
	assert.NoError(t, err, "SignUp should not return an error")
	assert.NotNil(t, account, "SignUp should return an account")
	assert.Equal(
		t,
		login,
		account.GetLogin(),
		"SignUp should create an account with the specified login",
	)
}

func TestNewAuthService(t *testing.T) {
	type args struct {
		tokenProvider contract.ITokenProvider
		repo          contract.IAuthRepository
		passHasher    contract.IPassHasher
	}
	tests := []struct {
		name string
		args args
		want *AuthService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAuthService(tt.args.tokenProvider, tt.args.repo, tt.args.passHasher); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAuthService() = %v, want %v", got, tt.want)
			}
		})
	}
}
