package models

import "time"

type PullRequest struct {
	ID                string
	Name              string
	AuthorID          string
	Status            string
	AssignedReviewers []string
	CreatedAt         time.Time
	MergedAt          *time.Time
}

type ReviewStats struct {
	UserID   string
	OpenPRs  int
	TotalReviews int
}

type PullRequestShort struct {
	ID       string
	Name     string
	AuthorID string
	Status   string
}

type PRStats struct {
    TotalPRs  int `json:"total_prs"`
    OpenPRs   int `json:"open_prs"`
    MergedPRs int `json:"merged_prs"`
}

type SimpleStats struct {
    TotalPRs      int           `json:"total_prs"`
    OpenPRs       int           `json:"open_prs"`
    MergedPRs     int           `json:"merged_prs"`
    TopReviewers  []ReviewStats `json:"top_reviewers"`
}
