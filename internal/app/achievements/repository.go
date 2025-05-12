package achievements

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type achievementRepository struct {
	db *pgxpool.Pool
}

func NewAchievementRepository(db *pgxpool.Pool) *achievementRepository {
	return &achievementRepository{
		db: db,
	}
}

func (r *achievementRepository) beginTransaction(ctx context.Context) (transaction, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (r *achievementRepository) getAchievement(ctx context.Context, id int) (achievement, error) {
	var (
		query = `
		SELECT
			id,
			name,
			created_at,
			updated_at
		FROM
			achievements
		WHERE
			id = $1
			AND deleted_at IS NULL;
		`
		response achievement
	)

	row := r.db.QueryRow(ctx, query, id)

	if err := row.Scan(
		&response.id,
		&response.name,
		&response.createdAt,
		&response.updatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return achievement{}, ErrAchievementNotFound
		}
		return achievement{}, err
	}

	return response, nil
}

func (r *achievementRepository) createAchievement(ctx context.Context, name string) (achievement, error) {
	var (
		query = `
		INSERT INTO 
			achievements(name) 
		VALUES 
			($1) 
		RETURNING 
			id, 
			name, 
			created_at, 
			updated_at; 
		`
		response achievement
	)

	row := r.db.QueryRow(ctx, query, name)
	if err := row.Scan(
		&response.id,
		&response.name,
		&response.createdAt,
		&response.updatedAt,
	); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return achievement{}, ErrAchievementNameExists
		}
		return achievement{}, err
	}

	return response, nil
}

func (r *achievementRepository) updateAchievement(ctx context.Context, dto UpdateAchievementDTO) (string, error) {
	var (
		query = `
			UPDATE 
				achievements 
			SET 
				name = $1, 
				updated_at = NOW() 
			WHERE id = $2
			AND deleted_at IS NULL 
			RETURNING name;
		`
		response string
	)

	row := r.db.QueryRow(ctx, query, dto.NewName, dto.ID)
	if err := row.Scan(&response); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return "", ErrAchievementNameExists
		}

		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrAchievementNotFound
		}
		return "", err
	}

	return response, nil
}

func (r *achievementRepository) deleteAchievement(ctx context.Context, id int) error {
	var (
		query = `
		UPDATE 
			achievements
		SET
			deleted_at = NOW()
		WHERE
			id = $1 AND deleted_at IS NULL;
		`
	)

	commandTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return ErrAchievementNotFound
	}

	return nil
}

func (r *achievementRepository) getLevelsByAchievementID(ctx context.Context, achID int) ([]AchievementLevel, error) {
	var (
		query = `
		SELECT
			ach.id,
			ach.name, 
			l.level, 
			l.description, 
			l.threshold, 
			l.created_at, 
			l.updated_at
		FROM 
			achievements_levels l
		INNER JOIN 
			achievements ach ON (l.achievement_id = ach.id)
		WHERE 
			ach.id = $1
			AND ach.deleted_at IS NULL
			AND l.deleted_at IS NULL;
		`
		response []AchievementLevel
	)

	rows, err := r.db.Query(ctx, query, achID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var level AchievementLevel
		if err := rows.Scan(
			&level.achID,
			&level.achName,
			&level.level,
			&level.description,
			&level.threshold,
			&level.createdAt,
			&level.updatedAt,
		); err != nil {
			return nil, err
		}

		response = append(response, level)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return response, nil
}

func (r *achievementRepository) getLastLevelAndTreshold(ctx context.Context, achID int) (LevelData, error) {
	var (
		query = `
		SELECT  
			COALESCE(l.level, 0),
			COALESCE(l.threshold, 0)
		FROM 
			achievements_levels l
		RIGHT JOIN 
			achievements ach ON (l.achievement_id = ach.id)
		WHERE 
			ach.id = $1
			AND ach.deleted_at IS NULL
			AND l.deleted_at IS NULL
		ORDER BY level DESC
		LIMIT 1;
		`
		response LevelData
	)

	if err := r.db.QueryRow(ctx, query, achID).Scan(&response.Level, &response.Threshold); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return LevelData{}, ErrAchievementNotFound
		}
		return LevelData{}, err
	}

	return response, nil
}

func (r *achievementRepository) addLevel(ctx context.Context, dto AddAchievementLevelDTO) error {
	var (
		query = `
		INSERT INTO 
			achievements_levels(achievement_id, level, description, threshold)
		VALUES
			($1, $2, $3 ,$4);
		`
	)

	_, err := r.db.Exec(ctx, query, dto.AchID, dto.Level, dto.Description, dto.Threshold)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrLevelExists
		}
	}
	return nil
}

func (r *achievementRepository) deleteLevel(ctx context.Context, tx transaction, dto DeleteAchievementLevelDTO) error {
	var (
		query = `
		DELETE FROM achievements_levels
		WHERE achievement_id = $1 AND level = $2;
		`
	)

	commandTag, err := tx.Exec(ctx, query, dto.AchID, dto.Level)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return ErrAchievementLevelNotFound
	}

	return nil
}

func (r *achievementRepository) decrementUpperLevels(ctx context.Context, tx transaction, dto DeleteAchievementLevelDTO) error {
	var (
		query = `
		UPDATE achievements_levels
		SET level = level - 1
		WHERE achievement_id = $1 AND level > $2;
		`
	)

	_, err := tx.Exec(ctx, query, dto.AchID, dto.Level)
	return err
}

func (r *achievementRepository) getLevel(ctx context.Context, achID, level int) (AchievementLevel, error) {
	var (
		query = `
		SELECT
			ach.id,
			ach.name, 
			l.level, 
			l.description, 
			l.threshold, 
			l.created_at, 
			l.updated_at
		FROM
			achievements_levels l
		INNER JOIN
			achievements ach ON (l.achievement_id = ach.id)
		WHERE
			ach.id = $1
			AND l.level = $2
			AND ach.deleted_at IS NULL;
		`
		result AchievementLevel
	)

	if err := r.db.QueryRow(ctx, query, achID, level).Scan(
		&result.achID,
		&result.achName,
		&result.level,
		&result.description,
		&result.threshold,
		&result.createdAt,
		&result.updatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return AchievementLevel{}, ErrAchievementLevelNotFound
		}
		return AchievementLevel{}, err
	}
	return result, nil
}
