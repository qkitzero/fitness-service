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
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	authv1 "github.com/qkitzero/auth-service/gen/go/auth/v1"
	customerv1 "github.com/qkitzero/fitness-service/gen/go/customer/v1"
	judgmentv1 "github.com/qkitzero/fitness-service/gen/go/judgment/v1"
	measurementv1 "github.com/qkitzero/fitness-service/gen/go/measurement/v1"
	measurementitemv1 "github.com/qkitzero/fitness-service/gen/go/measurementitem/v1"
	organizationv1 "github.com/qkitzero/fitness-service/gen/go/organization/v1"
	tenantv1 "github.com/qkitzero/fitness-service/gen/go/tenant/v1"
	trainingv1 "github.com/qkitzero/fitness-service/gen/go/training/v1"
	appcustomer "github.com/qkitzero/fitness-service/internal/application/customer"
	appjudgment "github.com/qkitzero/fitness-service/internal/application/judgment"
	appmeasurement "github.com/qkitzero/fitness-service/internal/application/measurement"
	appmeasurementitem "github.com/qkitzero/fitness-service/internal/application/measurementitem"
	apporganization "github.com/qkitzero/fitness-service/internal/application/organization"
	apptenant "github.com/qkitzero/fitness-service/internal/application/tenant"
	apptraining "github.com/qkitzero/fitness-service/internal/application/training"
	apiauth "github.com/qkitzero/fitness-service/internal/infrastructure/api/auth"
	apiuser "github.com/qkitzero/fitness-service/internal/infrastructure/api/user"
	infracustomer "github.com/qkitzero/fitness-service/internal/infrastructure/customer"
	"github.com/qkitzero/fitness-service/internal/infrastructure/db"
	infrajudgment "github.com/qkitzero/fitness-service/internal/infrastructure/judgment"
	inframeasurement "github.com/qkitzero/fitness-service/internal/infrastructure/measurement"
	inframeasurementitem "github.com/qkitzero/fitness-service/internal/infrastructure/measurementitem"
	infraorganization "github.com/qkitzero/fitness-service/internal/infrastructure/organization"
	infrastandard "github.com/qkitzero/fitness-service/internal/infrastructure/standard"
	infratenant "github.com/qkitzero/fitness-service/internal/infrastructure/tenant"
	infratraining "github.com/qkitzero/fitness-service/internal/infrastructure/training"
	grpccustomer "github.com/qkitzero/fitness-service/internal/interface/grpc/customer"
	grpcjudgment "github.com/qkitzero/fitness-service/internal/interface/grpc/judgment"
	grpcmeasurement "github.com/qkitzero/fitness-service/internal/interface/grpc/measurement"
	grpcmeasurementitem "github.com/qkitzero/fitness-service/internal/interface/grpc/measurementitem"
	grpcorganization "github.com/qkitzero/fitness-service/internal/interface/grpc/organization"
	grpctenant "github.com/qkitzero/fitness-service/internal/interface/grpc/tenant"
	grpctraining "github.com/qkitzero/fitness-service/internal/interface/grpc/training"
	groupv1 "github.com/qkitzero/user-service/gen/go/group/v1"
)

const shutdownTimeout = 15 * time.Second

type config struct {
	Env             string
	Port            string
	DBHost          string
	DBUser          string
	DBPassword      string
	DBName          string
	DBPort          string
	DBSSLMode       string
	AuthServiceHost string
	AuthServicePort string
	UserServiceHost string
	UserServicePort string
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
		{"AUTH_SERVICE_HOST", &cfg.AuthServiceHost},
		{"AUTH_SERVICE_PORT", &cfg.AuthServicePort},
		{"USER_SERVICE_HOST", &cfg.UserServiceHost},
		{"USER_SERVICE_PORT", &cfg.UserServicePort},
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

	var dialOpt grpc.DialOption
	switch cfg.Env {
	case "production":
		dialOpt = grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, ""))
	default:
		dialOpt = grpc.WithTransportCredentials(insecure.NewCredentials())
	}

	authConn, err := grpc.NewClient(cfg.AuthServiceHost+":"+cfg.AuthServicePort, dialOpt)
	if err != nil {
		return fmt.Errorf("auth client: %w", err)
	}
	defer func() { _ = authConn.Close() }()

	userConn, err := grpc.NewClient(cfg.UserServiceHost+":"+cfg.UserServicePort, dialOpt)
	if err != nil {
		return fmt.Errorf("user client: %w", err)
	}
	defer func() { _ = userConn.Close() }()

	server := grpc.NewServer()

	authServiceClient := authv1.NewAuthServiceClient(authConn)
	groupServiceClient := groupv1.NewGroupServiceClient(userConn)
	ageGroupStandardRepository := infrastandard.NewAgeGroupStandardRepository(gormDB)
	customerRepository := infracustomer.NewCustomerRepository(gormDB)
	judgmentRepository := infrajudgment.NewJudgmentRepository(gormDB)
	measurementRepository := inframeasurement.NewMeasurementRepository(gormDB)
	measurementItemRepository := inframeasurementitem.NewMeasurementItemRepository(gormDB)
	organizationRepository := infraorganization.NewOrganizationRepository(gormDB)
	rankStandardRepository := infrastandard.NewRankStandardRepository(gormDB)
	tenantProfileRepository := infratenant.NewProfileRepository(gormDB)
	trainingMenuRepository := infratraining.NewTrainingMenuRepository(gormDB)

	authService := apiauth.NewAuthService(authServiceClient)
	userService := apiuser.NewUserService(groupServiceClient)
	customerUsecase := appcustomer.NewCustomerUsecase(authService, userService, customerRepository, organizationRepository)
	judgmentUsecase := appjudgment.NewJudgmentUsecase(authService, userService, judgmentRepository, measurementRepository, customerRepository, measurementItemRepository, ageGroupStandardRepository, rankStandardRepository)
	measurementUsecase := appmeasurement.NewMeasurementUsecase(authService, userService, measurementRepository, customerRepository, measurementItemRepository)
	measurementItemUsecase := appmeasurementitem.NewMeasurementItemUsecase(authService, measurementItemRepository)
	organizationUsecase := apporganization.NewOrganizationUsecase(authService, userService, organizationRepository)
	tenantProfileUsecase := apptenant.NewProfileUsecase(authService, userService, tenantProfileRepository)
	trainingMenuUsecase := apptraining.NewTrainingMenuUsecase(authService, trainingMenuRepository)

	healthServer := health.NewServer()
	customerHandler := grpccustomer.NewCustomerHandler(customerUsecase)
	judgmentHandler := grpcjudgment.NewJudgmentHandler(judgmentUsecase)
	measurementHandler := grpcmeasurement.NewMeasurementHandler(measurementUsecase)
	measurementItemHandler := grpcmeasurementitem.NewMeasurementItemHandler(measurementItemUsecase)
	organizationHandler := grpcorganization.NewOrganizationHandler(organizationUsecase)
	tenantProfileHandler := grpctenant.NewProfileHandler(tenantProfileUsecase)
	trainingMenuHandler := grpctraining.NewTrainingMenuHandler(trainingMenuUsecase)

	grpc_health_v1.RegisterHealthServer(server, healthServer)
	customerv1.RegisterCustomerServiceServer(server, customerHandler)
	judgmentv1.RegisterJudgmentServiceServer(server, judgmentHandler)
	measurementv1.RegisterMeasurementServiceServer(server, measurementHandler)
	measurementitemv1.RegisterMeasurementItemServiceServer(server, measurementItemHandler)
	organizationv1.RegisterOrganizationServiceServer(server, organizationHandler)
	tenantv1.RegisterProfileServiceServer(server, tenantProfileHandler)
	trainingv1.RegisterTrainingMenuServiceServer(server, trainingMenuHandler)

	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
	healthServer.SetServingStatus("customer", grpc_health_v1.HealthCheckResponse_SERVING)
	healthServer.SetServingStatus("judgment", grpc_health_v1.HealthCheckResponse_SERVING)
	healthServer.SetServingStatus("measurement", grpc_health_v1.HealthCheckResponse_SERVING)
	healthServer.SetServingStatus("measurementitem", grpc_health_v1.HealthCheckResponse_SERVING)
	healthServer.SetServingStatus("organization", grpc_health_v1.HealthCheckResponse_SERVING)
	healthServer.SetServingStatus("tenant", grpc_health_v1.HealthCheckResponse_SERVING)
	healthServer.SetServingStatus("training", grpc_health_v1.HealthCheckResponse_SERVING)

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
