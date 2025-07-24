package dto

type CreateTeamRequest struct {
	TeamName string `json:"teamName"`
	Managers []struct {
		ManagerID   string `json:"managerId"`
		ManagerName string `json:"managerName"`
	} `json:"managers"`
	Members []struct {
		MemberID   string `json:"memberId"`
		MemberName string `json:"memberName"`
	} `json:"members"`
}

type UserIDRequest struct {
	UserID string `json:"user_id" binding:"required"`
}
