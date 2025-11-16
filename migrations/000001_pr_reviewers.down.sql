DROP INDEX IF EXISTS idx_pull_request_reviewers_user;
DROP INDEX IF EXISTS idx_pull_requests_author;
DROP INDEX IF EXISTS idx_pull_requests_status;
DROP INDEX IF EXISTS idx_users_team_active;

DROP TABLE IF EXISTS pull_request_reviewers;
DROP TABLE IF EXISTS pull_requests;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS teams;