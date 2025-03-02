package admin

import (
	"encoding/json"
	"uiren/internal/app/modules"
	"uiren/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

func (app *App) getModule(c *fiber.Ctx) error {
	var (
		ctx = c.Context()
		req = c.Params("code")
	)
	logger.Info("app.getModule handler")
	if req == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": ErrBadRequest})
	}

	resp, err := app.modulesService.GetModule(ctx, req)
	if err != nil {
		logger.Error("app.getModule modulesService.GetModule: ", err)
		switch err {
		case modules.ErrNotFound:
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": modules.ErrNotFound.Error()})
		}
		return fiberInternalServerError(c)
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (app *App) createModule(c *fiber.Ctx) error {
	var (
		ctx     = c.Context()
		req     modules.CreateModuleDTO
		rawData map[string]interface{}
	)
	logger.Info("app.createModule handler")

	if err := json.Unmarshal(c.Body(), &rawData); err != nil {
		logger.Error("app.updateModule json.Unmarshal: ", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := validateRewardsAndRequirements(rawData); err != nil {
		logger.Error("app.updateModule validateRewardsAndRequirements: ", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := c.BodyParser(&req); err != nil {
		logger.Error("app.CreateModule c.BodyParser: ", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := validateCode(req.Code); err != nil {
		logger.Error("app.CreateModule validateCode: ", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	id, err := app.modulesService.CreateModule(ctx, req)
	if err != nil {
		logger.Error("app.createModule modulesService.CreateModule: ", err)
		switch err {
		case modules.ErrInvalidCode:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": modules.ErrInvalidCode.Error()})
		case modules.ErrCodeAlreadyExists:
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": modules.ErrCodeAlreadyExists.Error()})
		default:
			return fiberInternalServerError(c)
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"id": id.Hex()})
}

func (app *App) deleteModule(c *fiber.Ctx) error {
	var (
		ctx = c.Context()
		req = c.Params("code")
	)
	logger.Info("app.deleteModule handler")

	if err := app.modulesService.DeleteModule(ctx, req); err != nil {
		logger.Error("app.deleteModule modulesService.DeleteModule: ", err)
		switch err {
		case modules.ErrNotFound:
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": modules.ErrNotFound.Error()})
		default:
			return fiberInternalServerError(c)
		}
	}

	return fiberOK(c)
}

func (app *App) updateModule(c *fiber.Ctx) error {
	var (
		ctx     = c.Context()
		code    = c.Params("code")
		rawData map[string]interface{}
		req     modules.UpdateModuleDTO
	)
	logger.Info("app.updateModule handler")

	if err := json.Unmarshal(c.Body(), &rawData); err != nil {
		logger.Error("app.updateModule json.Unmarshal: ", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := validateRewardsAndRequirements(rawData); err != nil {
		logger.Error("app.updateModule validateRewardsAndRequirements: ", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := c.BodyParser(&req); err != nil {
		logger.Error("app.updateModule c.BodyParser: ", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := app.modulesService.UpdateModule(ctx, code, req); err != nil {
		logger.Error("app.updateModule modulesService.UpdateModule: ", err)
		switch err {
		case modules.ErrNoFieldsToUpdate:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": modules.ErrNoFieldsToUpdate.Error()})
		case modules.ErrNotFound:
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": modules.ErrNotFound.Error()})
		default:
			return fiberInternalServerError(c)
		}
	}

	return fiberOK(c)
}

func (app *App) addLessonToList(c *fiber.Ctx) error {
	var (
		ctx        = c.Context()
		moduleCode = c.Params("code")
		lessonCode = c.Params("lessonCode")
	)

	if lessonCode == "" || moduleCode == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": ErrBadRequest})
	}

	if err := app.modulesService.AddLessonToList(ctx, moduleCode, lessonCode); err != nil {
		logger.Error("app.addLessonToList modulesService.AddLessonToList:", err)
		switch err {
		case modules.ErrNotFound:
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": modules.ErrNotFound.Error()})
		case modules.ErrLessonAlreadyInSet:
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": modules.ErrLessonAlreadyInSet.Error()})
		default:
			return fiberInternalServerError(c)
		}
	}

	return fiberOK(c)
}

func (app *App) deleteLessonFromList(c *fiber.Ctx) error {
	var (
		ctx        = c.Context()
		moduleCode = c.Params("code")
		lessonCode = c.Params("lessonCode")
	)

	if lessonCode == "" || moduleCode == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": ErrBadRequest})
	}

	if err := app.modulesService.DeleteLessonFromList(ctx, moduleCode, lessonCode); err != nil {
		logger.Error("app.deleteLessonFromList modulesService.DeleteLessonFromList:", err)
		switch err {
		case modules.ErrNotFound:
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": modules.ErrNotFound.Error()})
		case modules.ErrLessonNotInList:
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": modules.ErrLessonNotInList.Error()})
		default:
			return fiberInternalServerError(c)
		}
	}

	return fiberOK(c)
}
