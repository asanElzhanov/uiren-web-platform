package admin

import (
	"strconv"
	"uiren/internal/app/achievements"
	"uiren/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

func (app *App) createAchievement(c *fiber.Ctx) error {
	var (
		ctx = c.Context()
		req CreateAchievementRequest
	)
	logger.Info("app.createAchievement handler")

	if err := c.BodyParser(&req); err != nil {
		logger.Error("app.createAchievement c.BodyParser: ", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": ErrBadRequest + ", name required"})
	}

	achievement, err := app.achievementService.CreateAchievement(ctx, req.Name)
	if err != nil {
		logger.Error("app.createAchievement achievementService.CreateAchievement: ", err)
		switch err {
		case achievements.ErrAchievementNameExists:
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": achievements.ErrAchievementNameExists.Error()})
		default:
			return fiberInternalServerError(c)
		}
	}

	return c.Status(fiber.StatusOK).JSON(achievement)
}

func (app *App) updateAchievement(c *fiber.Ctx) error {
	var (
		ctx = c.Context()
		req UpdateAchievementRequest
	)
	logger.Info("app.updateAchievement handler")

	if err := c.BodyParser(&req); err != nil {
		logger.Error("app.updateAchievement c.BodyParser: ", err)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"error": err.Error()})
	}

	if req.ID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": ErrBadRequest + ", id required"})
	}
	if req.NewName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": ErrBadRequest + ", name required"})
	}

	newName, err := app.achievementService.UpdateAchievement(ctx, achievements.UpdateAchievementDTO{
		ID:      req.ID,
		NewName: req.NewName,
	})
	if err != nil {
		logger.Error("app.updateAchievement achievementService.UpdateAchievement: ", err)
		switch err {
		case achievements.ErrAchievementNameExists:
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": achievements.ErrAchievementNameExists.Error()})
		case achievements.ErrAchievementNotFound:
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": achievements.ErrAchievementNotFound.Error()})
		default:
			return fiberInternalServerError(c)
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"new_name": newName})
}

func (app *App) deleteAchievement(c *fiber.Ctx) error {
	var (
		ctx = c.Context()
		req = c.Query("id")
	)

	if req == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": ErrBadRequest + ", id required"})
	}

	id, err := strconv.Atoi(req)
	if err != nil {
		logger.Error("app.deleteAchievement strconv.Atoi: ", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": ErrBadRequest + ", id invalid"})
	}

	if err := app.achievementService.DeleteAchievement(ctx, id); err != nil {
		logger.Error("app.deleteAchievement achievementService.DeleteAchievement: ", err)
		switch err {
		case achievements.ErrAchievementNotFound:
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": achievements.ErrAchievementNotFound.Error()})
		default:
			return fiberInternalServerError(c)
		}
	}

	return fiberOK(c)
}

func (app *App) getAchievement(c *fiber.Ctx) error {
	var (
		ctx = c.Context()
		req = c.Params("id")
	)

	if req == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": ErrBadRequest + ", id required"})
	}

	id, err := strconv.Atoi(req)
	if err != nil {
		logger.Error("app.deleteAchievement strconv.Atoi: ", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": ErrBadRequest + ", id invalid"})
	}

	achievement, err := app.achievementService.GetAchievement(ctx, id)
	if err != nil {
		logger.Error("app.deleteAchievement achievementService.GetAchievement: ", err)
		switch err {
		case achievements.ErrAchievementNotFound:
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": achievements.ErrAchievementNotFound.Error()})
		default:
			return fiberInternalServerError(c)
		}
	}

	return c.Status(fiber.StatusOK).JSON(achievement)
}

func (app *App) getAllAchievements(c *fiber.Ctx) error {
	var (
		ctx = c.Context()
	)

	resp, err := app.achievementService.GetAllAchievements(ctx)
	if err != nil {
		logger.Error("app.getAllAchievements achievementService.GetAllAchievements: ", err)
		return fiberInternalServerError(c)
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (app *App) addAchievementLevel(c *fiber.Ctx) error {
	var (
		ctx = c.Context()
		req AddAchievementLevelRequest
	)

	if err := c.BodyParser(&req); err != nil {
		logger.Error("app.addAchievementLevel c.BodyParser: ", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if req.AchID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": ErrBadRequest + ", id required"})
	}
	if req.Description == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": ErrBadRequest + ", description required"})
	}
	if req.Threshold == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": ErrBadRequest + ", threshold required"})
	}

	if err := app.achievementService.AddAchievementLevel(ctx, achievements.AddAchievementLevelDTO{
		AchID:       req.AchID,
		Description: req.Description,
		Threshold:   req.Threshold,
	}); err != nil {
		logger.Error("app.addAchievementLevel achievementService.AddAchievementLevel: ", err)
		switch err {
		case achievements.ErrAchievementNotFound:
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": achievements.ErrAchievementNotFound.Error()})
		case achievements.ErrInvalidThreshold:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": achievements.ErrInvalidThreshold.Error()})
		case achievements.ErrLowThreshold:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": achievements.ErrLowThreshold.Error()})
		case achievements.ErrLevelExists:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": achievements.ErrLevelExists.Error()})
		default:
			return fiberInternalServerError(c)
		}
	}

	return fiberOK(c)
}

func (app *App) deleteAchievementLevel(c *fiber.Ctx) error {
	var (
		ctx = c.Context()
		req DeleteAchievementLevelRequest
	)

	if err := c.BodyParser(&req); err != nil {
		logger.Error("app.deleteAchievementLevel c.BodyParser: ", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if req.AchID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": ErrBadRequest + ", id required"})
	}
	if req.Level == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": ErrBadRequest + ", level required"})
	}

	if err := app.achievementService.DeleteAchievementLevel(ctx, achievements.DeleteAchievementLevelDTO{
		AchID: req.AchID,
		Level: req.Level,
	}); err != nil {
		logger.Error("app.deleteAchievementLevel achievementService.DeleteAchievementLevel: ", err)
		switch err {
		case achievements.ErrAchievementLevelNotFound:
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": achievements.ErrAchievementLevelNotFound.Error()})
		default:
			return fiberInternalServerError(c)
		}
	}

	return fiberOK(c)
}
