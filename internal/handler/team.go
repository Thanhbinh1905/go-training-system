package handler

import (
	"net/http"

	"github.com/Thanhbinh1905/go-training-system/internal/dto"
	"github.com/Thanhbinh1905/go-training-system/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TeamHandler struct {
	service service.TeamService
}

func NewTeamHandler(s service.TeamService) *TeamHandler {
	return &TeamHandler{service: s}
}

func (h *TeamHandler) CreateTeam(c *gin.Context) {
	var req dto.CreateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "INVALID_REQUEST",
			"success": false,
			"message": "Invalid request payload",
			"errors":  []string{err.Error()},
		})
		return
	}

	// Lấy userID từ context do middleware đã gán
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    "UNAUTHORIZED",
			"success": false,
			"message": "Unauthorized: userID not found in context",
		})
		return
	}

	userIDStr, ok := userIDRaw.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "INVALID_CONTEXT",
			"success": false,
			"message": "Invalid user ID format in context",
		})
		return
	}

	createdBy, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "INVALID_UUID",
			"success": false,
			"message": "Cannot parse user ID to UUID",
		})
		return
	}

	if err := h.service.CreateTeam(c.Request.Context(), createdBy, &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "CREATE_TEAM_FAILED",
			"success": false,
			"message": "Failed to create team",
			"errors":  []string{err.Error()},
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code":    "TEAM_CREATED",
		"success": true,
		"message": "Team created successfully",
	})
}

func (h *TeamHandler) AddMember(c *gin.Context) {
	teamID, err := uuid.Parse(c.Param("teamId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "INVALID_TEAM_ID",
			"success": false,
			"message": "Invalid team ID",
		})
		return
	}

	var req dto.UserIDRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "INVALID_REQUEST",
			"success": false,
			"message": "Invalid body",
		})
		return
	}

	memberID, err := uuid.Parse(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "INVALID_USER_ID",
			"success": false,
			"message": "User ID must be a valid UUID",
		})
		return
	}

	// Lấy user_id từ context để làm createdBy
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    "UNAUTHORIZED",
			"success": false,
			"message": "Unauthorized: userID not found in context",
		})
		return
	}

	userIDStr, ok := userIDRaw.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "INVALID_CONTEXT",
			"success": false,
			"message": "Invalid user ID format in context",
		})
		return
	}

	createdBy, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "INVALID_UUID",
			"success": false,
			"message": "Cannot parse user ID to UUID",
		})
		return
	}

	err = h.service.AddMember(c.Request.Context(), teamID, memberID, createdBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "ADD_MEMBER_FAILED",
			"success": false,
			"message": "Failed to add member",
			"errors":  []string{err.Error()},
		})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *TeamHandler) AddManager(c *gin.Context) {
	teamIDStr := c.Param("teamId")
	teamID, err := uuid.Parse(teamIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "INVALID_TEAM_ID",
			"success": false,
			"message": "Invalid or missing team_id in URL path",
		})
		return
	}

	var req dto.UserIDRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "INVALID_REQUEST",
			"success": false,
			"message": "Invalid body",
		})
		return
	}

	managerID, err := uuid.Parse(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "INVALID_USER_ID",
			"success": false,
			"message": "User ID must be a valid UUID",
		})
		return
	}

	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    "UNAUTHORIZED",
			"success": false,
			"message": "Unauthorized: userID not found in context",
		})
		return
	}

	userIDStr, ok := userIDRaw.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "INVALID_CONTEXT",
			"success": false,
			"message": "Invalid user ID format in context",
		})
		return
	}

	createdBy, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "INVALID_UUID",
			"success": false,
			"message": "Cannot parse user ID to UUID",
		})
		return
	}

	err = h.service.AddManager(c.Request.Context(), teamID, managerID, createdBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "ADD_MANAGER_FAILED",
			"success": false,
			"message": "Failed to add manager",
			"errors":  []string{err.Error()},
		})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *TeamHandler) RemoveMember(c *gin.Context) {
	teamID, _ := uuid.Parse(c.Param("teamId"))
	userID, _ := uuid.Parse(c.Param("memberId"))

	err := h.service.RemoveMember(c.Request.Context(), teamID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to remove member"})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *TeamHandler) RemoveManager(c *gin.Context) {
	teamID, _ := uuid.Parse(c.Param("teamId"))
	userID, _ := uuid.Parse(c.Param("managerId"))

	err := h.service.RemoveManager(c.Request.Context(), teamID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to remove manager"})
		return
	}
	c.Status(http.StatusNoContent)
}
