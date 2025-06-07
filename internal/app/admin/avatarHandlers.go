package admin

import (
	"uiren/internal/app/avatars"
	"uiren/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

func (app *App) uploadAvatar(c *fiber.Ctx) error {
	var (
		ctx = c.Context()
	)

	file, err := c.FormFile("avatar")
	if err != nil {
		logger.Error("app.uploadAvatar c.FormFile: ", err)
		return fiberFormFileError(c, err)
	}

	idVal := c.Locals("id")
	id, ok := idVal.(string)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": ErrBadRequest + ", incorrect id"})
	}

	reader, err := file.Open()
	if err != nil {
		logger.Error("app.uploadAvatar file.Open: ", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": ErrBadRequest})
	}

	if err := app.avatarService.UploadAvatar(ctx, avatars.UploadAvatarRequest{
		UserId: id,
		File:   reader,
	}); err != nil {
		logger.Error("app.uploadAvatar avatarService.UploadAvatar: ", err)
		return err
	}

	return fiberOK(c)
}
