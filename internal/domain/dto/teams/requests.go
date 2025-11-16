package teams

type CreateRequest struct {
	TeamName string         `json:"team_name"`
	Members  []MemberCreate `json:"members"`
}

type MemberCreate struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

type GetParams struct {
	TeamName string `form:"team_name"`
}

type BulkDeactivateRequest struct {
    TeamName string `json:"team_name"`
}