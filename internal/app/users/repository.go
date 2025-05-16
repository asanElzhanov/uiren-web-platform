package users

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *userRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) createUser(ctx context.Context, params CreateUserDTO) (string, error) {
	var (
		query = `
		INSERT INTO 
			users(username, email, password)
		VALUES 
			($1, $2, $3)
		RETURNING 
			id;
		`
		id uuid.UUID
	)

	row := r.db.QueryRow(ctx, query, params.Username, params.Email, params.Password)

	if err := row.Scan(&id); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			switch pgErr.ConstraintName {
			case "users_username_key":
				return "", ErrUsernameExists
			case "users_email_key":
				return "", ErrEmailExists
			default:
				return "", err
			}
		}
		return "", err
	}

	return id.String(), nil
}

func (r *userRepository) getUserByID(ctx context.Context, id string) (UserDTO, error) {
	var (
		query = `
		SELECT 
			id,
			username,
			email,
			password,
			first_name,
			last_name,
			phone,
			is_active,
			is_admin,
			created_at,
			updated_at
		FROM 
			users
		WHERE 
			id = $1
			AND is_active = true 
			AND deleted_at IS NULL;
		`
	)

	row := r.db.QueryRow(ctx, query, id)

	var user user
	if err := row.Scan(
		&user.id,
		&user.username,
		&user.email,
		&user.password,
		&user.firstname,
		&user.lastname,
		&user.phone,
		&user.isActive,
		&user.isAdmin,
		&user.createdAt,
		&user.updatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return UserDTO{}, ErrUserNotFound
		}
		return UserDTO{}, err
	}

	return user.ToDTO(), nil
}

func (r *userRepository) getUserByUsername(ctx context.Context, username string) (UserDTO, error) {
	var (
		query = `
		SELECT 
			id,
			username,
			email,
			password,
			first_name,
			last_name,
			phone,
			is_active,
			is_admin,
			created_at,
			updated_at
		FROM 
			users
		WHERE 
			username = $1
			AND is_active = true 
			AND deleted_at IS NULL;
		`
	)

	row := r.db.QueryRow(ctx, query, username)

	var user user
	if err := row.Scan(
		&user.id,
		&user.username,
		&user.email,
		&user.password,
		&user.firstname,
		&user.lastname,
		&user.phone,
		&user.isActive,
		&user.isAdmin,
		&user.createdAt,
		&user.updatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return UserDTO{}, ErrUserNotFound
		}
		return UserDTO{}, err
	}

	return user.ToDTO(), nil
}

func (r *userRepository) getUserByEmail(ctx context.Context, email string) (UserDTO, error) {
	var (
		query = `
		SELECT 
			id,
			username,
			email,
			password,
			first_name,
			last_name,
			phone,
			is_active,
			is_admin,
			created_at,
			updated_at
		FROM 
			users
		WHERE 
			email = $1 
			AND is_active = true 
			AND deleted_at IS NULL;
		`
	)

	row := r.db.QueryRow(ctx, query, email)

	var user user
	if err := row.Scan(
		&user.id,
		&user.username,
		&user.email,
		&user.password,
		&user.firstname,
		&user.lastname,
		&user.phone,
		&user.isActive,
		&user.isAdmin,
		&user.createdAt,
		&user.updatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return UserDTO{}, ErrUserNotFound
		}
		return UserDTO{}, err
	}

	return user.ToDTO(), nil
}

func (r *userRepository) updateUser(ctx context.Context, dto UpdateUserDTO) (UserDTO, error) {
	var (
		query = `
		UPDATE 
			users 
		SET 
			phone = $1,
			first_name = $2,
			last_name = $3,
			updated_at = $4
		WHERE 
			id = $5
		RETURNING
			id,
			username,
			first_name,
			last_name,
			phone,
			updated_at;`

		updatedUser user
	)

	row := r.db.QueryRow(ctx, query, dto.Phone, dto.Firstname, dto.Lastname, time.Now(), dto.ID)

	if err := row.Scan(
		&updatedUser.id,
		&updatedUser.username,
		&updatedUser.firstname,
		&updatedUser.lastname,
		&updatedUser.phone,
		&updatedUser.updatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return UserDTO{}, ErrUserNotFound
		}
		return UserDTO{}, err
	}

	return updatedUser.ToDTO(), nil
}

func (r *userRepository) enableUser(ctx context.Context, username string) error {
	var (
		query = `
		UPDATE
			users
		SET
			is_active = true
		WHERE
			username = $1;
		`
	)

	_, err := r.db.Exec(ctx, query, username)
	return err
}

func (r *userRepository) checkUserExists(ctx context.Context, username string) error {
	var (
		query = `
		SELECT 
			true
		FROM 
			users
		WHERE 
			username = $1 
			AND is_active = true 
			AND deleted_at IS NULL;
		`
		exists bool
	)

	if err := r.db.QueryRow(ctx, query, username).Scan(&exists); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrUserNotFound
		}
		return err
	}

	return nil
}

func (r *userRepository) getAllUsers(ctx context.Context) ([]UserDTO, error) {
	var (
		query = `
		SELECT 
			id,
			username,
			email,
			password,
			first_name,
			last_name,
			phone,
			is_active,
			is_admin,
			created_at,
			updated_at
		FROM 
			users
		WHERE 
			is_active = true 
			AND deleted_at IS NULL;
		`
	)

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []UserDTO
	for rows.Next() {
		var user user
		if err := rows.Scan(
			&user.id,
			&user.username,
			&user.email,
			&user.password,
			&user.firstname,
			&user.lastname,
			&user.phone,
			&user.isActive,
			&user.isAdmin,
			&user.createdAt,
			&user.updatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, user.ToDTO())
	}

	return users, nil
}
