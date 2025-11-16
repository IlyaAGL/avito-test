package postgres

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/IlyaAGL/avito_autumn_2025/internal/models"
)

type postgresPRRepo struct {
	db *sql.DB
}

func NewPostgresPullRequestRepository(db *sql.DB) *postgresPRRepo {
	return &postgresPRRepo{db: db}
}

func (repo *postgresPRRepo) CreatePR(ctx context.Context, pr *models.PullRequest) error {
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
		`INSERT INTO pull_requests (pull_request_id, pull_request_name, author_id, status) 
         VALUES ($1, $2, $3, $4)`,
		pr.ID, pr.Name, pr.AuthorID, "OPEN",
	)
	if err != nil {
		return err
	}

	for _, reviewerID := range pr.AssignedReviewers {
		_, err = tx.ExecContext(ctx,
			"INSERT INTO pull_request_reviewers (pull_request_id, user_id) VALUES ($1, $2)",
			pr.ID, reviewerID,
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

func (repo *postgresPRRepo) GetPR(ctx context.Context, prID string) (*models.PullRequest, error) {
	var pr models.PullRequest

	err := repo.db.QueryRowContext(ctx,
		`SELECT pull_request_id, pull_request_name, author_id, status, created_at, merged_at
         FROM pull_requests WHERE pull_request_id = $1`,
		prID,
	).Scan(&pr.ID, &pr.Name, &pr.AuthorID, &pr.Status, &pr.CreatedAt, &pr.MergedAt)
	if err != nil {
		return nil, err
	}

	rows, err := repo.db.QueryContext(ctx,
		"SELECT user_id FROM pull_request_reviewers WHERE pull_request_id = $1 ORDER BY assigned_at",
		prID,
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
		var reviewerID string
		if err := rows.Scan(&reviewerID); err != nil {
			return nil, err
		}
		pr.AssignedReviewers = append(pr.AssignedReviewers, reviewerID)
	}

	return &pr, nil
}

func (repo *postgresPRRepo) PRExists(ctx context.Context, prID string) (bool, error) {
	var exists bool
	err := repo.db.QueryRowContext(ctx,
		"SELECT EXISTS(SELECT 1 FROM pull_requests WHERE pull_request_id = $1)",
		prID,
	).Scan(&exists)
	return exists, err
}

func (repo *postgresPRRepo) MergePR(ctx context.Context, prID string) error {
	mergedAt := time.Now()
	_, err := repo.db.ExecContext(ctx,
		"UPDATE pull_requests SET status = 'MERGED', merged_at = $1, updated_at = NOW() WHERE pull_request_id = $2 AND status != 'MERGED'",
		mergedAt, prID,
	)
	return err
}

func (repo *postgresPRRepo) UpdatePRReviewers(ctx context.Context, prID string, reviewerIDs []string) error {
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
		"DELETE FROM pull_request_reviewers WHERE pull_request_id = $1",
		prID,
	)
	if err != nil {
		return err
	}

	for _, reviewerID := range reviewerIDs {
		_, err = tx.ExecContext(ctx,
			"INSERT INTO pull_request_reviewers (pull_request_id, user_id) VALUES ($1, $2)",
			prID, reviewerID,
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

func (repo *postgresPRRepo) GetPRsByReviewer(ctx context.Context, userID string) ([]models.PullRequestShort, error) {
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

func (repo *postgresPRRepo) GetReviewStats(ctx context.Context) ([]models.ReviewStats, error) {
    query := `
        SELECT 
            u.user_id,
            u.username,
            COUNT(prr.pull_request_id) as total_reviews
        FROM users u
        LEFT JOIN pull_request_reviewers prr ON u.user_id = prr.user_id
        WHERE u.is_active = true
        GROUP BY u.user_id, u.username
        ORDER BY total_reviews DESC
    `

    rows, err := repo.db.QueryContext(ctx, query)
    if err != nil {
        return nil, err
    }
    
	defer func() {
		if err := rows.Close(); err != nil {
			log.Println("Failed to close the rows")
		}
	}()

    var stats []models.ReviewStats
    for rows.Next() {
        var stat models.ReviewStats
        var username string
        if err := rows.Scan(&stat.UserID, &username, &stat.TotalReviews); err != nil {
            return nil, err
        }
        stats = append(stats, stat)
    }

    return stats, nil
}