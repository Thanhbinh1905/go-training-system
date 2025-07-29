package main

import (
	"log"

	"github.com/Thanhbinh1905/go-training-system/services/team-service/config"
	"github.com/Thanhbinh1905/go-training-system/services/team-service/internal/authclient"
	"github.com/Thanhbinh1905/go-training-system/services/team-service/internal/handler"
	"github.com/Thanhbinh1905/go-training-system/services/team-service/internal/repository"
	"github.com/Thanhbinh1905/go-training-system/services/team-service/internal/service"
	"github.com/Thanhbinh1905/go-training-system/shared/db"
	"github.com/Thanhbinh1905/go-training-system/shared/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const graphqlURL = "http://user-service:8080/graphql"

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("failed to Load config")
	}

	logger.InitLogger(cfg.Production)

	conn, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		logger.Log.Error("failed to connect to database", zap.Error(err))
	}
	defer db.Close(conn)

	teamRepo := repository.NewTeamRepository(conn)
	authClient := authclient.NewAuthServiceClient(graphqlURL)
	teamService := service.NewTeamService(teamRepo, *authClient)
	teamHandler := handler.NewTeamHandler(teamService)

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	r.POST("/teams", teamHandler.CreateTeam)
	r.POST("/teams/:teamID/managers", teamHandler.AddManager)
	r.POST("/teams/:teamID/members", teamHandler.AddMember)
	r.DELETE("/teams/:teamID/managers/:managerID", teamHandler.RemoveManager)
	r.DELETE("/teams/:teamID/members/:memberID", teamHandler.RemoveMember)

	logger.Log.Info("Starting server on port " + "8080")
	r.Run()
}
