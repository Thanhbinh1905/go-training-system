package handler

import (
	"net/http"

	"github.com/Thanhbinh1905/go-training-system/services/team-service/internal/dto"
	"github.com/Thanhbinh1905/go-training-system/services/team-service/internal/service"
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
	var req dto.CreateTeamInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
		return
	}

	if err := h.service.CreateTeam(c.Request.Context(), token, &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Team created successfully"})
}

func (h *TeamHandler) AddManager(c *gin.Context) {
	teamID := c.Param("teamID")

	if teamID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Team ID is required"})
		return
	}

	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
		return
	}

	var req *dto.AddManagerInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.AddManager(c.Request.Context(), token, uuid.MustParse(teamID), req.ManagerIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Manager added successfully"})
}

func (s *TeamHandler) AddMember(c *gin.Context) {
	teamID := c.Param("teamID")

	if teamID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Team ID is required"})
		return
	}

	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
		return
	}

	var req *dto.AddMemberInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.service.AddMember(c.Request.Context(), token, uuid.MustParse(teamID), req.MemberIDs); err != nil {
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

	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
		return
	}

	if err := s.service.RemoveManager(c.Request.Context(), token, uuid.MustParse(teamID), uuid.MustParse(managerID)); err != nil {
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

	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
		return
	}

	if err := s.service.RemoveMember(c.Request.Context(), token, uuid.MustParse(teamID), uuid.MustParse(memberID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Member removed successfully"})
}
