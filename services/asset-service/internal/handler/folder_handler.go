package handler

import (
	"net/http"

	"github.com/Thanhbinh1905/go-training-system/services/asset-service/internal/dto"
	"github.com/Thanhbinh1905/go-training-system/services/asset-service/internal/service"
	ctxKey "github.com/Thanhbinh1905/go-training-system/shared/contextkey"
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
	userID, err := ctxKey.GetUserIDFromContext(c.Request.Context())
	if err != nil {
		respondError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	var input dto.CreateFolderInput
	if err := c.ShouldBindJSON(&input); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.folderService.CreateFolder(c.Request.Context(), userID, &input); err != nil {
		respondError(c, http.StatusInternalServerError, "failed to create folder")
		return
	}
	c.JSON(http.StatusCreated, input)
}

// GET /folders/:folderId
func (h *FolderHandler) GetFolder(c *gin.Context) {
	id, err := uuid.Parse(c.Param("folderId"))
	if err != nil {
		respondError(c, http.StatusBadRequest, "invalid folder id")
		return
	}

	folder, err := h.folderService.GetFolderByID(c.Request.Context(), id)
	if err != nil {
		respondError(c, http.StatusNotFound, "folder not found")
		return
	}
	c.JSON(http.StatusOK, folder)
}

// PUT /folders/:folderId
func (h *FolderHandler) UpdateFolder(c *gin.Context) {
	userID, err := ctxKey.GetUserIDFromContext(c.Request.Context())
	if err != nil {
		respondError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	folderID, err := uuid.Parse(c.Param("folderId"))
	if err != nil {
		respondError(c, http.StatusBadRequest, "invalid folder id")
		return
	}

	var updated dto.UpdateFolderInput
	if err := c.ShouldBindJSON(&updated); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.folderService.UpdateFolder(c.Request.Context(), userID, folderID, &updated); err != nil {
		if err.Error() == "forbidden: not the owner" {
			respondError(c, http.StatusForbidden, err.Error())
			return
		}
		respondError(c, http.StatusInternalServerError, "failed to update folder")
		return
	}
	c.JSON(http.StatusOK, updated)
}

// DELETE /folders/:folderId
func (h *FolderHandler) DeleteFolder(c *gin.Context) {
	userID, err := ctxKey.GetUserIDFromContext(c.Request.Context())
	if err != nil {
		respondError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	id, err := uuid.Parse(c.Param("folderId"))
	if err != nil {
		respondError(c, http.StatusBadRequest, "invalid folder id")
		return
	}

	if err := h.folderService.DeleteFolder(c.Request.Context(), userID, id); err != nil {
		if err.Error() == "forbidden: not the owner" {
			respondError(c, http.StatusForbidden, err.Error())
			return
		}
		respondError(c, http.StatusInternalServerError, "failed to delete folder")
		return
	}
	c.Status(http.StatusNoContent)
}
