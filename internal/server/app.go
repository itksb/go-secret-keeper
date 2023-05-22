package server

import (
	"context"
	"fmt"
	"github.com/itksb/go-secret-keeper/pkg/contract"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"time"
)

var _ contract.IApplication = &ServerApp{}

// ServerApp - server application
type ServerApp struct {
	cfg         Config
	l           contract.IApplicationLogger
	deferredOps []func() error
}

// NewServerApp - constructor
func NewServerApp(cfg Config) *ServerApp {
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

	return &ServerApp{
		cfg: cfg,
		l:   sugared,
		deferredOps: append([]func() error{}, func() error {
			sugared.Info("executing deferred operation logger.Sync()")
			return logger.Sync()
		}),
	}
}

// Run - run server
func (s *ServerApp) Run() error {
	time.Sleep(1 * time.Second)
	return nil
}

// Stop - stop server
func (s *ServerApp) Stop() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancelFunc()
	g, _ := errgroup.WithContext(ctx)

	for _, op := range s.deferredOps {
		g.Go(op)
	}

	return g.Wait()
}