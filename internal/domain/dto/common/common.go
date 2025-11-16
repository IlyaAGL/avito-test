package common

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type StatsResponse struct {
	UserStats []ReviewStats `json:"user_stats"`
}

type ReviewStats struct {
	UserID   string `json:"user_id"`
	OpenPRs  int    `json:"open_prs"`
	TotalPRs int    `json:"total_prs"`
}
