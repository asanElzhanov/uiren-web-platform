package friendship

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repository struct {
	db *pgxpool.Pool
}

func NewFriendshipRepository(db *pgxpool.Pool) *repository {
	return &repository{
		db: db,
	}
}

func (r *repository) createFriendshipStatus(ctx context.Context, req FriendshipRequestDTO) (Friendship, error) {
	var (
		query = `
		WITH inserted AS(
			INSERT INTO friendships (user1_username, user2_username, status, recipient)
			VALUES ($1, $2, 'pending', $3)
			ON CONFLICT (user1_username, user2_username)
			DO NOTHING
			RETURNING user1_username, user2_username, status, recipient
		)

		SELECT user1_username, user2_username, status, recipient from inserted
		UNION ALL
		SELECT user1_username, user2_username, status, recipient from friendships
		WHERE user1_username = $1 AND user2_username = $2;
		`
		response Friendship
	)

	if err := r.db.QueryRow(ctx, query, req.RequesterUsername, req.RecipientUsername, req.Recipient).Scan(
		&response.Username1,
		&response.Username2,
		&response.Status,
		&response.Recipient,
	); err != nil {
		return Friendship{}, err
	}

	return response, nil
}

func (r *repository) changeFriendshipStatus(ctx context.Context, req FriendshipRequestDTO) (Friendship, error) {
	var (
		query = `
		UPDATE friendships SET status = $1
		WHERE 
			user1_username = $2 
			AND user2_username = $3 
			AND status = 'pending'
			AND recipient = $4 
		RETURNING user1_username, user2_username, status, recipient; 
		`
		response Friendship
	)

	if err := r.db.QueryRow(ctx, query, req.Status, req.RequesterUsername, req.RecipientUsername, req.Recipient).Scan(
		&response.Username1,
		&response.Username2,
		&response.Status,
		&response.Recipient,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Friendship{}, ErrFriendshipNotFound
		}
	}

	return response, nil
}

func (r *repository) getFriendList(ctx context.Context, username string) (FriendList, error) {
	var (
		query = `
		SELECT u.username FROM friendships f
		JOIN users u ON f.user1_username = u.username OR f.user2_username = u.username
		WHERE (f.user1_username = $1 OR f.user2_username = $1) AND f.status = 'accepted' AND u.username != $1;
		`
		response FriendList
	)

	rows, err := r.db.Query(ctx, query, username)
	if err != nil {
		return FriendList{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var username string
		if err := rows.Scan(&username); err != nil {
			return FriendList{}, err
		}
		response.Usernames = append(response.Usernames, username)
		response.Total++
	}
	if err := rows.Err(); err != nil {
		return FriendList{}, err
	}

	return response, nil
}

func (r *repository) getRequestList(ctx context.Context, username string) (FriendList, error) {
	var (
		query = `
		SELECT u.username FROM friendships f
		JOIN users u ON f.user1_username = u.username OR f.user2_username = u.username
		WHERE 
		(f.user1_username = $1 OR f.user2_username = $1) 
		AND f.status = 'pending' 
		AND u.username != $1
		AND recipient = $1;
		`
		response FriendList
	)

	rows, err := r.db.Query(ctx, query, username)
	if err != nil {
		return FriendList{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var username string
		if err := rows.Scan(&username); err != nil {
			return FriendList{}, err
		}
		response.Usernames = append(response.Usernames, username)
		response.Total++
	}
	if err := rows.Err(); err != nil {
		return FriendList{}, err
	}

	return response, nil
}

func (r *repository) getFriendshipRecipient(ctx context.Context, username1, username2 string) (string, error) {
	var (
		query = `
		SELECT recipient FROM friendships
		WHERE 
			(user1_username = $1 AND user2_username = $2)
			OR (user1_username = $2 AND user2_username = $1);
		`
		recipient string
	)

	if err := r.db.QueryRow(ctx, query, username1, username2).Scan(&recipient); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrFriendshipNotFound
		}
		return "", err
	}

	return recipient, nil
}
