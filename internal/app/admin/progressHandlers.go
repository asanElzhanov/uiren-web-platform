package admin

import (
	"encoding/json"
	"uiren/internal/app/achievements"
	"uiren/internal/app/progress"

	"github.com/gofiber/fiber/v2"
)

func (app *App) updateProgress(c *fiber.Ctx) error {
	var (
		ctx     = c.Context()
		req     progress.UpdateUserProgressRequest
		rawData map[string]interface{}
	)

	if err := json.Unmarshal(c.Body(), &rawData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": ErrBadRequest + ", invalid request"})
	}
	if err := validateAchievementProgress(rawData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": ErrBadRequest + ", invalid request"})
	}

	if req.UserID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": ErrBadRequest + ", userID must be provided"})
	}
	if req.XP < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": ErrBadRequest + ", XP must be greater than 0"})
	}

	if err := app.progressService.UpdateUserProgress(ctx, req); err != nil {
		switch err {
		case progress.ErrBadgeNotProvided:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": progress.ErrBadgeNotProvided.Error()})
		case progress.ErrUserHasBadge:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": progress.ErrUserHasBadge.Error()})
		case progress.ErrBadgeNotExists:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": progress.ErrBadgeNotExists.Error()})
		case progress.ErrNegativeProgress:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": progress.ErrNegativeProgress.Error()})
		case achievements.ErrAchievementNotFound:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": achievements.ErrAchievementNotFound.Error()})
		default:
			return fiberInternalServerError(c)
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "user progress updated"})
}

func (app *App) registerBadge(c *fiber.Ctx) error {
	var (
		ctx = c.Context()
		req progress.InsertBadgeRequest
	)

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": ErrBadRequest + ", invalid request"})
	}
	if req.Badge == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": ErrBadRequest + ", badge required"})
	}
	if req.Description == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": ErrBadRequest + ", description must be provided"})
	}

	if err := app.progressService.RegisterNewBadge(ctx, req); err != nil {
		switch err {
		case progress.ErrUserHasBadge:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": progress.ErrUserHasBadge.Error()})
		default:
			return fiberInternalServerError(c)
		}
	}

	return fiberOK(c)
}
