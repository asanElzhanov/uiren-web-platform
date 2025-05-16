package admin

import (
	"strconv"
	"uiren/internal/app/data"
	"uiren/internal/app/exercises"
	"uiren/internal/app/lessons"
	"uiren/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

func (app *App) mainPageModules(c *fiber.Ctx) error {
	var (
		ctx = c.Context()
	)

	modules, err := app.dataService.GetPublicModules(ctx)
	if err != nil {
		logger.Error("app.getModulesForMainPage dataService.GetModules: ", err)
		return fiberInternalServerError(c)
	}

	return c.Status(fiber.StatusOK).JSON(modules)
}

func (app *App) getLessonToPass(c *fiber.Ctx) error {
	var (
		ctx = c.Context()
		req = c.Query("code")
	)

	if req == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": ErrBadRequest + ", code required"})
	}

	lesson, err := app.dataService.GetPublicLesson(ctx, req)
	if err != nil {
		logger.Error("app.getLessonToPass dataService.GetPublicLesson: ", err)
		switch err {
		case lessons.ErrNotFound:
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": lessons.ErrNotFound.Error()})
		default:
			return fiberInternalServerError(c)
		}
	}

	return c.Status(fiber.StatusOK).JSON(lesson)
}

func (app *App) getExerciseToPass(c *fiber.Ctx) error {
	var (
		ctx = c.Context()
		req = c.Query("code")
	)
	if req == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": ErrBadRequest + ", code required"})
	}

	exercise, err := app.dataService.GetPublicExercise(ctx, req)
	if err != nil {
		logger.Error("app.getExerciseToPass dataService.GetPublicExercise: ", err)
		switch err {
		case exercises.ErrNotFound:
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": exercises.ErrNotFound.Error()})
		default:
			return fiberInternalServerError(c)
		}
	}

	return c.Status(fiber.StatusOK).JSON(exercise)
}

func (app *App) getUserInfo(c *fiber.Ctx) error {
	var (
		ctx               = c.Context()
		username          = c.Query("username")
		withProgressQuery = c.Query("withProgress")
		err               error
		userInfo          data.UserInfo
	)
	if username == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": ErrBadRequest + ", username required"})
	}
	if withProgressQuery == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": ErrBadRequest + ", withProgress required"})
	}

	withProgress, err := strconv.ParseBool(withProgressQuery)
	if err != nil {
		logger.Error("app.getUserByUsername ParseBool: ", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": ErrBadRequest + ", withProgress invalid format"})
	}

	if withProgress {
		userInfo, err = app.dataService.GetUserWithProgress(ctx, username)
	} else {
		userInfo, err = app.dataService.GetUserWithoutProgress(ctx, username)
	}
	if err != nil {
		logger.Error("app.getUserByUsername dataService.GetUserWithProgress: ", err)
		return getUserError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(userInfo)
}

func (app *App) getXPLeaderboard(c *fiber.Ctx) error {
	var (
		ctx = c.Context()
	)

	leaderboard, err := app.dataService.GetXPLeaderboard(ctx)
	if err != nil {
		logger.Error("app.getXPLeaderboard dataService.GetXPLeaderboard: ", err)
		return fiberInternalServerError(c)
	}

	return c.Status(fiber.StatusOK).JSON(leaderboard.Board)
}
