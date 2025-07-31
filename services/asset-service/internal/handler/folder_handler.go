package handler

import (
	"net/http"

	"github.com/Thanhbinh1905/go-training-system/services/asset-service/internal/model"
	"github.com/Thanhbinh1905/go-training-system/services/asset-service/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FolderHandler struct {
	folderService service.FolderService
}

func NewFolderHandler(folderService service.FolderService) *FolderHandler {
	return &FolderHandler{folderService}
}

// POST /folders
func (h *FolderHandler) CreateFolder(c *gin.Context) {
	var folder model.Folder
	if err := c.ShouldBindJSON(&folder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.folderService.CreateFolder(c.Request.Context(), &folder); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create folder"})
		return
	}
	c.JSON(http.StatusCreated, folder)
}

// GET /folders/:folderId
func (h *FolderHandler) GetFolder(c *gin.Context) {
	id, err := uuid.Parse(c.Param("folderId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid folder id"})
		return
	}
	folder, err := h.folderService.GetFolder(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "folder not found"})
		return
	}
	c.JSON(http.StatusOK, folder)
}

// PUT /folders/:folderId
func (h *FolderHandler) UpdateFolder(c *gin.Context) {
	id, err := uuid.Parse(c.Param("folderId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid folder id"})
		return
	}

	var updated model.Folder
	if err := c.ShouldBindJSON(&updated); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updated.ID = id

	if err := h.folderService.UpdateFolder(c.Request.Context(), &updated); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update folder"})
		return
	}
	c.JSON(http.StatusOK, updated)
}

// DELETE /folders/:folderId
func (h *FolderHandler) DeleteFolder(c *gin.Context) {
	id, err := uuid.Parse(c.Param("folderId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid folder id"})
		return
	}
	if err := h.folderService.DeleteFolder(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete folder"})
		return
	}
	c.Status(http.StatusNoContent)
}
