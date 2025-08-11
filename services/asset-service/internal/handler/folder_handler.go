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
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		respondError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	var folder model.Folder
	if err := c.ShouldBindJSON(&folder); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Set ownerID nếu cần
	folder.OwnerID = userID

	if err := h.folderService.CreateFolder(c, &folder); err != nil {
		respondError(c, http.StatusInternalServerError, "failed to create folder")
		return
	}
	c.JSON(http.StatusCreated, folder)
}

// GET /folders/:folderId
func (h *FolderHandler) GetFolder(c *gin.Context) {
	_, err := GetUserIDFromContext(c)
	if err != nil {
		respondError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	id, err := uuid.Parse(c.Param("folderId"))
	if err != nil {
		respondError(c, http.StatusBadRequest, "invalid folder id")
		return
	}

	folder, err := h.folderService.GetFolder(c, id)
	if err != nil {
		respondError(c, http.StatusNotFound, "folder not found")
		return
	}
	c.JSON(http.StatusOK, folder)
}

// PUT /folders/:folderId
func (h *FolderHandler) UpdateFolder(c *gin.Context) {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		respondError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	id, err := uuid.Parse(c.Param("folderId"))
	if err != nil {
		respondError(c, http.StatusBadRequest, "invalid folder id")
		return
	}

	var updated model.Folder
	if err := c.ShouldBindJSON(&updated); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	updated.ID = id
	updated.OwnerID = userID

	if err := h.folderService.UpdateFolder(c, &updated); err != nil {
		respondError(c, http.StatusInternalServerError, "failed to update folder")
		return
	}
	c.JSON(http.StatusOK, updated)
}

// DELETE /folders/:folderId
func (h *FolderHandler) DeleteFolder(c *gin.Context) {
	_, err := GetUserIDFromContext(c)
	if err != nil {
		respondError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	id, err := uuid.Parse(c.Param("folderId"))
	if err != nil {
		respondError(c, http.StatusBadRequest, "invalid folder id")
		return
	}

	if err := h.folderService.DeleteFolder(c, id); err != nil {
		respondError(c, http.StatusInternalServerError, "failed to delete folder")
		return
	}
	c.Status(http.StatusNoContent)
}
