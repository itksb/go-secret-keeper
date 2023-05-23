package server

import (
	"context"
	"fmt"
	"github.com/itksb/go-secret-keeper/internal/server/db"
	"github.com/itksb/go-secret-keeper/migrate"
	"github.com/itksb/go-secret-keeper/pkg/contract"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"time"
)

var _ contract.IApplication = &ServerApp{}

// ServerApp - server application
type ServerApp struct {
	cfg         Config
	l           contract.IApplicationLogger
	GRPCServer  *grpc.Server
	db          *sqlx.DB
	appMigrator migrate.IAppMigrator
	deferredOps []func() error
}

// NewServerApp - constructor
func NewServerApp(
	cfg Config,
	appMigrator migrate.IAppMigrator,
) *ServerApp {
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
		cfg:         cfg,
		l:           sugared,
		appMigrator: appMigrator,
		deferredOps: append([]func() error{}, func() error {
			sugared.Info("executing deferred operation logger.Sync()")
			return logger.Sync()
		}),
	}
}

// Run - run server
func (s *ServerApp) Run() error {
	// run migrations
	err := s.appMigrator.Migrate(s.cfg.Dsn, migrate.Migrations)
	if err != nil {
		s.l.Errorf("migration error: %s", err.Error())
		return err
	}
	s.db, err = db.NewPostgresDbPool(s.cfg.Dsn, s.l)
	if err != nil {
		return err
	}

	s.deferredOps = append(s.deferredOps, func() error {
		s.l.Infof("closing db connection")
		return s.db.Close()
	})

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
