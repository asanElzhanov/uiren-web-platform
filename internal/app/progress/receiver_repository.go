package progress

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type progressReceiverRepository struct {
	db *pgxpool.Pool
}

func NewProgressReceiverRepository(db *pgxpool.Pool) *progressReceiverRepository {
	return &progressReceiverRepository{
		db: db,
	}
}

func (r *progressReceiverRepository) getBadges(ctx context.Context, id string) ([]string, error) {
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

func (r *progressReceiverRepository) getXP(ctx context.Context, id string) (int, error) {
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

func (r *progressReceiverRepository) getAchievementsProgress(ctx context.Context, user_id string) ([]UserAchievement, error) {
	var (
		query = `
		SELECT 
    		ach.name, 
    		prg.achievement_level, 
    		lvl.description, 
    		prg.progress, 
    		lvl.threshold
		FROM 
    		users_achievements prg
		LEFT JOIN
    		achievements ach ON prg.achievement_id = ach.id
		LEFT JOIN
    		achievements_levels lvl ON lvl.achievement_id = prg.achievement_id AND lvl.level = prg.achievement_level
		WHERE
    		prg.user_id = $1;

		`
		result []UserAchievement
	)

	rows, err := r.db.Query(ctx, query, user_id)
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

func (r *progressReceiverRepository) getAchievementProgress(ctx context.Context, userID string, achID int) (UserAchievement, error) {
	var (
		query = `
		SELECT 
    		ach.name, 
    		prg.achievement_level, 
    		lvl.description, 
    		prg.progress, 
    		lvl.threshold
		FROM 
    		users_achievements prg
		LEFT JOIN
    		achievements ach ON prg.achievement_id = ach.id
		LEFT JOIN
    		achievements_levels lvl ON lvl.achievement_id = prg.achievement_id AND lvl.level = prg.achievement_level
		WHERE
    		prg.user_id = $1
			AND prg.achievement_id = $2;

		`
		achievement UserAchievement
	)

	if err := r.db.QueryRow(ctx, query, userID, achID).Scan(
		&achievement.AchievementName,
		&achievement.Level,
		&achievement.LevelDescription,
		&achievement.Progress,
		&achievement.Threshold,
	); err != nil {
		if err == pgx.ErrNoRows {
			return UserAchievement{}, ErrAchievementProgressNotFound
		}
		return UserAchievement{}, err
	}

	return achievement, nil
}

func (r *progressReceiverRepository) getXPLeaderboard(ctx context.Context, limit int) (XPLeaderboard, error) {
	var (
		query = `
		WITH leaderboard AS(
			SELECT 
				user_id as uid, 
				username, 
				xp, 
				row_number() OVER (ORDER BY xp DESC) as rank_number
			FROM users_progress up
			JOIN users u ON up.user_id = u.id
			WHERE u.deleted_at IS NULL
		)

		SELECT uid, username, xp, rank_number 
		FROM leaderboard
		ORDER BY rank_number
		LIMIT $1;
		`

		counter = 0
		result  XPLeaderboard
	)

	rows, err := r.db.Query(ctx, query, limit)
	if err != nil {
		return XPLeaderboard{}, err
	}

	for rows.Next() {
		var entry XPLeaderboardEntry
		if err := rows.Scan(
			&entry.UserID,
			&entry.Username,
			&entry.XP,
			&entry.Rank,
		); err != nil {
			return XPLeaderboard{}, err
		}

		counter++
		result.Leaders = append(result.Leaders, entry)
	}

	err = rows.Err()
	if err != nil {
		return XPLeaderboard{}, err
	}

	result.Total = counter
	return result, nil
}
