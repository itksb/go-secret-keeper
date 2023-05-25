package keeper_test

import (
	"crypto/sha256"
	"fmt"
	"github.com/itksb/go-secret-keeper/internal/client/keeper"
)

// GetPrivateKeyFunc - example of GetPrivateKeyFunc
func ExampleGetPrivateKeyFunc() {
	key := []byte("1234567890") // source from config or env

	f := keeper.GetPrivateKeyFunc(func() ([]byte, error) {
		sum := sha256.Sum256(key)
		return sum[:], nil
	})

	result, _ := f()
	fmt.Println(result)
}
