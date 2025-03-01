package admin

import (
	"errors"
	"fmt"
	"uiren/internal/app/auth"
	"uiren/internal/app/users"
	"uiren/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

func (app *App) signIn(c *fiber.Ctx) error {
	var (
		ctx    = c.Context()
		params SignInParams
	)
	logger.Info("app.signIn handler")

	if err := c.BodyParser(&params); err != nil {
		logger.Error("app.signIn error: ", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": ErrBadRequest})
	}

	if params.Identificator == "" || params.Password == "" {
		logger.Error("app.signIn error: ", ErrBadRequest)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": ErrBadRequest})
	}

	token, err := app.authService.SignIn(ctx, auth.LoginParams{
		Identificator: params.Identificator,
		Password:      params.Password,
	})

	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			logger.Error("app.signIn error: ", ErrInvalidCredentials)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": ErrInvalidCredentials})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"token": token})
}

func (app *App) register(c *fiber.Ctx) error {
	var (
		ctx = c.Context()
		req CreateUserReq
	)
	logger.Info("app.createUser handler")

	if err := c.BodyParser(&req); err != nil {
		logger.Error("app.getUser BodyParser: ", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": ErrBadRequest})
	}

	userID, err := app.authService.Register(ctx, auth.RegisterParams{
		DTO: users.CreateUserDTO{
			Username: req.Username,
			Email:    req.Email,
			Password: req.Password,
		},
	})
	if err != nil {
		return returnCreateUserError(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"id": userID})
}

func (app *App) verification(c *fiber.Ctx) error {
	logger.Info("app.register handler")
	var (
		ctx      = c.Context()
		username = c.Params("username")
		code     = c.Params("code")
	)

	if username == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": ErrBadRequest + ", username required"})
	}

	if err := app.authService.VerifyUser(ctx, username, code); err != nil {
		logger.Error("app.verification error: ", err)
		switch err {
		case auth.ErrVerificationInvalid:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": auth.ErrVerificationInvalid.Error()})
		case auth.ErrVerificationExpired:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": auth.ErrVerificationExpired.Error()})
		default:
			return fiberInternalServerError(c)
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": fmt.Sprintf("%s verified", username)})
}
