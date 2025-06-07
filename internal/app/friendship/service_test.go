package friendship

import (
	"context"
	"errors"
	"testing"
	"uiren/internal/app/users"
	"uiren/pkg/logger"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func init() {
	logger.InitLogger("info")
}

func Test_friendshipService_SendFriendRequest_success(t *testing.T) {
	t.Parallel()
	var (
		ctx         = context.TODO()
		ctrl        = gomock.NewController(t)
		userService = NewMockuserService(ctrl)
		repo        = NewMockfriendshipRepository(ctrl)
		srv         = NewFriendshipService(repo, userService)
		req         = FriendshipRequestDTO{
			RequesterUsername: "user1",
			RecipientUsername: "user2",
		}
	)
	userService.EXPECT().CheckUserExists(ctx, req.RequesterUsername).Return(nil)
	userService.EXPECT().CheckUserExists(ctx, req.RecipientUsername).Return(nil)

	newReq := FriendshipRequestDTO{
		RequesterUsername: req.RequesterUsername,
		RecipientUsername: req.RecipientUsername,
	}
	newReq.Recipient = req.RecipientUsername
	newReq.normalize()

	res := Friendship{
		Username1: newReq.RequesterUsername,
		Username2: newReq.RecipientUsername,
		Status:    statusPending,
		Recipient: newReq.Recipient,
	}
	repo.EXPECT().createFriendshipStatus(ctx, newReq).Return(res, nil)

	friendship, err := srv.SendFriendRequest(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, friendship, res)
}

func Test_friendshipService_SendFriendRequest_fail(t *testing.T) {
	t.Parallel()
	var (
		ctx         = context.TODO()
		ctrl        = gomock.NewController(t)
		userService = NewMockuserService(ctrl)
		repo        = NewMockfriendshipRepository(ctrl)
		srv         = NewFriendshipService(repo, userService)
	)
	//userService.EXPECT().CheckUserExists(ctx, req.RequesterUsername).Return(nil)
	//userService.EXPECT().CheckUserExists(ctx, req.RecipientUsername).Return(nil)

	t.Run("same user", func(t *testing.T) {
		req := FriendshipRequestDTO{
			RequesterUsername: "user",
			RecipientUsername: "user",
		}
		_, err := srv.SendFriendRequest(ctx, req)
		assert.Error(t, err)
		assert.Equal(t, err, ErrSameUser)
	})

	t.Run("requester does not exist", func(t *testing.T) {
		req := FriendshipRequestDTO{
			RequesterUsername: "user1",
			RecipientUsername: "user2",
		}
		userService.EXPECT().CheckUserExists(ctx, req.RequesterUsername).Return(users.ErrUserNotFound)
		userService.EXPECT().CheckUserExists(ctx, req.RecipientUsername).Return(nil).Times(0)

		_, err := srv.SendFriendRequest(ctx, req)
		assert.Error(t, err)
		assert.Equal(t, err, ErrRequesterNotFound)
	})

	t.Run("recipient does not exist", func(t *testing.T) {
		req := FriendshipRequestDTO{
			RequesterUsername: "user1",
			RecipientUsername: "user2",
		}
		userService.EXPECT().CheckUserExists(ctx, req.RequesterUsername).Return(nil)
		userService.EXPECT().CheckUserExists(ctx, req.RecipientUsername).Return(users.ErrUserNotFound)

		_, err := srv.SendFriendRequest(ctx, req)
		assert.Error(t, err)
		assert.Equal(t, err, ErrRecipientNotFound)
	})

	t.Run("create friendship status repo error", func(t *testing.T) {
		repoError := errors.New("database exploded")
		req := FriendshipRequestDTO{
			RequesterUsername: "user1",
			RecipientUsername: "user2",
		}
		userService.EXPECT().CheckUserExists(ctx, req.RequesterUsername).Return(nil)
		userService.EXPECT().CheckUserExists(ctx, req.RecipientUsername).Return(nil)

		newReq := FriendshipRequestDTO{
			RequesterUsername: req.RequesterUsername,
			RecipientUsername: req.RecipientUsername,
		}
		newReq.Recipient = req.RecipientUsername
		newReq.normalize()

		repo.EXPECT().createFriendshipStatus(ctx, newReq).Return(Friendship{}, repoError)

		_, err := srv.SendFriendRequest(ctx, req)
		assert.Error(t, err)
		assert.Equal(t, err, repoError)
	})
}

func Test_friendshipService_HandleFriendRequest_success(t *testing.T) {
	t.Parallel()
	var (
		ctx         = context.TODO()
		ctrl        = gomock.NewController(t)
		userService = NewMockuserService(ctrl)
		repo        = NewMockfriendshipRepository(ctrl)
		srv         = NewFriendshipService(repo, userService)
		req         = FriendshipRequestDTO{
			RequesterUsername: "user1",
			RecipientUsername: "user2",
			Status:            statusAccepted,
		}
	)
	userService.EXPECT().CheckUserExists(ctx, req.RequesterUsername).Return(nil)
	userService.EXPECT().CheckUserExists(ctx, req.RecipientUsername).Return(nil)
	repo.EXPECT().getFriendshipRecipient(ctx, req.RequesterUsername, req.RecipientUsername).Return(req.RequesterUsername, nil)

	newReq := FriendshipRequestDTO{
		RequesterUsername: req.RequesterUsername,
		RecipientUsername: req.RecipientUsername,
		Status:            statusAccepted,
	}
	newReq.Recipient = req.RequesterUsername
	newReq.normalize()

	res := Friendship{
		Username1: newReq.RequesterUsername,
		Username2: newReq.RecipientUsername,
		Status:    statusAccepted,
		Recipient: newReq.Recipient,
	}
	repo.EXPECT().changeFriendshipStatus(ctx, newReq).Return(res, nil)

	friendship, err := srv.HandleFriendRequest(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, friendship, res)
}

func Test_friendshipService_HandleFriendRequest_fail(t *testing.T) {
	t.Parallel()
	var (
		ctx         = context.TODO()
		ctrl        = gomock.NewController(t)
		repo        = NewMockfriendshipRepository(ctrl)
		userService = NewMockuserService(ctrl)
		srv         = NewFriendshipService(repo, userService)
	)
	defer ctrl.Finish()
	//userService.EXPECT().CheckUserExists(ctx, req.RequesterUsername).Return(nil)
	//userService.EXPECT().CheckUserExists(ctx, req.RecipientUsername).Return(nil)

	t.Run("same user", func(t *testing.T) {
		req := FriendshipRequestDTO{
			RequesterUsername: "user1",
			RecipientUsername: "user1",
			Status:            statusAccepted,
		}

		_, err := srv.HandleFriendRequest(ctx, req)
		assert.Error(t, err)
		assert.Equal(t, err, ErrSameUser)
	})
	t.Run("invalid status", func(t *testing.T) {
		req := FriendshipRequestDTO{
			RequesterUsername: "user1",
			RecipientUsername: "user1",
			Status:            "some status",
		}

		_, err := srv.HandleFriendRequest(ctx, req)
		assert.Error(t, err)
		assert.Equal(t, err, ErrInvalidStatus)
	})

	t.Run("requester does not exist", func(t *testing.T) {
		req := FriendshipRequestDTO{
			RequesterUsername: "user1",
			RecipientUsername: "user2",
			Status:            statusDeclined,
		}
		userService.EXPECT().CheckUserExists(ctx, req.RequesterUsername).Return(users.ErrUserNotFound)

		_, err := srv.HandleFriendRequest(ctx, req)
		assert.Error(t, err)
		assert.Equal(t, err, ErrRequesterNotFound)
	})

	t.Run("recipient does not exist", func(t *testing.T) {
		req := FriendshipRequestDTO{
			RequesterUsername: "user1",
			RecipientUsername: "user2",
			Status:            statusDeclined,
		}
		userService.EXPECT().CheckUserExists(ctx, req.RequesterUsername).Return(nil)
		userService.EXPECT().CheckUserExists(ctx, req.RecipientUsername).Return(users.ErrUserNotFound)

		_, err := srv.HandleFriendRequest(ctx, req)
		assert.Error(t, err)
		assert.Equal(t, err, ErrRecipientNotFound)
	})

	t.Run("CheckUserExists fail#1", func(t *testing.T) {
		someErr := errors.New("some error")
		req := FriendshipRequestDTO{
			RequesterUsername: "user1",
			RecipientUsername: "user2",
			Status:            statusDeclined,
		}
		userService.EXPECT().CheckUserExists(ctx, req.RequesterUsername).Return(someErr)
		userService.EXPECT().CheckUserExists(ctx, req.RecipientUsername).Return(nil).Times(0)

		_, err := srv.HandleFriendRequest(ctx, req)
		assert.Error(t, err)
		assert.Equal(t, err, someErr)
	})

	t.Run("CheckUserExists fail#2", func(t *testing.T) {
		someErr := errors.New("some error")
		req := FriendshipRequestDTO{
			RequesterUsername: "user1",
			RecipientUsername: "user2",
			Status:            statusDeclined,
		}
		userService.EXPECT().CheckUserExists(ctx, req.RequesterUsername).Return(nil)
		userService.EXPECT().CheckUserExists(ctx, req.RecipientUsername).Return(someErr)

		_, err := srv.HandleFriendRequest(ctx, req)
		assert.Error(t, err)
		assert.Equal(t, err, someErr)
	})

	t.Run("getFriendshipRecipient error", func(t *testing.T) {
		someErr := errors.New("some error")
		req := FriendshipRequestDTO{
			RequesterUsername: "user1",
			RecipientUsername: "user2",
			Status:            statusDeclined,
		}
		userService.EXPECT().CheckUserExists(ctx, req.RequesterUsername).Return(nil)
		userService.EXPECT().CheckUserExists(ctx, req.RecipientUsername).Return(nil)

		repo.EXPECT().getFriendshipRecipient(ctx, req.RequesterUsername, req.RecipientUsername).Return("", someErr)

		_, err := srv.HandleFriendRequest(ctx, req)
		assert.Error(t, err)
		assert.Equal(t, err, someErr)
	})

	t.Run("requester is not recipient", func(t *testing.T) {
		req := FriendshipRequestDTO{
			RequesterUsername: "user1",
			RecipientUsername: "user2",
			Status:            statusDeclined,
		}
		userService.EXPECT().CheckUserExists(ctx, req.RequesterUsername).Return(nil)
		userService.EXPECT().CheckUserExists(ctx, req.RecipientUsername).Return(nil)

		repo.EXPECT().getFriendshipRecipient(ctx, req.RequesterUsername, req.RecipientUsername).Return(req.RecipientUsername, nil)

		_, err := srv.HandleFriendRequest(ctx, req)
		assert.Error(t, err)
		assert.Equal(t, err, ErrNotRecipient)
	})

	t.Run("change frienship status error", func(t *testing.T) {
		someErr := errors.New("some error")
		req := FriendshipRequestDTO{
			RequesterUsername: "user1",
			RecipientUsername: "user2",
			Status:            statusDeclined,
		}
		userService.EXPECT().CheckUserExists(ctx, req.RequesterUsername).Return(nil)
		userService.EXPECT().CheckUserExists(ctx, req.RecipientUsername).Return(nil)

		repo.EXPECT().getFriendshipRecipient(ctx, req.RequesterUsername, req.RecipientUsername).Return(req.RequesterUsername, nil)

		newReq := FriendshipRequestDTO{
			RequesterUsername: req.RequesterUsername,
			RecipientUsername: req.RecipientUsername,
			Status:            statusDeclined,
		}
		newReq.Recipient = req.RequesterUsername
		newReq.normalize()

		repo.EXPECT().changeFriendshipStatus(ctx, newReq).Return(Friendship{}, someErr)

		_, err := srv.HandleFriendRequest(ctx, req)
		assert.Error(t, err)
		assert.Equal(t, err, someErr)
	})
}

func Test_FriendshipService_GetFriendList_success(t *testing.T) {
	t.Parallel()
	var (
		ctx   = context.TODO()
		ctrl  = gomock.NewController(t)
		repo  = NewMockfriendshipRepository(ctrl)
		srv   = NewFriendshipService(repo, nil)
		user  = "test_user"
		fList = FriendList{
			Friends: []FriendListEntity{
				{
					Username: "requester 1",
				},
				{
					Username: "requester 2",
				},
			},
			Total: 2,
		}
	)

	repo.EXPECT().getFriendList(ctx, user).Return(fList, nil)

	result, err := srv.GetFriendList(ctx, user)
	assert.NoError(t, err)
	assert.Equal(t, fList, result)
}

func Test_FriendshipService_GetFriendList_repoError(t *testing.T) {
	t.Parallel()
	var (
		ctx       = context.TODO()
		ctrl      = gomock.NewController(t)
		repo      = NewMockfriendshipRepository(ctrl)
		srv       = NewFriendshipService(repo, nil)
		user      = "test_user"
		repoError = errors.New("database error")
	)

	repo.EXPECT().getFriendList(ctx, user).Return(FriendList{}, repoError)

	result, err := srv.GetFriendList(ctx, user)
	assert.Error(t, err)
	assert.Equal(t, repoError, err)
	assert.Equal(t, FriendList{}, result)
}

func Test_FriendshipService_GetRequestList_success(t *testing.T) {
	t.Parallel()
	var (
		ctx   = context.TODO()
		ctrl  = gomock.NewController(t)
		repo  = NewMockfriendshipRepository(ctrl)
		srv   = NewFriendshipService(repo, nil)
		user  = "test_user"
		fList = FriendList{
			Friends: []FriendListEntity{
				{
					Username: "requester 1",
				},
				{
					Username: "requester 2",
				},
			},
			Total: 2,
		}
	)

	repo.EXPECT().getRequestList(ctx, user).Return(fList, nil)

	result, err := srv.GetRequestList(ctx, user)
	assert.NoError(t, err)
	assert.Equal(t, fList, result)
}

func Test_FriendshipService_GetRequestList_repoError(t *testing.T) {
	t.Parallel()
	var (
		ctx       = context.TODO()
		ctrl      = gomock.NewController(t)
		repo      = NewMockfriendshipRepository(ctrl)
		srv       = NewFriendshipService(repo, nil)
		user      = "test_user"
		repoError = errors.New("database error")
	)

	repo.EXPECT().getRequestList(ctx, user).Return(FriendList{}, repoError)

	result, err := srv.GetRequestList(ctx, user)
	assert.Error(t, err)
	assert.Equal(t, repoError, err)
	assert.Equal(t, FriendList{}, result)
}
