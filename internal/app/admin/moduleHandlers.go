package admin

import (
	"uiren/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

func (app *App) getModuleInfo(c *fiber.Ctx) error {
	var (
		ctx = c.Context()
		req = c.Params("code")
	)
	logger.Info("app.getModuleInfo handler")

	resp, err := app.modulesService.GetModule(ctx, req)
	if err != nil {
		logger.Error("app.modulesService.GetModule: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}
