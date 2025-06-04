package avatars

import (
	"context"
	"uiren/pkg/logger"
)

type avatarRepository interface {
	uploadAvatar(ctx context.Context, req UploadAvatarRequest) error
}

type AvatarService struct {
	repo avatarRepository
}

func NewAvatarService(avatarRepository avatarRepository) *AvatarService {
	return &AvatarService{
		repo: avatarRepository,
	}
}

func (s *AvatarService) UploadAvatar(ctx context.Context, req UploadAvatarRequest) error {
	logger.Info("AvatarService.UploadAvatar new request")

	if err := s.repo.uploadAvatar(ctx, req); err != nil {
		logger.Error("AvatarService.UploadAvatar repo.uploadAvatar: ", err)
		return err
	}

	return nil
}
