package models

type Team struct {
	Name    string
	Members []Member
}

type Member struct {
	UserID   string
	Username string
	IsActive bool
}
