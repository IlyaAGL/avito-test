package teams

type CreateResponse struct {
	Team TeamResponse `json:"team"`
}

type TeamResponse struct {
	TeamName string           `json:"team_name"`
	Members  []MemberResponse `json:"members"`
}

type MemberResponse struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}
