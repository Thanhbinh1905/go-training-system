package handler

import (
	"net/http"

	"github.com/Thanhbinh1905/go-training-system/services/asset-service/internal/model"
	"github.com/Thanhbinh1905/go-training-system/services/asset-service/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SharingHandler struct {
	sharingService service.SharingService
}

func NewSharingHandler(sharingService service.SharingService) *SharingHandler {
	return &SharingHandler{sharingService}
}

// POST /folders/:folderId/share
func (h *SharingHandler) ShareFolder(c *gin.Context) {
	folderID, err := uuid.Parse(c.Param("folderId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid folder ID"})
		return
	}

	var share model.FolderShare
	if err := c.ShouldBindJSON(&share); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	share.FolderID = folderID
	if err := h.sharingService.ShareFolder(c.Request.Context(), &share); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to share folder"})
		return
	}
	c.JSON(http.StatusCreated, share)
}

// DELETE /folders/:folderId/share/:userId
func (h *SharingHandler) RevokeFolderShare(c *gin.Context) {
	folderID, err := uuid.Parse(c.Param("folderId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid folder ID"})
		return
	}
	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	if err := h.sharingService.RevokeFolderShare(c.Request.Context(), folderID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to revoke share"})
		return
	}
	c.Status(http.StatusNoContent)
}

// POST /notes/:noteId/share
func (h *SharingHandler) ShareNote(c *gin.Context) {
	noteID, err := uuid.Parse(c.Param("noteId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid note ID"})
		return
	}

	var share model.NoteShare
	if err := c.ShouldBindJSON(&share); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	share.NoteID = noteID
	if err := h.sharingService.ShareNote(c.Request.Context(), &share); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to share note"})
		return
	}
	c.JSON(http.StatusCreated, share)
}

// DELETE /notes/:noteId/share/:userId
func (h *SharingHandler) RevokeNoteShare(c *gin.Context) {
	noteID, err := uuid.Parse(c.Param("noteId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid note ID"})
		return
	}
	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	if err := h.sharingService.RevokeNoteShare(c.Request.Context(), noteID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to revoke share"})
		return
	}
	c.Status(http.StatusNoContent)
}
