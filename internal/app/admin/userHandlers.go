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
	//todo поменять просто на getBYid
	var (
		ctx         = c.Context()
		indentifier = c.Query("identifier")
	)
	logger.Info("app.getUser handler")

	if indentifier == "" {
		logger.Error("app.getUser identifier required")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": ErrBadRequest + ", identifier required"})
	}
	userDTO, err := app.userService.GetUserForLogin(ctx, indentifier)
	if err != nil {
		logger.Error("app.getUser app.userService.GetUserForLogin: ", err)
		if err == users.ErrUserNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": ErrUserNotFound})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": ErrInternalServerError})
	}

	resp := User{
		ID:        userDTO.ID,
		Username:  userDTO.Username,
		Firstname: userDTO.Firstname,
		Lastname:  userDTO.Lastname,
		Email:     userDTO.Email,
		Phone:     userDTO.Phone,
	}
	return c.JSON(resp)
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
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": ErrInternalServerError})
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
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": ErrInternalServerError})
	}
}
