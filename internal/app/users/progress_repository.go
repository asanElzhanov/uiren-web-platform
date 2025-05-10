package users

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userProgressRepository struct {
	db *pgxpool.Pool
}

func NewUserProgressRepository(db *pgxpool.Pool) *userProgressRepository {
	return &userProgressRepository{
		db: db,
	}
}

func (r *userProgressRepository) getBadges(ctx context.Context, id string) ([]string, error) {
	var (
		query = `
		SELECT 
			badge
		FROM
			users_badges
		WHERE
			user_id = $1;
		`
		result []string
	)

	rows, err := r.db.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var badge string
		if err := rows.Scan(&badge); err != nil {
			return nil, err
		}
		result = append(result, badge)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *userProgressRepository) getXP(ctx context.Context, id string) (int, error) {
	var (
		query = `
		SELECT 
			xp
		FROM 
			users_progress
		WHERE
			user_id = $1;
		`

		result int
	)

	if err := r.db.QueryRow(ctx, query, id).Scan(&result); err != nil {
		if err == pgx.ErrNoRows {
			return 0, nil
		}
		return -1, err
	}

	return result, nil
}

func (r *userProgressRepository) getAchievements(ctx context.Context, id string) ([]UserAchievement, error) {
	var (
		query = `
		SELECT 
    		ach.name, 
    		prg.achievement_level, 
    		lvl.description, 
    		prg.progress, 
    		lvl.threshold
		FROM 
    		user_achievements prg
		LEFT JOIN
    		achievements ach ON prg.achievement_id = ach.id
		LEFT JOIN
    		achievements_levels lvl ON lvl.achievement_id = prg.achievement_id AND lvl.level = prg.achievement_level
		WHERE
    		prg.user_id = $1;

		`
		result []UserAchievement
	)

	rows, err := r.db.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var achievement UserAchievement
		if err := rows.Scan(
			&achievement.AchievementName,
			&achievement.Level,
			&achievement.LevelDescription,
			&achievement.Progress,
			&achievement.Threshold,
		); err != nil {
			return nil, err
		}
		result = append(result, achievement)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return result, nil
}
