package app

import (
	"net"

	"github.com/ElishaFlacon/fast-sobes-auth/internal/config"
	"github.com/ElishaFlacon/fast-sobes-auth/internal/domain"
	"gorm.io/gorm"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type App interface {
	Run() error
	Stop() error
	GracefulStop() error
}

type app struct {
	cfg      *config.Config
	db       *gorm.DB
	log      domain.Logger
	provider *Provider
	server   *grpc.Server
}

func NewApp(cfg *config.Config, db *gorm.DB, log domain.Logger) App {
	a := &app{
		cfg: cfg,
		db:  db,
		log: log,
	}

	a.server = grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_prometheus.UnaryServerInterceptor,
		)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_prometheus.StreamServerInterceptor,
		)),
	)
	reflection.Register(a.server)
	grpc_prometheus.Register(a.server)
	grpc_prometheus.EnableHandlingTimeHistogram()

	a.provider = NewProvider(cfg, db, log)

	// === REGISTER HANDLERS ===
	a.provider.HelloHandler().RegisterImplementation(a.server)
	a.provider.AuthHandler().RegisterImplementation(a.server)
	a.provider.SettingsHandler().RegisterImplementation(a.server)

	return a
}

func (a *app) Run() error {
	lis, err := net.Listen("tcp", ":"+a.cfg.Grpc.Port)
	if err != nil {
		a.log.Errorf("failed to listen: %v", err)
		return err
	}

	if err := a.server.Serve(lis); err != nil {
		a.log.Errorf("failed to serve: %v", err)
		return err
	}

	a.log.Infof("gRPC server started on :%s", a.cfg.Grpc.Port)

	return nil
}

func (a *app) GracefulStop() error {
	a.server.GracefulStop()
	a.log.Stop()
	return nil
}

func (a *app) Stop() error {
	a.server.Stop()
	a.log.Stop()
	return nil
}
