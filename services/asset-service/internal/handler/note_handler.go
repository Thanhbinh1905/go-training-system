package handler

import (
	"net/http"

	"github.com/Thanhbinh1905/go-training-system/services/asset-service/internal/dto"
	"github.com/Thanhbinh1905/go-training-system/services/asset-service/internal/service"
	ctxKey "github.com/Thanhbinh1905/go-training-system/shared/contextkey"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type NoteHandler struct {
	noteService service.NoteService
}

func NewNoteHandler(noteService service.NoteService) *NoteHandler {
	return &NoteHandler{noteService}
}

// POST /notes
func (h *NoteHandler) CreateNote(c *gin.Context) {
	userID, err := ctxKey.GetUserIDFromContext(c.Request.Context())
	if err != nil {
		respondError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	var input dto.CreateNoteInput
	if err := c.ShouldBindJSON(&input); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.noteService.CreateNote(c.Request.Context(), userID, &input); err != nil {
		respondError(c, http.StatusInternalServerError, "failed to create note")
		return
	}
	c.JSON(http.StatusCreated, input)
}

// GET /notes/:noteId
func (h *NoteHandler) GetNote(c *gin.Context) {
	id, err := uuid.Parse(c.Param("noteId"))
	if err != nil {
		respondError(c, http.StatusBadRequest, "invalid note id")
		return
	}

	note, err := h.noteService.GetNote(c.Request.Context(), id)
	if err != nil {
		respondError(c, http.StatusNotFound, "note not found")
		return
	}
	c.JSON(http.StatusOK, note)
}

// PUT /notes/:noteId
func (h *NoteHandler) UpdateNote(c *gin.Context) {
	userID, err := ctxKey.GetUserIDFromContext(c.Request.Context())
	if err != nil {
		respondError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	id, err := uuid.Parse(c.Param("noteId"))
	if err != nil {
		respondError(c, http.StatusBadRequest, "invalid note id")
		return
	}

	var input dto.UpdateNoteInput
	if err := c.ShouldBindJSON(&input); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.noteService.UpdateNote(c.Request.Context(), userID, id, &input); err != nil {
		if err.Error() == "forbidden: not the owner" {
			respondError(c, http.StatusForbidden, err.Error())
			return
		}
		respondError(c, http.StatusInternalServerError, "failed to update note")
		return
	}
	c.JSON(http.StatusOK, input)
}

// DELETE /notes/:noteId
func (h *NoteHandler) DeleteNote(c *gin.Context) {
	userID, err := ctxKey.GetUserIDFromContext(c.Request.Context())
	if err != nil {
		respondError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	id, err := uuid.Parse(c.Param("noteId"))
	if err != nil {
		respondError(c, http.StatusBadRequest, "invalid note id")
		return
	}

	if err := h.noteService.DeleteNote(c.Request.Context(), userID, id); err != nil {
		if err.Error() == "forbidden: not the owner" {
			respondError(c, http.StatusForbidden, err.Error())
			return
		}
		respondError(c, http.StatusInternalServerError, "failed to delete note")
		return
	}
	c.Status(http.StatusNoContent)
}

// GET /folders/:folderId/notes
func (h *NoteHandler) GetNotesByFolder(c *gin.Context) {
	folderID, err := uuid.Parse(c.Param("folderId"))
	if err != nil {
		respondError(c, http.StatusBadRequest, "invalid folder id")
		return
	}

	notes, err := h.noteService.GetNotesByFolder(c.Request.Context(), folderID)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "failed to get notes")
		return
	}
	c.JSON(http.StatusOK, notes)
}
