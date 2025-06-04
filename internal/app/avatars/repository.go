package avatars

import (
	"context"
	"io"
	"os"
	"path/filepath"
)

type repository struct {
	folderPath string
}

func NewAvatarRepository(folderPath string) *repository {
	folderPath = "./storage" + folderPath + "/"
	return &repository{
		folderPath: folderPath,
	}
}

func (r *repository) uploadAvatar(ctx context.Context, req UploadAvatarRequest) error {
	var (
		filename = generateAvatarFilename(req.UserId)
		fullPath = r.folderPath + filepath.Join(filename)
	)

	file, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, req.File)
	return err
}
