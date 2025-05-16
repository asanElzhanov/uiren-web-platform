package admin

import (
	"uiren/internal/app/users"
	"uiren/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

func (app *App) createUser(c *fiber.Ctx) error {
	var (
		ctx = c.Context()
		req CreateUserReq
	)
	logger.Info("app.createUser handler")

	if err := c.BodyParser(&req); err != nil {
		logger.Error("app.getUser BodyParser: ", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": ErrBadRequest})
	}

	userID, err := app.userService.CreateUser(ctx, users.CreateUserDTO{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return returnCreateUserError(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"id": userID})
}

func (app *App) getUser(c *fiber.Ctx) error {
	var (
		ctx    = c.Context()
		userID = c.Params("id")
	)
	logger.Info("app.getUser handler")

	if userID == "" {
		logger.Error("app.getUser: id required")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": ErrBadRequest + ", id required"})
	}

	user, err := app.userService.GetUserByID(ctx, userID)
	if err != nil {
		logger.Error("app.getUser GetUser: ", err)
		return getUserError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

func (app *App) updateUser(c *fiber.Ctx) error {
	var (
		ctx    = c.Context()
		req    UpdateUserReq
		userID = c.Params("id")
	)
	logger.Info("app.updateUser handler")

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

func (app *App) getAllUsers(c *fiber.Ctx) error {
	var (
		ctx = c.Context()
	)
	logger.Info("app.getAllUsers handler")

	users, err := app.userService.GetAllUsers(ctx)
	if err != nil {
		logger.Error("app.getAllUsers GetAllUsers: ", err)
		return fiberInternalServerError(c)
	}

	return c.Status(fiber.StatusOK).JSON(users)
}

func returnCreateUserError(c *fiber.Ctx, err error) error {
	switch err {
	case users.ErrIncorrectEmail:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": users.ErrIncorrectEmail.Error()})
	case users.ErrIncorrectPhone:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": users.ErrIncorrectPhone.Error()})
	case users.ErrIncorrectPassword:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": users.ErrIncorrectPassword.Error()})
	case users.ErrUsernameExists:
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": users.ErrUsernameExists.Error()})
	case users.ErrEmailExists:
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": users.ErrEmailExists.Error()})
	default:
		return fiberInternalServerError(c)
	}
}

func getUserError(c *fiber.Ctx, err error) error {
	switch err {
	case users.ErrUserNotFound:
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": users.ErrUserNotFound.Error()})
	default:
		return fiberInternalServerError(c)
	}
}
