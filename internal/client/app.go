package client

import (
	"context"
	"fmt"
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
	time.Sleep(1 * time.Second)
	return nil
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
