package handler

import (
	"net/http"

	"github.com/Thanhbinh1905/go-training-system/services/asset-service/internal/dto"
	"github.com/Thanhbinh1905/go-training-system/services/asset-service/internal/service"
	ctxKey "github.com/Thanhbinh1905/go-training-system/shared/contextkey"
	errorhelper "github.com/Thanhbinh1905/go-training-system/shared/errors"
	"github.com/Thanhbinh1905/go-training-system/shared/httpresponse"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AssetHandler struct {
	assetService service.AssetService
}

func NewAssetHandler(assetService service.AssetService) *AssetHandler {
	return &AssetHandler{assetService}
}

// -------------------- Folder --------------------

// POST /folders
func (h *AssetHandler) CreateFolder(c *gin.Context) {
	userID, err := ctxKey.GetUserIDFromContext(c.Request.Context())
	if err != nil {
		httpresponse.RespondError(c, err)
		return
	}

	var input dto.CreateFolderInput
	if err := c.ShouldBindJSON(&input); err != nil {
		httpresponse.RespondError(c, errorhelper.ErrInvalidInput)
		return
	}

	if err := h.assetService.CreateFolder(c.Request.Context(), userID, &input); err != nil {
		httpresponse.RespondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, input)
}

// GET /folders/:folderId
func (h *AssetHandler) GetFolder(c *gin.Context) {
	userID, err := ctxKey.GetUserIDFromContext(c.Request.Context())
	if err != nil {
		httpresponse.RespondError(c, err)
		return
	}

	folderID, err := uuid.Parse(c.Param("folderId"))
	if err != nil {
		httpresponse.RespondError(c, errorhelper.ErrInvalidInput)
		return
	}

	folder, err := h.assetService.GetFolderByID(c.Request.Context(), userID, folderID)
	if err != nil {
		httpresponse.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, folder)
}

// PUT /folders/:folderId
func (h *AssetHandler) UpdateFolder(c *gin.Context) {
	userID, err := ctxKey.GetUserIDFromContext(c.Request.Context())
	if err != nil {
		httpresponse.RespondError(c, err)
		return
	}

	folderID, err := uuid.Parse(c.Param("folderId"))
	if err != nil {
		httpresponse.RespondError(c, errorhelper.ErrInvalidInput)
		return
	}

	var updated dto.UpdateFolderInput
	if err := c.ShouldBindJSON(&updated); err != nil {
		httpresponse.RespondError(c, errorhelper.ErrInvalidInput)
		return
	}

	if err := h.assetService.UpdateFolder(c.Request.Context(), userID, folderID, &updated); err != nil {
		httpresponse.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, updated)
}

// DELETE /folders/:folderId
func (h *AssetHandler) DeleteFolder(c *gin.Context) {
	userID, err := ctxKey.GetUserIDFromContext(c.Request.Context())
	if err != nil {
		httpresponse.RespondError(c, err)
		return
	}

	folderID, err := uuid.Parse(c.Param("folderId"))
	if err != nil {
		httpresponse.RespondError(c, errorhelper.ErrInvalidInput)
		return
	}

	if err := h.assetService.DeleteFolder(c.Request.Context(), userID, folderID); err != nil {
		httpresponse.RespondError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// -------------------- Note --------------------
// POST /notes
func (h *AssetHandler) CreateNote(c *gin.Context) {
	userID, err := ctxKey.GetUserIDFromContext(c.Request.Context())
	if err != nil {
		httpresponse.RespondError(c, err) // service trả ErrUnauthorized
		return
	}

	var input dto.CreateNoteInput
	if err := c.ShouldBindJSON(&input); err != nil {
		httpresponse.RespondError(c, errorhelper.ErrInvalidInput)
		return
	}

	if err := h.assetService.CreateNote(c.Request.Context(), userID, &input); err != nil {
		httpresponse.RespondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, input)
}

// GET /notes/:noteId
func (h *AssetHandler) GetNote(c *gin.Context) {
	userID, err := ctxKey.GetUserIDFromContext(c.Request.Context())
	if err != nil {
		httpresponse.RespondError(c, err)
		return
	}

	noteID, err := uuid.Parse(c.Param("noteId"))
	if err != nil {
		httpresponse.RespondError(c, errorhelper.ErrInvalidInput)
		return
	}

	note, err := h.assetService.GetNoteByID(c.Request.Context(), userID, noteID)
	if err != nil {
		httpresponse.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, note)
}

// PUT /notes/:noteId
func (h *AssetHandler) UpdateNote(c *gin.Context) {
	userID, err := ctxKey.GetUserIDFromContext(c.Request.Context())
	if err != nil {
		httpresponse.RespondError(c, err)
		return
	}

	noteID, err := uuid.Parse(c.Param("noteId"))
	if err != nil {
		httpresponse.RespondError(c, errorhelper.ErrInvalidInput)
		return
	}

	var input dto.UpdateNoteInput
	if err := c.ShouldBindJSON(&input); err != nil {
		httpresponse.RespondError(c, errorhelper.ErrInvalidInput)
		return
	}

	if err := h.assetService.UpdateNote(c.Request.Context(), userID, noteID, &input); err != nil {
		httpresponse.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, input)
}

// DELETE /notes/:noteId
func (h *AssetHandler) DeleteNote(c *gin.Context) {
	userID, err := ctxKey.GetUserIDFromContext(c.Request.Context())
	if err != nil {
		httpresponse.RespondError(c, err)
		return
	}

	noteID, err := uuid.Parse(c.Param("noteId"))
	if err != nil {
		httpresponse.RespondError(c, errorhelper.ErrInvalidInput)
		return
	}

	if err := h.assetService.DeleteNote(c.Request.Context(), userID, noteID); err != nil {
		httpresponse.RespondError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// -------------------- Sharing --------------------

// POST /folders/:folderId/share
func (h *AssetHandler) ShareFolder(c *gin.Context) {
	var input dto.CreateFolderShareInput
	if err := c.ShouldBindJSON(&input); err != nil {
		httpresponse.RespondError(c, err)
		return
	}

	folderID, err := uuid.Parse(c.Param("folderId"))
	if err != nil {
		httpresponse.RespondError(c, err)
		return
	}

	userID, err := ctxKey.GetUserIDFromContext(c.Request.Context())
	if err != nil {
		httpresponse.RespondError(c, err) // service trả ErrUnauthorized
		return
	}

	if err := h.assetService.ShareFolder(c.Request.Context(), userID, folderID, &input); err != nil {
		httpresponse.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "folder shared successfully"})
}

// POST /notes/:noteId/share
func (h *AssetHandler) ShareNote(c *gin.Context) {
	var input dto.CreateNoteShareInput
	if err := c.ShouldBindJSON(&input); err != nil {
		httpresponse.RespondError(c, err)
		return
	}

	noteID, err := uuid.Parse(c.Param("noteId"))
	if err != nil {
		httpresponse.RespondError(c, err)
		return
	}

	userID, err := ctxKey.GetUserIDFromContext(c.Request.Context())
	if err != nil {
		httpresponse.RespondError(c, err) // service trả ErrUnauthorized
		return
	}

	if err := h.assetService.ShareNote(c.Request.Context(), userID, noteID, &input); err != nil {
		httpresponse.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "note shared successfully"})
}

// DELETE /folders/:folderId/share/:userId
func (h *AssetHandler) RevokeFolderShare(c *gin.Context) {
	folderID, err := uuid.Parse(c.Param("folderId"))
	if err != nil {
		httpresponse.RespondError(c, err)
		return
	}

	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		httpresponse.RespondError(c, err)
		return
	}

	if err := h.assetService.RevokeFolderShare(c.Request.Context(), folderID, userID); err != nil {
		httpresponse.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "folder share revoked"})
}

// DELETE /notes/:noteId/share/:userId
func (h *AssetHandler) RevokeNoteShare(c *gin.Context) {
	noteID, err := uuid.Parse(c.Param("noteId"))
	if err != nil {
		httpresponse.RespondError(c, err)
		return
	}

	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		httpresponse.RespondError(c, err)
		return
	}

	if err := h.assetService.RevokeNoteShare(c.Request.Context(), noteID, userID); err != nil {
		httpresponse.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "note share revoked"})
}

// GET /teams/:teamId/assets
func (h *AssetHandler) GetTeamAssets(c *gin.Context) {
	requestedID, err := ctxKey.GetUserIDFromContext(c.Request.Context())
	if err != nil {
		httpresponse.RespondError(c, errorhelper.ErrUnauthorized)
		return
	}

	teamID, err := uuid.Parse(c.Param("teamId"))
	if err != nil {
		httpresponse.RespondError(c, errorhelper.ErrInvalidInput)
		return
	}

	resp, err := h.assetService.GetTeamAssets(c.Request.Context(), requestedID, teamID)
	if err != nil {
		httpresponse.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GET /users/:userId/assets
func (h *AssetHandler) GetUserAssets(c *gin.Context) {
	requestedID, err := ctxKey.GetUserIDFromContext(c.Request.Context())
	if err != nil {
		httpresponse.RespondError(c, errorhelper.ErrUnauthorized)
		return
	}

	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		httpresponse.RespondError(c, errorhelper.ErrInvalidInput)
		return
	}

	resp, err := h.assetService.GetUserAssets(c.Request.Context(), requestedID, userID)
	if err != nil {
		httpresponse.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}
