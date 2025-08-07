package handler

import (
	"fmt"
	"net/http"

	"github.com/Thanhbinh1905/go-training-system/services/team-service/internal/dto"
	"github.com/Thanhbinh1905/go-training-system/services/team-service/internal/service"
	"github.com/Thanhbinh1905/go-training-system/services/team-service/pkg/contextkey"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TeamHandler struct {
	service service.TeamService
}

func NewTeamHandler(service service.TeamService) *TeamHandler {
	return &TeamHandler{
		service: service,
	}
}

func (h *TeamHandler) CreateTeam(c *gin.Context) {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req dto.CreateTeamInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateTeam(c.Request.Context(), userID, &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Team created successfully"})
}

func (h *TeamHandler) AddManager(c *gin.Context) {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	teamIDStr := c.Param("teamID")
	teamID, err := uuid.Parse(teamIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid team ID"})
		return
	}

	var req dto.AddManagerInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.AddManager(c.Request.Context(), userID, teamID, req.ManagerIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Manager added successfully"})
}

func (s *TeamHandler) AddMember(c *gin.Context) {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	teamID := c.Param("teamID")

	if teamID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Team ID is required"})
		return
	}

	var req *dto.AddMemberInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.service.AddMember(c.Request.Context(), userID, uuid.MustParse(teamID), req.MemberIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Member added successfully"})
}

func (s *TeamHandler) RemoveManager(c *gin.Context) {
	teamID := c.Param("teamID")
	managerID := c.Param("managerID")

	if teamID == "" || managerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Team ID and Manager ID are required"})
		return
	}

	if err := s.service.RemoveManager(c.Request.Context(), uuid.MustParse(teamID), uuid.MustParse(managerID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Manager removed successfully"})
}

func (s *TeamHandler) RemoveMember(c *gin.Context) {
	teamID := c.Param("teamID")
	memberID := c.Param("memberID")

	if teamID == "" || memberID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Team ID and Member ID are required"})
		return
	}

	if err := s.service.RemoveMember(c.Request.Context(), uuid.MustParse(teamID), uuid.MustParse(memberID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Member removed successfully"})
}

func GetUserIDFromContext(c *gin.Context) (uuid.UUID, error) {
	ctx := c.Request.Context()
	userIDStr, ok := ctx.Value(contextkey.CtxUserIDKey()).(string)
	if !ok || userIDStr == "" {
		return uuid.Nil, fmt.Errorf("userID not found in context")
	}
	return uuid.Parse(userIDStr)
}
