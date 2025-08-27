package app

import (
	"context"
	"net/http"
	"time"

	"github.com/Thanhbinh1905/go-training-system/services/asset-service/config"
	"github.com/Thanhbinh1905/go-training-system/services/asset-service/internal/client"
	"github.com/Thanhbinh1905/go-training-system/services/asset-service/internal/handler"
	"github.com/Thanhbinh1905/go-training-system/services/asset-service/internal/middleware"
	"github.com/Thanhbinh1905/go-training-system/services/asset-service/internal/repository"
	"github.com/Thanhbinh1905/go-training-system/services/asset-service/internal/service"
	"github.com/Thanhbinh1905/go-training-system/shared/db/postgres"
	"github.com/Thanhbinh1905/go-training-system/shared/db/redis"
	"github.com/Thanhbinh1905/go-training-system/shared/graceful"
	"github.com/Thanhbinh1905/go-training-system/shared/kafka"
	"github.com/Thanhbinh1905/go-training-system/shared/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	ginzap "github.com/gin-contrib/zap"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

const (
	GRATEFUL_TIMEOUT = 30 * time.Second
)

const (
	userGRPCURL = "user-service:50051"
	teamGRPCURL = "team-service:50052"
	redisTTL    = 10 * time.Minute
)

func Run(cfg *config.Config) {
	log := logger.InitLogger("services/asset-service/logs/asset-service.log", "asset-service")
	defer log.Sync()

	conn, err := postgres.Connect(cfg.DatabaseURL, log)
	if err != nil {
		log.Error("failed to connect to database", zap.Error(err))
		return
	}
	defer postgres.Close(conn, log)

	redisClient, err := redis.Init(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB, log)
	if err != nil {
		log.Fatal("failed to init redis:", zap.Error(err))
	}
	defer redisClient.Close(log)

	assetDBRepo := repository.NewAssetDBRepo(conn)
	assetCacheRepo := repository.NewAssetCache(redisClient.Client, redisTTL)

	assetRepo := repository.NewCachedAssetRepo(assetDBRepo, assetCacheRepo)

	// Initialize Kafka producer
	kafkaProducer, err := kafka.NewProducer([]string{"kafka:29092"})
	if err != nil {
		log.Fatal("failed to create kafka producer", zap.Error(err))
	}
	defer kafkaProducer.Close()

	userClient := client.NewUserGRPCClient(userGRPCURL)
	teamClient := client.NewTeamGRPCClient(teamGRPCURL)

	assetSvc := service.NewAssetService(assetRepo, *userClient, *teamClient, kafkaProducer)

	// Health checks
	dbCheck := func() error {
		sqlDB, err := conn.DB()
		if err != nil {
			return err
		}
		return sqlDB.Ping()
	}
	redisCheck := func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		return redisClient.Client.Ping(ctx).Err()
	}

	// Start HTTP server
	httpServer, err := runHTTPServer(assetSvc, *userClient, log, dbCheck, redisCheck)
	if err != nil {
		log.Fatal("Failed to start HTTP server", zap.Error(err))
	}

	// Create graceful shutdown server
	gracefulServer := graceful.NewGracefulServer(GRATEFUL_TIMEOUT)

	// Add services for graceful shutdown
	gracefulServer.AddService(&gracefulService{
		name:       "Asset Service",
		httpServer: httpServer,
		log:        log,
	})

	// Start graceful shutdown listener
	gracefulServer.Start()
}

func runHTTPServer(assetService service.AssetService, userClient client.UserGRPCClient, log *zap.Logger, dbCheck func() error, redisCheck func() error) (*http.Server, error) {
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
		if redisCheck != nil {
			if err := redisCheck(); err != nil {
				health["redis"] = "down"
				health["redis_error"] = err.Error()
			} else {
				health["redis"] = "up"
			}
		}
		c.JSON(200, health)
	})

	assetHandler := handler.NewAssetHandler(assetService)
	api := r.Group("/api/v1")
	assetGroup := api.Group("/assets")
	assetGroup.Use(middleware.AuthMiddleware(userClient.Client))
	{
		assetGroup.POST("/folders", assetHandler.CreateFolder)
		assetGroup.GET("/folders/:folderId", assetHandler.GetFolder)
		assetGroup.PUT("/folders/:folderId", assetHandler.GetFolder)
		assetGroup.DELETE("/folders/:folderId", assetHandler.GetFolder)

		assetGroup.POST("/notes", assetHandler.CreateNote)
		assetGroup.GET("/notes/:noteId", assetHandler.GetNote)
		assetGroup.PUT("/notes/:noteId", assetHandler.GetNote)
		assetGroup.DELETE("/notes/:noteId", assetHandler.GetNote)

		assetGroup.POST("/folders/:folderId/share", assetHandler.ShareFolder)
		assetGroup.DELETE("/folders/:folderId/share/:userId", assetHandler.RevokeFolderShare)
		assetGroup.POST("/notes/:noteId/share", assetHandler.ShareNote)
		assetGroup.DELETE("/notes/:noteId/share/:userId", assetHandler.RevokeNoteShare)

		assetGroup.GET("/teams/:teamId/assets", assetHandler.GetTeamAssets)
		assetGroup.GET("/users/:userId/assets", assetHandler.GetUserAssets)
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
	name       string
	httpServer *http.Server
	log        *zap.Logger
}

func (s *gracefulService) Shutdown(ctx context.Context) error {
	s.log.Info("Shutting down " + s.name)

	// Shutdown HTTP server
	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.log.Error("HTTP server shutdown error", zap.Error(err))
	}

	s.log.Info(s.name + " shutdown completed")
	return nil
}
