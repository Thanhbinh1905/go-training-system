// internal/handler/note_handler.go
package handler

import (
	"net/http"

	"github.com/Thanhbinh1905/go-training-system/services/asset-service/internal/model"
	"github.com/Thanhbinh1905/go-training-system/services/asset-service/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type NoteHandler struct {
	noteService service.NoteService
}

func NewNoteHandler(noteService service.NoteService) *NoteHandler {
	return &NoteHandler{noteService}
}

// POST /folders/:folderId/notes
func (h *NoteHandler) CreateNote(c *gin.Context) {
	folderID, err := uuid.Parse(c.Param("folderId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid folder ID"})
		return
	}

	var note model.Note
	if err := c.ShouldBindJSON(&note); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	note.FolderID = folderID

	if err := h.noteService.CreateNote(c.Request.Context(), &note); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create note"})
		return
	}
	c.JSON(http.StatusCreated, note)
}

// GET /notes/:noteId
func (h *NoteHandler) GetNote(c *gin.Context) {
	id, err := uuid.Parse(c.Param("noteId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid note ID"})
		return
	}
	note, err := h.noteService.GetNote(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "note not found"})
		return
	}
	c.JSON(http.StatusOK, note)
}

// PUT /notes/:noteId
func (h *NoteHandler) UpdateNote(c *gin.Context) {
	id, err := uuid.Parse(c.Param("noteId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid note ID"})
		return
	}

	var note model.Note
	if err := c.ShouldBindJSON(&note); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	note.ID = id

	if err := h.noteService.UpdateNote(c.Request.Context(), &note); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update note"})
		return
	}
	c.JSON(http.StatusOK, note)
}

// DELETE /notes/:noteId
func (h *NoteHandler) DeleteNote(c *gin.Context) {
	id, err := uuid.Parse(c.Param("noteId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid note ID"})
		return
	}
	if err := h.noteService.DeleteNote(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete note"})
		return
	}
	c.Status(http.StatusNoContent)
}
