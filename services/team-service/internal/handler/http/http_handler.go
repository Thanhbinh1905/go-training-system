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
	return &TeamHandler{service: service}
}

func (h *TeamHandler) CreateTeam(c *gin.Context) {
	var req dto.CreateTeamInput
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.CreateTeam(c.Request.Context(), &req); err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	respondMessage(c, http.StatusCreated, "Team created successfully")
}

func (h *TeamHandler) AddManager(c *gin.Context) {
	teamID, err := uuid.Parse(c.Param("teamID"))
	if err != nil {
		respondError(c, http.StatusBadRequest, "invalid team ID")
		return
	}

	var req dto.AddManagerInput
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.AddManager(c.Request.Context(), teamID, req.ManagerIDs); err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	respondMessage(c, http.StatusOK, "Manager added successfully")
}

func (s *TeamHandler) AddMember(c *gin.Context) {
	teamID := c.Param("teamID")
	if teamID == "" {
		respondError(c, http.StatusBadRequest, "Team ID is required")
		return
	}

	var req dto.AddMemberInput
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := s.service.AddMember(c.Request.Context(), uuid.MustParse(teamID), req.MemberIDs); err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	respondMessage(c, http.StatusOK, "Member added successfully")
}

func (s *TeamHandler) RemoveManager(c *gin.Context) {
	teamID := c.Param("teamID")
	managerID := c.Param("managerID")

	if teamID == "" || managerID == "" {
		respondError(c, http.StatusBadRequest, "Team ID and Manager ID are required")
		return
	}

	if err := s.service.RemoveManager(c.Request.Context(), uuid.MustParse(teamID), uuid.MustParse(managerID)); err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	respondMessage(c, http.StatusOK, "Manager removed successfully")
}

func (s *TeamHandler) RemoveMember(c *gin.Context) {
	teamID := c.Param("teamID")
	memberID := c.Param("memberID")

	if teamID == "" || memberID == "" {
		respondError(c, http.StatusBadRequest, "Team ID and Member ID are required")
		return
	}

	if err := s.service.RemoveMember(c.Request.Context(), uuid.MustParse(teamID), uuid.MustParse(memberID)); err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	respondMessage(c, http.StatusOK, "Member removed successfully")
}

func respondError(c *gin.Context, status int, msg string) {
	c.JSON(status, gin.H{"error": msg})
}

func respondMessage(c *gin.Context, status int, msg string) {
	c.JSON(status, gin.H{"message": msg})
}
