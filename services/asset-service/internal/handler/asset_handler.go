package handler

import (
	"fmt"
	"net/http"

	"github.com/Thanhbinh1905/go-training-system/services/asset-service/internal/service"
	"github.com/Thanhbinh1905/go-training-system/shared/contextkey"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AssetHandler struct {
	assetService service.AssetService
}

func NewAssetHandler(assetService service.AssetService) *AssetHandler {
	return &AssetHandler{assetService}
}

// GET /users/:userId/assets
func (h *AssetHandler) GetUserAssets(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	folders, notes, err := h.assetService.GetUserAssets(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user assets"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"folders": folders, "notes": notes})
}

// GET /teams/:teamId/assets
func (h *AssetHandler) GetTeamAssets(c *gin.Context) {
	teamID, err := uuid.Parse(c.Param("teamId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid team ID"})
		return
	}

	folders, notes, err := h.assetService.GetTeamAssets(c.Request.Context(), teamID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get team assets"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"folders": folders, "notes": notes})
}

func GetUserIDFromContext(c *gin.Context) (uuid.UUID, error) {
	ctx := c.Request.Context()
	userIDStr, ok := ctx.Value(contextkey.CtxUserIDKey()).(string)
	if !ok || userIDStr == "" {
		return uuid.Nil, fmt.Errorf("userID not found in context")
	}
	return uuid.Parse(userIDStr)
}

func respondError(c *gin.Context, status int, msg string) {
	c.JSON(status, gin.H{"error": msg})
}

func respondMessage(c *gin.Context, status int, msg string) {
	c.JSON(status, gin.H{"message": msg})
}
