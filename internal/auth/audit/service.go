package audit

import (
	"context"
	"net"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AuditService struct {
	db *pgxpool.Pool
}

type LoginAttempt struct {
	Email         string
	IPAddress     net.IP
	UserAgent     string
	Success       bool
	FailureReason string
	Provider      string
	CreatedAt     time.Time
}

func NewAuditService(db *pgxpool.Pool) *AuditService {
	return &AuditService{db: db}
}

// LogLoginAttempt records a login attempt
func (s *AuditService) LogLoginAttempt(ctx context.Context, attempt LoginAttempt) error {
	query := `
		INSERT INTO login_attempts (email, ip_address, user_agent, success, failure_reason, provider, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := s.db.Exec(ctx, query,
		attempt.Email,
		attempt.IPAddress,
		attempt.UserAgent,
		attempt.Success,
		attempt.FailureReason,
		attempt.Provider,
		time.Now(),
	)

	return err
}

// UpdateLastLogin updates the user's last login timestamp
func (s *AuditService) UpdateLastLogin(ctx context.Context, userID int64) error {
	query := `UPDATE users SET last_login_at = NOW() WHERE id = $1`
	_, err := s.db.Exec(ctx, query, userID)
	return err
}

// IncrementLoginAttempts increments the failed login attempts counter
func (s *AuditService) IncrementLoginAttempts(ctx context.Context, email string) error {
	query := `
		UPDATE users 
		SET login_attempts = login_attempts + 1 
		WHERE email = $1
	`
	_, err := s.db.Exec(ctx, query, email)
	return err
}

// ResetLoginAttempts resets the counter after a successful login
func (s *AuditService) ResetLoginAttempts(ctx context.Context, email string) error {
	query := `
		UPDATE users 
		SET login_attempts = 0, account_locked_until = NULL 
		WHERE email = $1
	`
	_, err := s.db.Exec(ctx, query, email)
	return err
}

// LockAccount locks the account for a specified duration
func (s *AuditService) LockAccount(ctx context.Context, email string, lockDuration time.Duration) error {
	query := `
		UPDATE users 
		SET account_locked_until = $1 
		WHERE email = $2
	`
	lockUntil := time.Now().Add(lockDuration)
	_, err := s.db.Exec(ctx, query, lockUntil, email)
	return err
}

// IsAccountLocked checks if the account is locked
func (s *AuditService) IsAccountLocked(ctx context.Context, email string) (bool, error) {
	query := `
		SELECT account_locked_until 
		FROM users 
		WHERE email = $1
	`

	var lockUntil *time.Time
	err := s.db.QueryRow(ctx, query, email).Scan(&lockUntil)
	if err != nil {
		return false, err
	}

	if lockUntil != nil && lockUntil.After(time.Now()) {
		return true, nil
	}

	return false, nil
}

// GetLoginAttempts gets the number of login attempts
func (s *AuditService) GetLoginAttempts(ctx context.Context, email string) (int, error) {
	query := `SELECT login_attempts FROM users WHERE email = $1`

	var attempts int
	err := s.db.QueryRow(ctx, query, email).Scan(&attempts)
	if err != nil {
		return 0, err
	}

	return attempts, nil
}
