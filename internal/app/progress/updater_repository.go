package progress

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type progressUpdaterRepository struct {
	db *pgxpool.Pool
}

func NewProgressUpdaterRepository(db *pgxpool.Pool) *progressUpdaterRepository {
	return &progressUpdaterRepository{
		db: db,
	}
}

func (r *progressUpdaterRepository) beginTransaction(ctx context.Context) (transaction, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (r *progressUpdaterRepository) addBadges(ctx context.Context, tx transaction, req AddBadgesRequest) error {
	var (
		query = `
		INSERT INTO
			users_badges(user_id, badge)
		VALUES

		`
	)

	for _, badge := range req.Badges {
		query += fmt.Sprintf(`('%s', '%s'),`, req.UserID, badge)
	}
	query = query[:len(query)-1] + ";"

	_, err := tx.Exec(ctx, query)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrUserHasBadge
		}
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return ErrBadgeNotExists
		}
		return err
	}
	return nil
}

func (r *progressUpdaterRepository) addXP(ctx context.Context, tx transaction, req AddXPRequest) error {
	var (
		query = `
		INSERT INTO 
			users_progress(user_id, xp)
		VALUES
			($1, $2)
		ON CONFLICT ON CONSTRAINT unique_user_id
		DO UPDATE 
			SET xp = users_progress.xp + EXCLUDED.xp
		`
	)

	_, err := tx.Exec(ctx, query, req.UserID, req.XP)
	if err != nil {
		return err
	}
	return nil
}

func (r *progressUpdaterRepository) updateAchievementProgress(ctx context.Context, tx transaction, req UpdateAchievementProgressRequest) error {
	var (
		query = `
		INSERT INTO
			users_achievements(user_id, achievement_id, achievement_level, progress)
		VALUES
			($1, $2, $3, $4)
		ON CONFLICT ON CONSTRAINT unique_ach DO UPDATE
			SET progress = users_achievements.progress + EXCLUDED.progress, achievement_level = EXCLUDED.achievement_level`
	)

	_, err := tx.Exec(ctx, query, req.UserID, req.Progress.AchievementID, req.Progress.NewLevel, req.Progress.EarnedProgress)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23514" {
			return ErrNegativeProgress
		}
		return err
	}

	return nil
}

func (r *progressUpdaterRepository) insertBadge(ctx context.Context, req InsertBadgeRequest) error {
	var (
		query = `
		INSERT INTO
			badges(badge, description)
		VALUES
			($1, $2);`
	)

	_, err := r.db.Exec(ctx, query, req.Badge, req.Description)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrBadgeAlreadyExists
		}
		return err
	}

	return nil
}
