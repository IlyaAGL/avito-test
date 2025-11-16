package users

type SetActiveRequest struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

type ReviewParams struct {
	UserID string `form:"user_id"`
}

type DeactivateRequest struct {
	TeamName string   `json:"team_name"`
	UserIDs  []string `json:"user_ids"`
}
