package postgres

import (
	"context"
	"database/sql"
	"log"

	"github.com/IlyaAGL/avito_autumn_2025/internal/models"
)

type postgresTeamRepo struct {
	db *sql.DB
}

func NewPostgresTeamRepository(db *sql.DB) *postgresTeamRepo {
	return &postgresTeamRepo{db: db}
}

func (repo *postgresTeamRepo) CreateTeam(ctx context.Context, team *models.Team) error {
	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err := tx.Rollback(); err != nil {
			log.Println("Failed to rollback")
		}
	}()

	_, err = tx.ExecContext(ctx,
		"INSERT INTO teams (team_name) VALUES ($1) ON CONFLICT (team_name) DO NOTHING",
		team.Name,
	)
	if err != nil {
		return err
	}

	for _, member := range team.Members {
		_, err = tx.ExecContext(ctx,
			`INSERT INTO users (user_id, username, team_name, is_active) 
             VALUES ($1, $2, $3, $4)
             ON CONFLICT (user_id) 
             DO UPDATE SET username = $2, team_name = $3, is_active = $4, updated_at = NOW()`,
			member.UserID, member.Username, team.Name, member.IsActive,
		)
		if err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (repo *postgresTeamRepo) GetTeam(ctx context.Context, teamName string) (*models.Team, error) {
	var team models.Team
	team.Name = teamName

	rows, err := repo.db.QueryContext(ctx,
		"SELECT user_id, username, is_active FROM users WHERE team_name = $1 ORDER BY user_id",
		teamName,
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
		var member models.Member
		if err := rows.Scan(&member.UserID, &member.Username, &member.IsActive); err != nil {
			return nil, err
		}
		team.Members = append(team.Members, member)
	}

	if len(team.Members) == 0 {
		return nil, sql.ErrNoRows
	}

	return &team, nil
}

func (repo *postgresTeamRepo) TeamExists(ctx context.Context, teamName string) (bool, error) {
	var exists bool

	err := repo.db.QueryRowContext(ctx,
		"SELECT EXISTS(SELECT 1 FROM teams WHERE team_name = $1)",
		teamName,
	).Scan(&exists)
	
	return exists, err
}

func (repo *postgresTeamRepo) BulkDeactivateUsers(ctx context.Context, teamName string) error {
    tx, err := repo.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    
	defer func() {
		if err := tx.Rollback(); err != nil {
			log.Println("Failed to rollback")
		}
	}()

    _, err = tx.ExecContext(ctx, `
        DELETE FROM pull_request_reviewers prr
        USING pull_requests pr, users u
        WHERE prr.pull_request_id = pr.pull_request_id
        AND prr.user_id = u.user_id
        AND u.team_name = $1
        AND pr.status = 'OPEN'
    `, teamName)
    if err != nil {
        return err
    }

    _, err = tx.ExecContext(ctx,
        "UPDATE users SET is_active = false, updated_at = NOW() WHERE team_name = $1",
        teamName,
    )
    if err != nil {
        return err
    }

    if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}