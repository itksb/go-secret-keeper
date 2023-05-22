package keeper

import (
	"context"
	"time"
)

type IKeeperRepository interface {
	SaveSecret(ctx context.Context, userID string, secret UserSecret) error
	GetLastUpdate(ctx context.Context, userID string) (time.Time, error)
}
