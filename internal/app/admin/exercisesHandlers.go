package admin

import (
	"encoding/json"
	"uiren/internal/app/exercises"
	"uiren/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

func (app *App) getAllExercises(c *fiber.Ctx) error {
	var (
		ctx = c.Context()
	)

	resp, err := app.exerciseService.GetAllExercises(ctx)
	if err != nil {
		logger.Error("app.getAllExercises exerciseService.GetAllExercises: ", err)
		return fiberInternalServerError(c)
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (app *App) getExercise(c *fiber.Ctx) error {
	var (
		ctx = c.Context()
		req = c.Params("code")
	)
	logger.Info("app.getExercise handler")

	exercise, err := app.exerciseService.GetExercise(ctx, req)
	if err != nil {
		logger.Error("app.getExercise exerciseService.GetExercise: ", err)
		switch err {
		case exercises.ErrNotFound:
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": exercises.ErrNotFound.Error()})
		default:
			return fiberInternalServerError(c)
		}
	}

	return c.Status(fiber.StatusOK).JSON(exercise)
}

func (app *App) createExercise(c *fiber.Ctx) error {
	var (
		ctx     = c.Context()
		req     exercises.CreateExerciseDTO
		rawData map[string]interface{}
	)
	logger.Info("app.createExercise handler")

	if err := json.Unmarshal(c.Body(), &rawData); err != nil {
		logger.Error("app.createExercise json.Unmarshal: ", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := validatePairs(rawData); err != nil {
		logger.Error("app.createExercise validatePairs: ", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := c.BodyParser(&req); err != nil {
		logger.Error("app.createExercise c.BodyParser: ", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := validateCode(req.Code); err != nil {
		logger.Error("app.createExercise validateCode: ", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	oid, err := app.exerciseService.CreateExercise(ctx, req)
	if err != nil {
		logger.Error("app.createExercise exerciseService.CreateExercise: ", err)
		switch err {
		case exercises.ErrCodeAlreadyExists:
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": exercises.ErrCodeAlreadyExists.Error()})
		case exercises.ErrIncorrectType:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": exercises.ErrIncorrectType.Error()})
		case exercises.ErrOptionsRequired:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": exercises.ErrOptionsRequired.Error()})
		case exercises.ErrCorrectAnswerRequired:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": exercises.ErrCorrectAnswerRequired.Error()})
		case exercises.ErrPairsRequired:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": exercises.ErrPairsRequired.Error()})
		case exercises.ErrCorrectOrderRequired:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": exercises.ErrCorrectOrderRequired.Error()})
		default:
			return fiberInternalServerError(c)
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"id": oid.Hex()})
}

func (app *App) deleteExercise(c *fiber.Ctx) error {
	var (
		ctx = c.Context()
		req = c.Params("code")
	)
	logger.Info("app.deleteExercise handler")

	if err := app.exerciseService.DeleteExercise(ctx, req); err != nil {
		logger.Error("app.deleteExercise exerciseService.DeleteExercise: ", err)
		switch err {
		case exercises.ErrNotFound:
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": exercises.ErrNotFound.Error()})
		default:
			return fiberInternalServerError(c)
		}
	}

	return fiberOK(c)
}

func (app *App) updateExercise(c *fiber.Ctx) error {
	var (
		ctx     = c.Context()
		code    = c.Params("code")
		req     exercises.UpdateExerciseDTO
		rawData map[string]interface{}
	)
	logger.Info("app.updateExercise handler")

	if err := json.Unmarshal(c.Body(), &rawData); err != nil {
		logger.Error("app.updateExercise json.Unmarshal: ", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := validatePairs(rawData); err != nil {
		logger.Error("app.updateExercise validatePairs: ", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := c.BodyParser(&req); err != nil {
		logger.Error("app.updateExercise c.BodyParser: ", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := app.exerciseService.UpdateExercise(ctx, code, req); err != nil {
		logger.Error("app.updateService exerciseService.UpdateExercise: ", err)
		switch err {
		case exercises.ErrNotFound:
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": exercises.ErrNotFound.Error()})
		case exercises.ErrIncorrectType:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": exercises.ErrIncorrectType.Error()})
		case exercises.ErrNoFieldsToUpdate:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": exercises.ErrNoFieldsToUpdate.Error()})
		case exercises.ErrOptionsRequired:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": exercises.ErrOptionsRequired.Error()})
		case exercises.ErrCorrectAnswerRequired:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": exercises.ErrCorrectAnswerRequired.Error()})
		case exercises.ErrPairsRequired:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": exercises.ErrPairsRequired.Error()})
		case exercises.ErrCorrectOrderRequired:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": exercises.ErrCorrectOrderRequired.Error()})
		default:
			return fiberInternalServerError(c)
		}
	}

	return fiberOK(c)
}
