package main

import (
	"log"
	"time"

	"github.com/Thanhbinh1905/go-training-system/services/team-service/config"
	"github.com/Thanhbinh1905/go-training-system/services/team-service/internal/authclient"
	"github.com/Thanhbinh1905/go-training-system/services/team-service/internal/handler"
	"github.com/Thanhbinh1905/go-training-system/services/team-service/internal/repository"
	"github.com/Thanhbinh1905/go-training-system/services/team-service/internal/service"
	"github.com/Thanhbinh1905/go-training-system/services/team-service/pkg/logger"
	"github.com/Thanhbinh1905/go-training-system/shared/db"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"go.uber.org/zap"
)

const graphqlURL = "http://user-service:8080/graphql"

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("failed to Load config")
	}

	log := logger.NewLogger("logs/team-service.log", "team-service")
	defer log.Sync() // flush

	conn, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Error("failed to connect to database", zap.Error(err))
		return
	}
	defer db.Close(conn)

	teamRepo := repository.NewTeamRepository(conn)
	authClient := authclient.NewAuthServiceClient(graphqlURL)
	teamService := service.NewTeamService(teamRepo, *authClient)
	teamHandler := handler.NewTeamHandler(teamService)

	r := gin.Default()

	r.Use(ginzap.Ginzap(log, time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(log, true))

	p := ginprometheus.NewPrometheus("team_service")
	p.Use(r)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	r.POST("/teams", teamHandler.CreateTeam)
	r.POST("/teams/:teamID/managers", teamHandler.AddManager)
	r.POST("/teams/:teamID/members", teamHandler.AddMember)
	r.DELETE("/teams/:teamID/managers/:managerID", teamHandler.RemoveManager)
	r.DELETE("/teams/:teamID/members/:memberID", teamHandler.RemoveMember)

	log.Info("Starting server on port " + "8080")
	r.Run()
}
