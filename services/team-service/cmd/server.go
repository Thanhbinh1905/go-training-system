package main

import (
	"log"
	"time"

	"github.com/Thanhbinh1905/go-training-system/services/team-service/config"
	"github.com/Thanhbinh1905/go-training-system/services/team-service/internal/client"
	"github.com/Thanhbinh1905/go-training-system/services/team-service/internal/handler"
	"github.com/Thanhbinh1905/go-training-system/services/team-service/internal/middleware"
	"github.com/Thanhbinh1905/go-training-system/services/team-service/internal/repository"
	"github.com/Thanhbinh1905/go-training-system/services/team-service/internal/service"
	"github.com/Thanhbinh1905/go-training-system/services/team-service/pkg/logger"
	"github.com/Thanhbinh1905/go-training-system/shared/db"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"go.uber.org/zap"
)

const gRPC_URL = "user-service:9090"

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("failed to load config")
	}

	log := logger.NewLogger("logs/team-service.log", "team-service")
	defer log.Sync()

	conn, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Error("failed to connect to database", zap.Error(err))
		return
	}
	defer db.Close(conn)

	teamRepo := repository.NewTeamRepository(conn)
	userClient := client.NewUserGRPCClient(gRPC_URL)
	teamService := service.NewTeamService(teamRepo, *userClient)
	teamHandler := handler.NewTeamHandler(teamService)

	r := gin.Default()
	r.Use(ginzap.Ginzap(log, time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(log, true))

	// Prometheus metrics
	p := ginprometheus.NewPrometheus("team_service")
	p.Use(r)

	// Healthcheck
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Authenticated routes
	teamGroup := r.Group("/api/v1/teams")
	teamGroup.Use(middleware.AuthMiddleware(userClient.Client, true)) // manager-only
	{
		teamGroup.POST("/", teamHandler.CreateTeam)
		teamGroup.POST("/:teamID/managers", teamHandler.AddManager)
		teamGroup.POST("/:teamID/members", teamHandler.AddMember)
		teamGroup.DELETE("/:teamID/managers/:managerID", teamHandler.RemoveManager)
		teamGroup.DELETE("/:teamID/members/:memberID", teamHandler.RemoveMember)
	}

	log.Info("Starting server on port 8080")
	r.Run(":8080")
}
