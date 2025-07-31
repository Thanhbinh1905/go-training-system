package main

import (
	"log"

	"github.com/Thanhbinh1905/go-training-system/services/asset-service/config"
	"github.com/Thanhbinh1905/go-training-system/services/asset-service/internal/handler"
	"github.com/Thanhbinh1905/go-training-system/services/asset-service/internal/repository"
	"github.com/Thanhbinh1905/go-training-system/services/asset-service/internal/service"
	"github.com/Thanhbinh1905/go-training-system/services/asset-service/pkg/logger"
	"github.com/Thanhbinh1905/go-training-system/shared/db"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

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

	r := gin.Default()

	folderRepo := repository.NewFolderRepository(conn)
	folderService := service.NewFolderService(folderRepo)
	folderHandler := handler.NewFolderHandler(folderService)

	r.POST("/folders", folderHandler.CreateFolder)
	r.GET("/folders/:folderId", folderHandler.GetFolder)
	r.PUT("/folders/:folderId", folderHandler.UpdateFolder)
	r.DELETE("/folders/:folderId", folderHandler.DeleteFolder)

	noteRepo := repository.NewNoteRepository(conn)
	noteService := service.NewNoteService(noteRepo)
	noteHandler := handler.NewNoteHandler(noteService)

	r.POST("/folders/:folderId/notes", noteHandler.CreateNote)
	r.GET("/notes/:noteId", noteHandler.GetNote)
	r.PUT("/notes/:noteId", noteHandler.UpdateNote)
	r.DELETE("/notes/:noteId", noteHandler.DeleteNote)

	folderShareRepo := repository.NewFolderShareRepository(conn)
	noteShareRepo := repository.NewNoteShareRepository(conn)
	sharingService := service.NewSharingService(folderShareRepo, noteShareRepo)
	sharingHandler := handler.NewSharingHandler(sharingService)

	r.POST("/folders/:folderId/share", sharingHandler.ShareFolder)
	r.DELETE("/folders/:folderId/share/:userId", sharingHandler.RevokeFolderShare)
	r.POST("/notes/:noteId/share", sharingHandler.ShareNote)
	r.DELETE("/notes/:noteId/share/:userId", sharingHandler.RevokeNoteShare)

	assetService := service.NewAssetService(folderRepo, noteRepo, folderShareRepo, noteShareRepo)
	assetHandler := handler.NewAssetHandler(assetService)

	r.GET("/users/:userId/assets", assetHandler.GetUserAssets)
	r.GET("/teams/:teamId/assets", assetHandler.GetTeamAssets)
}
