package client

import (
	"context"
	"fmt"
	"github.com/itksb/go-secret-keeper/internal/client/auth"
	"github.com/itksb/go-secret-keeper/internal/client/cipher"
	"github.com/itksb/go-secret-keeper/internal/client/command"
	"github.com/itksb/go-secret-keeper/internal/client/gui/term"
	"github.com/itksb/go-secret-keeper/internal/client/keeper"
	"github.com/itksb/go-secret-keeper/internal/client/session"
	"github.com/itksb/go-secret-keeper/pkg/contract"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"time"
)

var _ contract.IApplication = &ClientApp{}

// ClientApp - client application
type ClientApp struct {
	cfg         Config
	l           contract.IApplicationLogger
	deferredOps []func() error
}

// NewClientApp - constructor
// composition root
func NewClientApp(cfg Config) *ClientApp {
	var logger *zap.Logger
	var err error
	if cfg.Debug {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}
	if err != nil {
		fmt.Printf("can't initialize zap logger: %v", err)
	}

	sugared := logger.Sugar()

	return &ClientApp{
		cfg: cfg,
		l:   sugared,
		deferredOps: append([]func() error{}, func() error {
			sugared.Info("executing deferred operation logger.Sync()")
			return logger.Sync()
		}),
	}
}

// Run - run client
func (c *ClientApp) Run() error {
	guiSession := session.NewAppSession()

	authService := auth.NewClientAuthService()

	cryptoService := cipher.NewCryptoService(func() ([]byte, error) {
		return []byte(c.cfg.CryptoKey), nil
	})

	keeperApi := keeper.NewAPIKeeper(cryptoService, c.l)
	keeperService := keeper.NewClientKeeper(keeperApi, c.l)

	gui := term.NewTerminalService(
		c.l,
		guiSession,
		command.LoginCmdAbstractFabric(authService, c.l),
		command.RegisterCmdAbstractFabric(authService, c.l),
		command.ListSecretsCmdAbstractFabric(c.l, keeperService),
	)

	return gui.Start()
}

// Stop - stop client
func (c *ClientApp) Stop() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancelFunc()
	g, _ := errgroup.WithContext(ctx)

	for _, op := range c.deferredOps {
		g.Go(op)
	}

	return g.Wait()
}
