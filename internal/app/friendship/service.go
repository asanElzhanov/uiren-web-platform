package friendship

import (
	"context"
	"errors"
	"uiren/internal/app/users"
	"uiren/pkg/logger"
)

//go:generate mockgen -source service.go -destination service_mock.go -package friendship

type friendshipRepository interface {
	createFriendshipStatus(ctx context.Context, req FriendshipRequestDTO) (Friendship, error)
	changeFriendshipStatus(ctx context.Context, req FriendshipRequestDTO) (Friendship, error)
	getFriendList(ctx context.Context, username string) (FriendList, error)
	getRequestList(ctx context.Context, username string) (FriendList, error)
	getFriendshipRecipient(ctx context.Context, username1, username2 string) (string, error)
}

type userService interface {
	CheckUserExists(ctx context.Context, username string) error
}

type FriendshipService struct {
	friendshipRepository friendshipRepository
	userService          userService
}

func NewFriendshipService(friendshipRepository friendshipRepository, userService userService) *FriendshipService {
	return &FriendshipService{
		friendshipRepository: friendshipRepository,
		userService:          userService,
	}
}

func (s *FriendshipService) SendFriendRequest(ctx context.Context, friendshipRequest FriendshipRequestDTO) (Friendship, error) {
	logger.Info("FriendshipService.SendFriendRequest new request")

	if friendshipRequest.RequesterUsername == friendshipRequest.RecipientUsername {
		return Friendship{}, ErrSameUser
	}

	if err := s.checkUsersExistance(ctx, friendshipRequest); err != nil {
		return Friendship{}, err
	}

	friendshipRequest.Recipient = friendshipRequest.RecipientUsername
	friendshipRequest.normalize()
	friendship, err := s.friendshipRepository.createFriendshipStatus(ctx, friendshipRequest)
	if err != nil {
		logger.Error("FriendshipService.SendFriendRequest friendshipRepository.setFriendshipStatus: ", err)
		return Friendship{}, err
	}
	return friendship, nil
}

func (s *FriendshipService) HandleFriendRequest(ctx context.Context, friendshipRequest FriendshipRequestDTO) (Friendship, error) {
	logger.Info("FriendshipService.AcceptFriendRequest new request")

	if !isValidStatus(friendshipRequest.Status) {
		return Friendship{}, ErrInvalidStatus
	}

	if friendshipRequest.RequesterUsername == friendshipRequest.RecipientUsername {
		return Friendship{}, ErrSameUser
	}

	if err := s.checkUsersExistance(ctx, friendshipRequest); err != nil {
		return Friendship{}, err
	}

	recipient, err := s.friendshipRepository.getFriendshipRecipient(ctx, friendshipRequest.RequesterUsername, friendshipRequest.RecipientUsername)
	if err != nil {
		logger.Error("FriendshipService.AcceptFriendRequest friendshipRepository.getFriendshipRecipient: ", err)
		return Friendship{}, err
	}

	// requester must be the one who got request
	if recipient != friendshipRequest.RequesterUsername {
		return Friendship{}, ErrNotRecipient
	}
	friendshipRequest.Recipient = recipient
	friendshipRequest.normalize()

	friendship, err := s.friendshipRepository.changeFriendshipStatus(ctx, friendshipRequest)
	if err != nil {
		logger.Error("FriendshipService.AcceptFriendRequest friendshipRepository.setFriendshipStatus: ", err)
		return Friendship{}, err
	}
	return friendship, nil
}

func (s *FriendshipService) checkUsersExistance(ctx context.Context, req FriendshipRequestDTO) error {
	err := s.userService.CheckUserExists(ctx, req.RequesterUsername)
	if err != nil {
		logger.Error("FriendshipService.SendFriendRequest GetUserByUsername error: ", err)
		if errors.Is(err, users.ErrUserNotFound) {
			return ErrRequesterNotFound
		}
		return err
	}

	err = s.userService.CheckUserExists(ctx, req.RecipientUsername)
	if err != nil {
		logger.Error("FriendshipService.SendFriendRequest GetUserByUsername error: ", err)
		if errors.Is(err, users.ErrUserNotFound) {
			return ErrRecipientNotFound
		}
		return err
	}

	return nil
}

func (s *FriendshipService) GetFriendList(ctx context.Context, username string) (FriendList, error) {
	logger.Info("FriendshipService.GetFriendList new request")

	friendList, err := s.friendshipRepository.getFriendList(ctx, username)
	if err != nil {
		logger.Error("FriendshipService.GetFriendList friendshipRepository.getFriendList: ", err)
		return FriendList{}, err
	}
	return friendList, nil
}

func (s *FriendshipService) GetRequestList(ctx context.Context, username string) (FriendList, error) {
	logger.Info("FriendshipService.GetRequestList new request")

	friendList, err := s.friendshipRepository.getRequestList(ctx, username)
	if err != nil {
		logger.Error("FriendshipService.GetRequestList friendshipRepository.getRequestList: ", err)
		return FriendList{}, err
	}
	return friendList, nil
}
