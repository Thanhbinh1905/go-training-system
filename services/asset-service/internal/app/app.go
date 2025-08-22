package app

import (
	"time"

	"github.com/Thanhbinh1905/go-training-system/services/asset-service/config"
	"github.com/Thanhbinh1905/go-training-system/services/asset-service/internal/client"
	"github.com/Thanhbinh1905/go-training-system/services/asset-service/internal/handler"
	"github.com/Thanhbinh1905/go-training-system/services/asset-service/internal/middleware"
	"github.com/Thanhbinh1905/go-training-system/services/asset-service/internal/repository"
	"github.com/Thanhbinh1905/go-training-system/services/asset-service/internal/service"
	"github.com/Thanhbinh1905/go-training-system/shared/db"
	"github.com/Thanhbinh1905/go-training-system/shared/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	ginzap "github.com/gin-contrib/zap"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

const (
	userGRPCURL = "user-service:50051"
	teamGRPCURL = "team-service:50052"
)

func Run(cfg *config.Config) {
	log := logger.InitLogger("logs/asset-service.log", "asset-service")
	defer log.Sync()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("failed to Load config")
	}

	conn, err := db.Connect(cfg.DatabaseURL, log)
	if err != nil {
		log.Error("failed to connect to database", zap.Error(err))
		return
	}
	defer db.Close(conn, log)

	assetRepo := repository.NewAssetRepo(conn)

	userClient := client.NewUserGRPCClient(userGRPCURL)
	teamClient := client.NewTeamGRPCClient(teamGRPCURL)

	assetSvc := service.NewAssetService(assetRepo, *userClient, *teamClient)

	runHTTPServer(assetSvc, *userClient, *teamClient, log)
}

func runHTTPServer(assetService service.AssetService, userClient client.UserGRPCClient, teamClient client.TeamGRPCClient, log *zap.Logger) {
	r := gin.Default()
	r.Use(ginzap.Ginzap(log, time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(log, true))

	// Metrics
	p := ginprometheus.NewPrometheus("team_service")
	p.Use(r)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
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

	log.Info("HTTP server started", zap.String("port", "8080"))
	if err := r.Run(":8080"); err != nil {
		log.Fatal("failed to run HTTP server", zap.Error(err))
	}

}
