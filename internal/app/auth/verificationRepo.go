package auth

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type verificationRepository struct {
	db *pgxpool.Pool
}

func NewVerificationRepository(db *pgxpool.Pool) *verificationRepository {
	return &verificationRepository{
		db: db,
	}
}

func (r *verificationRepository) createVerificationCode(ctx context.Context, req CreateVerificationCodeRequest) error {
	var (
		query = `
		INSERT INTO
			users_verification_codes(username, email, verification_code, expires_at)
		VALUES
			($1, $2, $3, $4);
		`
	)

	_, err := r.db.Exec(ctx, query, req.Username, req.Email, req.Code, time.Now().Add(req.Duration))
	return err
}

func (r *verificationRepository) getVerificationCode(ctx context.Context, username string) (Verification, error) {
	var (
		query = `
		SELECT 
			username, email, verification_code, expires_at
		FROM
			users_verification_codes
		WHERE
			username = $1;
		`
		response Verification
	)

	row := r.db.QueryRow(ctx, query, username)

	if err := row.Scan(&response.Username, &response.Email, &response.Code, &response.ExpiresAt); err != nil {
		return Verification{}, err
	}

	return response, nil
}
