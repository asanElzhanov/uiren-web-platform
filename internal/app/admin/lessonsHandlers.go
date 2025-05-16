package admin

import (
	"uiren/internal/app/exercises"
	"uiren/internal/app/lessons"
	"uiren/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

func (app *App) getAllLessons(c *fiber.Ctx) error {
	var (
		ctx = c.Context()
	)

	resp, err := app.lessonService.GetAllLessonsWithExercises(ctx)
	if err != nil {
		logger.Error("app.getLessons lessonService.GetAllLessonsWithExercises: ", err)
		return fiberInternalServerError(c)
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (app *App) getLesson(c *fiber.Ctx) error {
	var (
		ctx = c.Context()
		req = c.Params("code")
	)

	lesson, err := app.lessonService.GetLesson(ctx, req)
	if err != nil {
		switch err {
		case lessons.ErrNotFound:
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": lessons.ErrNotFound.Error()})
		default:
			return fiberInternalServerError(c)
		}
	}

	return c.Status(fiber.StatusOK).JSON(lesson)
}

func (app *App) createLesson(c *fiber.Ctx) error {
	var (
		ctx = c.Context()
		req lessons.CreateLessonDTO
	)
	logger.Info("app.createLesson handler")

	if err := c.BodyParser(&req); err != nil {
		logger.Error("app.CreateLesson c.BodyParser: ", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := validateCode(req.Code); err != nil {
		logger.Error("app.createLesson validateCode: ", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	id, err := app.lessonService.CreateLesson(ctx, req)
	if err != nil {
		logger.Error("app.createLesson lessonService.CreateLesson: ", err)
		switch err {
		case lessons.ErrCodeAlreadyExists:
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": lessons.ErrCodeAlreadyExists.Error()})
		default:
			return fiberInternalServerError(c)
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"id": id})
}

func (app *App) updateLesson(c *fiber.Ctx) error {
	var (
		ctx  = c.Context()
		req  lessons.UpdateLessonDTO
		code = c.Params("code")
	)
	logger.Info("app.updateLesson handler")

	if err := c.BodyParser(&req); err != nil {
		logger.Error("app.updateLesson c.BodyParser: ", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := app.lessonService.UpdateLesson(ctx, code, req); err != nil {
		logger.Error("app.updateLesson lessonService.UpdateLesson: ", err)
		switch err {
		case lessons.ErrNotFound:
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": lessons.ErrNotFound.Error()})
		case lessons.ErrNoFieldsToUpdate:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": lessons.ErrNoFieldsToUpdate.Error()})
		default:
			return fiberInternalServerError(c)
		}
	}

	return fiberOK(c)
}

func (app *App) deleteLesson(c *fiber.Ctx) error {
	var (
		ctx = c.Context()
		req = c.Params("code")
	)
	logger.Info("app.deleteLesson handler")

	if err := app.lessonService.DeleteLesson(ctx, req); err != nil {
		logger.Error("app.deleteLesson lessonService.DeleteLesson: ", err)
		switch err {
		case lessons.ErrNotFound:
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": lessons.ErrNotFound.Error()})
		default:
			return fiberInternalServerError(c)
		}
	}

	return fiberOK(c)
}

func (app *App) addExerciseToList(c *fiber.Ctx) error {
	var (
		ctx          = c.Context()
		code         = c.Params("code")
		exerciseCode = c.Params("exerciseCode")
	)

	if err := app.lessonService.AddExerciseToList(ctx, code, exerciseCode); err != nil {
		logger.Error("app.addExerciseToList lessonService.AddExerciseToList: ", err)
		switch err {
		case lessons.ErrNotFound:
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": lessons.ErrNotFound.Error()})
		case lessons.ErrExerciseAlreadyInSet:
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": lessons.ErrExerciseAlreadyInSet.Error()})
		case exercises.ErrNotFound:
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": exercises.ErrNotFound.Error()})
		default:
			return fiberInternalServerError(c)
		}
	}

	return fiberOK(c)
}

func (app *App) deleteExerciseFromList(c *fiber.Ctx) error {
	var (
		ctx          = c.Context()
		code         = c.Params("code")
		exerciseCode = c.Params("exerciseCode")
	)

	if err := app.lessonService.DeleteExerciseFromList(ctx, code, exerciseCode); err != nil {
		logger.Error("app.deleteExerciseFromList lessonService.DeleteExerciseFromList: ", err)
		switch err {
		case lessons.ErrNotFound:
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": lessons.ErrNotFound.Error()})
		case lessons.ErrExerciseNotInList:
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": lessons.ErrExerciseNotInList.Error()})
		default:
			return fiberInternalServerError(c)
		}
	}

	return fiberOK(c)
}
