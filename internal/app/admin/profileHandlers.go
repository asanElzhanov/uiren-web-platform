package admin

import (
	"uiren/internal/app/users"
	"uiren/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

func (app *App) updateProfile(c *fiber.Ctx) error {
	var (
		ctx = c.Context()
		req UpdateUserReq
	)

	userIDVal := c.Locals("id")
	userID, ok := userIDVal.(string)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid token claims"})
	}

	if err := c.BodyParser(&req); err != nil {
		logger.Error("app.updateUser BodyParser: ", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": ErrBadRequest})
	}

	if userID == "" {
		logger.Error("app.updateUser: id required")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": ErrBadRequest + ", id required"})
	}

	updatedUser, err := app.userService.UpdateUser(ctx, users.UpdateUserDTO{
		ID:          userID,
		Firstname:   req.Firstname,
		Lastname:    req.Lastname,
		Phone:       req.Phone,
		PhoneRegion: req.PhoneRegion,
	})
	if err != nil {
		logger.Error("app.updateUser UpdateUser: ", err)

		switch err {
		case users.ErrUserNotFound:
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": ErrUserNotFound})
		case users.ErrIncorrectPhone:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": users.ErrIncorrectPhone.Error()})
		default:
			return fiberInternalServerError(c)
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"id":          updatedUser.ID,
		"first_name":  updatedUser.Firstname,
		"last_name":   updatedUser.Lastname,
		"phone":       updatedUser.Phone,
		"update_time": updatedUser.UpdatedAt,
	})
}
