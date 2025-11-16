package pullrequests

type CreateResponse struct {
	PR PullRequestResponse `json:"pr"`
}

type MergeResponse struct {
	PR PullRequestResponse `json:"pr"`
}

type ReassignResponse struct {
	PR         PullRequestResponse `json:"pr"`
	ReplacedBy string              `json:"replaced_by"`
}

type PullRequestResponse struct {
	PullRequestID     string   `json:"pull_request_id"`
	PullRequestName   string   `json:"pull_request_name"`
	AuthorID          string   `json:"author_id"`
	Status            string   `json:"status"`
	AssignedReviewers []string `json:"assigned_reviewers"`
	MergedAt          *string  `json:"mergedAt,omitempty"`
}
