package app

import (
	"context"
	"net/http"
	"time"

	"github.com/Thanhbinh1905/go-training-system/services/team-service/config"
	"github.com/Thanhbinh1905/go-training-system/services/team-service/internal/client"
	grpcHandler "github.com/Thanhbinh1905/go-training-system/services/team-service/internal/handler/grpc"
	httpHandler "github.com/Thanhbinh1905/go-training-system/services/team-service/internal/handler/http"
	"github.com/Thanhbinh1905/go-training-system/services/team-service/internal/middleware"
	"github.com/Thanhbinh1905/go-training-system/services/team-service/internal/repository"
	"github.com/Thanhbinh1905/go-training-system/services/team-service/internal/service"
	teampb "github.com/Thanhbinh1905/go-training-system/services/team-service/pb/team"
	userpb "github.com/Thanhbinh1905/go-training-system/services/team-service/pb/user"
	"github.com/Thanhbinh1905/go-training-system/shared/db/postgres"
	"github.com/Thanhbinh1905/go-training-system/shared/db/redis"
	"github.com/Thanhbinh1905/go-training-system/shared/graceful"
	"github.com/Thanhbinh1905/go-training-system/shared/kafka"
	"github.com/Thanhbinh1905/go-training-system/shared/logger"

	"net"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

const (
	GRATEFUL_TIMEOUT = 30 * time.Second
)

func Run(cfg *config.Config) {
	log := logger.InitLogger("services/team-service/logs/team-service.log", "team-service")
	defer log.Sync()

	conn, err := postgres.Connect(cfg.DatabaseURL, log)
	if err != nil {
		log.Fatal("failed to connect to database", zap.Error(err))
	}
	defer postgres.Close(conn, log)

	redisClient, err := redis.Init(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB, log)
	if err != nil {
		log.Fatal("failed to init redis:", zap.Error(err))
	}
	defer redisClient.Close(log)

	teamDBRepo := repository.NewTeamDbRepository(conn)
	teamCacheRepo := repository.NewTeamCache(redisClient.Client)

	// Initialize Kafka producer
	kafkaProducer, err := kafka.NewProducer([]string{"kafka:29092"})
	if err != nil {
		log.Fatal("failed to create kafka producer", zap.Error(err))
	}
	defer kafkaProducer.Close()

	teamRepo := repository.NewTeamRepository(teamDBRepo, teamCacheRepo)
	userClient := client.NewUserGRPCClient(cfg.UserGRPCURL)
	teamService := service.NewTeamService(teamRepo, *userClient, kafkaProducer)

	// Start gRPC server
	grpcServer, grpcListener, err := runGRPCServer(teamService, cfg.GRPCPort, log)
	if err != nil {
		log.Fatal("Failed to start gRPC server", zap.Error(err))
	}

	// Health checks
	dbCheck := func() error {
		sqlDB, err := conn.DB()
		if err != nil {
			return err
		}
		return sqlDB.Ping()
	}
	grpcCheck := func() error {
		d, err := net.DialTimeout("tcp", ":"+cfg.GRPCPort, 2*time.Second)
		if err != nil {
			return err
		}
		_ = d.Close()
		return nil
	}

	// Start HTTP server
	httpServer, err := runHTTPServer(teamService, *userClient, log, dbCheck, grpcCheck)
	if err != nil {
		log.Fatal("Failed to start HTTP server", zap.Error(err))
	}

	// Create graceful shutdown server
	gracefulServer := graceful.NewGracefulServer(GRATEFUL_TIMEOUT)

	// Add services for graceful shutdown
	gracefulServer.AddService(&gracefulService{
		name:         "Team Service",
		grpcServer:   grpcServer,
		grpcListener: grpcListener,
		httpServer:   httpServer,
		log:          log,
	})

	// Start graceful shutdown listener
	gracefulServer.Start()
}

func runGRPCServer(teamService service.TeamService, port string, log *zap.Logger) (*grpc.Server, net.Listener, error) {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return nil, nil, err
	}
	server := grpc.NewServer()
	teampb.RegisterTeamServiceServer(server, grpcHandler.NewTeamgRPCHandler(teamService))

	log.Info("gRPC server listening", zap.String("port", port))

	// Start server in goroutine
	go func() {
		if err := server.Serve(lis); err != nil {
			log.Error("gRPC server error", zap.Error(err))
		}
	}()

	return server, lis, nil
}

func runHTTPServer(teamService service.TeamService, userClient client.UserGRPCClient, log *zap.Logger, dbCheck func() error, grpcCheck func() error) (*http.Server, error) {
	r := gin.Default()
	r.Use(ginzap.Ginzap(log, time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(log, true))

	// Metrics
	p := ginprometheus.NewPrometheus("team_service")
	p.Use(r)

	r.GET("/health", func(c *gin.Context) {
		health := gin.H{"status": "up"}
		if dbCheck != nil {
			if err := dbCheck(); err != nil {
				health["db"] = "down"
				health["db_error"] = err.Error()
			} else {
				health["db"] = "up"
			}
		}
		if grpcCheck != nil {
			if err := grpcCheck(); err != nil {
				health["grpc"] = "down"
				health["grpc_error"] = err.Error()
			} else {
				health["grpc"] = "up"
			}
		}
		c.JSON(200, health)
	})

	teamHandler := httpHandler.NewTeamHandler(teamService)

	api := r.Group("/api/v1")
	teamGroup := api.Group("/teams")
	teamGroup.Use(middleware.AuthMiddleware(userClient.Client, userpb.Role_MANAGER))
	{
		teamGroup.POST("/", teamHandler.CreateTeam)
		teamGroup.POST("/:teamID/managers", teamHandler.AddManager)
		teamGroup.POST("/:teamID/members", teamHandler.AddMember)
		teamGroup.DELETE("/:teamID/managers/:managerID", teamHandler.RemoveManager)
		teamGroup.DELETE("/:teamID/members/:memberID", teamHandler.RemoveMember)
	}

	// Create HTTP server
	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	log.Info("HTTP server started", zap.String("port", "8080"))

	// Start server in goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("failed to run HTTP server", zap.Error(err))
		}
	}()

	return server, nil
}

// gracefulService implements Shutdownable interface
type gracefulService struct {
	name         string
	grpcServer   *grpc.Server
	grpcListener net.Listener
	httpServer   *http.Server
	log          *zap.Logger
}

func (s *gracefulService) Shutdown(ctx context.Context) error {
	s.log.Info("Shutting down " + s.name)

	// Shutdown HTTP server
	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.log.Error("HTTP server shutdown error", zap.Error(err))
	}

	// Shutdown gRPC server
	s.grpcServer.GracefulStop()

	// Close gRPC listener
	if err := s.grpcListener.Close(); err != nil {
		s.log.Error("gRPC listener close error", zap.Error(err))
	}

	s.log.Info(s.name + " shutdown completed")
	return nil
}
