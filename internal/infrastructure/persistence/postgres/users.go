package postgres

import (
	"context"
	"database/sql"
	"log"

	"github.com/IlyaAGL/avito_autumn_2025/internal/models"
	"github.com/lib/pq"
)

type postgresUserRepo struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *postgresUserRepo {
	return &postgresUserRepo{db: db}
}

func (repo *postgresUserRepo) CreateOrUpdateUser(ctx context.Context, user *models.User) error {
	_, err := repo.db.ExecContext(ctx,
		`INSERT INTO users (user_id, username, team_name, is_active) 
         VALUES ($1, $2, $3, $4)
         ON CONFLICT (user_id) 
         DO UPDATE SET username = $2, team_name = $3, is_active = $4, updated_at = NOW()`,
		user.UserID, user.Username, user.TeamName, user.IsActive,
	)

	return err
}

func (repo *postgresUserRepo) GetUser(ctx context.Context, userID string) (*models.User, error) {
	var user models.User
	err := repo.db.QueryRowContext(ctx,
		"SELECT user_id, username, team_name, is_active FROM users WHERE user_id = $1",
		userID,
	).Scan(&user.UserID, &user.Username, &user.TeamName, &user.IsActive)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *postgresUserRepo) SetUserActive(ctx context.Context, userID string, isActive bool) (*models.User, error) {
	_, err := repo.db.ExecContext(ctx,
		"UPDATE users SET is_active = $1, updated_at = NOW() WHERE user_id = $2",
		isActive, userID,
	)
	if err != nil {
		return nil, err
	}

	return repo.GetUser(ctx, userID)
}

func (repo *postgresUserRepo) GetActiveTeamMembers(ctx context.Context, teamName string, excludeUserIDs []string) ([]models.User, error) {
	var users []models.User

	var query string
	
	var args []any

	if len(excludeUserIDs) > 0 {
		query = `SELECT user_id, username, team_name, is_active 
                 FROM users 
                 WHERE team_name = $1 AND is_active = true AND user_id != ALL($2)
                 ORDER BY user_id`
		args = []any{teamName, pq.Array(excludeUserIDs)}
	} else {
		query = `SELECT user_id, username, team_name, is_active 
                 FROM users 
                 WHERE team_name = $1 AND is_active = true
                 ORDER BY user_id`
		args = []any{teamName}
	}

	rows, err := repo.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			log.Println("Failed to close the rows")
		}
	}()

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.UserID, &user.Username, &user.TeamName, &user.IsActive); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (repo *postgresUserRepo) GetUserReviewPRs(ctx context.Context, userID string) ([]models.PullRequestShort, error) {
	var prs []models.PullRequestShort

	rows, err := repo.db.QueryContext(ctx,
		`SELECT pr.pull_request_id, pr.pull_request_name, pr.author_id, pr.status
         FROM pull_requests pr
         JOIN pull_request_reviewers prr ON pr.pull_request_id = prr.pull_request_id
         WHERE prr.user_id = $1
         ORDER BY pr.created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			log.Println("Failed to close the rows")
		}
	}()

	for rows.Next() {
		var pr models.PullRequestShort
		if err := rows.Scan(&pr.ID, &pr.Name, &pr.AuthorID, &pr.Status); err != nil {
			return nil, err
		}
		
		prs = append(prs, pr)
	}

	return prs, nil
}
