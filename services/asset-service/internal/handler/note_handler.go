package handler

import (
	"fmt"
	"net/http"

	"github.com/Thanhbinh1905/go-training-system/services/asset-service/internal/model"
	"github.com/Thanhbinh1905/go-training-system/services/asset-service/internal/service"
	"github.com/Thanhbinh1905/go-training-system/shared/contextkey"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type NoteHandler struct {
	noteService service.NoteService
}

func NewNoteHandler(noteService service.NoteService) *NoteHandler {
	return &NoteHandler{noteService}
}

func getUserIDFromContext(c *gin.Context) (uuid.UUID, error) {
	userIDStr, ok := c.Request.Context().Value(contextkey.CtxUserIDKey()).(string)
	if !ok || userIDStr == "" {
		return uuid.Nil, fmt.Errorf("userID not found in context")
	}
	return uuid.Parse(userIDStr)
}

// POST /folders/:folderId/notes
func (h *NoteHandler) CreateNote(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		respondError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	folderID, err := uuid.Parse(c.Param("folderId"))
	if err != nil {
		respondError(c, http.StatusBadRequest, "invalid folder ID")
		return
	}

	var note model.Note
	if err := c.ShouldBindJSON(&note); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	note.FolderID = folderID
	note.OwnerID = userID

	if err := h.noteService.CreateNote(c, &note); err != nil {
		respondError(c, http.StatusInternalServerError, "failed to create note")
		return
	}
	c.JSON(http.StatusCreated, note)
}

// GET /notes/:noteId
func (h *NoteHandler) GetNote(c *gin.Context) {
	_, err := getUserIDFromContext(c)
	if err != nil {
		respondError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	id, err := uuid.Parse(c.Param("noteId"))
	if err != nil {
		respondError(c, http.StatusBadRequest, "invalid note ID")
		return
	}

	note, err := h.noteService.GetNote(c, id)
	if err != nil {
		respondError(c, http.StatusNotFound, "note not found")
		return
	}
	c.JSON(http.StatusOK, note)
}

// PUT /notes/:noteId
func (h *NoteHandler) UpdateNote(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		respondError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	id, err := uuid.Parse(c.Param("noteId"))
	if err != nil {
		respondError(c, http.StatusBadRequest, "invalid note ID")
		return
	}

	var note model.Note
	if err := c.ShouldBindJSON(&note); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	note.ID = id
	note.OwnerID = userID

	if err := h.noteService.UpdateNote(c, &note); err != nil {
		respondError(c, http.StatusInternalServerError, "failed to update note")
		return
	}
	c.JSON(http.StatusOK, note)
}

// DELETE /notes/:noteId
func (h *NoteHandler) DeleteNote(c *gin.Context) {
	_, err := getUserIDFromContext(c)
	if err != nil {
		respondError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	id, err := uuid.Parse(c.Param("noteId"))
	if err != nil {
		respondError(c, http.StatusBadRequest, "invalid note ID")
		return
	}

	if err := h.noteService.DeleteNote(c, id); err != nil {
		respondError(c, http.StatusInternalServerError, "failed to delete note")
		return
	}
	c.Status(http.StatusNoContent)
}
