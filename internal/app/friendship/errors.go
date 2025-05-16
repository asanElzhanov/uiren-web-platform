package friendship

import "errors"

var (
	ErrRequesterNotFound  = errors.New("requester not found")
	ErrRecipientNotFound  = errors.New("recipient not found")
	ErrSameUser           = errors.New("same user")
	ErrFriendshipNotFound = errors.New("friendship does not exist")
	ErrInvalidStatus      = errors.New("invalid status")
	ErrNotRecipient       = errors.New("requester is not the friendship's recipient")
)
