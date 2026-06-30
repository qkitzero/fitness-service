package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	customerv1 "github.com/qkitzero/fitness-service/gen/go/customer/v1"
	appcustomer "github.com/qkitzero/fitness-service/internal/application/customer"
	infracustomer "github.com/qkitzero/fitness-service/internal/infrastructure/customer"
	"github.com/qkitzero/fitness-service/internal/infrastructure/db"
	grpccustomer "github.com/qkitzero/fitness-service/internal/interface/grpc/customer"
)

const shutdownTimeout = 15 * time.Second

type config struct {
	Env        string
	Port       string
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
	DBSSLMode  string
}

func loadConfig() (config, error) {
	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}
	cfg := config{Env: env}
	required := []struct {
		key string
		dst *string
	}{
		{"PORT", &cfg.Port},
		{"DB_HOST", &cfg.DBHost},
		{"DB_USER", &cfg.DBUser},
		{"DB_PASSWORD", &cfg.DBPassword},
		{"DB_NAME", &cfg.DBName},
		{"DB_PORT", &cfg.DBPort},
		{"DB_SSL_MODE", &cfg.DBSSLMode},
	}
	var missing []string
	for _, r := range required {
		v := os.Getenv(r.key)
		if v == "" {
			missing = append(missing, r.key)
			continue
		}
		*r.dst = v
	}
	if len(missing) > 0 {
		return cfg, fmt.Errorf("missing required env vars: %s", strings.Join(missing, ", "))
	}
	return cfg, nil
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("fitness-service: %v", err)
	}
}

func run() error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	gormDB, err := db.Init(cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode)
	if err != nil {
		return fmt.Errorf("db init: %w", err)
	}
	sqlDB, err := gormDB.DB()
	if err != nil {
		return fmt.Errorf("db handle: %w", err)
	}
	defer func() { _ = sqlDB.Close() }()

	listener, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}

	server := grpc.NewServer()

	customerRepository := infracustomer.NewCustomerRepository(gormDB)

	customerUsecase := appcustomer.NewCustomerUsecase(customerRepository)

	healthServer := health.NewServer()
	customerHandler := grpccustomer.NewCustomerHandler(customerUsecase)

	grpc_health_v1.RegisterHealthServer(server, healthServer)
	customerv1.RegisterCustomerServiceServer(server, customerHandler)

	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
	healthServer.SetServingStatus("customer", grpc_health_v1.HealthCheckResponse_SERVING)

	if cfg.Env == "development" {
		reflection.Register(server)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	serveErr := make(chan error, 1)
	go func() {
		log.Printf("gRPC server listening on %s", listener.Addr().String())
		serveErr <- server.Serve(listener)
	}()

	select {
	case err := <-serveErr:
		if err != nil {
			return fmt.Errorf("grpc serve: %w", err)
		}
		return nil
	case <-ctx.Done():
		log.Println("shutdown signal received, starting graceful stop")
		healthServer.Shutdown()

		stopped := make(chan struct{})
		go func() {
			server.GracefulStop()
			close(stopped)
		}()

		select {
		case <-stopped:
			log.Println("gRPC server stopped gracefully")
		case <-time.After(shutdownTimeout):
			log.Printf("graceful stop timed out after %s, forcing stop", shutdownTimeout)
			server.Stop()
			<-stopped
		}
		return nil
	}
}
