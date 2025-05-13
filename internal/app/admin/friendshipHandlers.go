package admin

import (
	"uiren/internal/app/friendship"
	"uiren/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

func (app *App) sendFriendRequest(c *fiber.Ctx) error {
	var (
		ctx = c.Context()
		req friendship.FriendshipRequestDTO
	)
	logger.Info("app.sendFriendRequest handler")

	if err := c.BodyParser(&req); err != nil {
		logger.Error("app.sendFriendRequest c.BodyParser: ", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	usernameVal := c.Locals("username")
	username, ok := usernameVal.(string)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": ErrBadRequest + ", incorrect token payload(missing username)"})
	}
	req.RequesterUsername = username

	if req.RequesterUsername == "" || req.RecipientUsername == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": ErrBadRequest + ", username required"})
	}

	newFriendship, err := app.friendshipService.SendFriendRequest(ctx, req)
	if err != nil {
		logger.Error("app.sendFriendRequest friendshipService.SendFriendRequest: ", err)
		switch err {
		case friendship.ErrSameUser:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": friendship.ErrSameUser.Error()})
		case friendship.ErrRequesterNotFound:
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": friendship.ErrRequesterNotFound.Error()})
		case friendship.ErrRecipientNotFound:
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": friendship.ErrRecipientNotFound.Error()})
		default:
			return fiberInternalServerError(c)
		}
	}

	return c.Status(fiber.StatusOK).JSON(newFriendship)
}

func (app *App) handleFriendRequest(c *fiber.Ctx) error {
	var (
		ctx = c.Context()
		req friendship.FriendshipRequestDTO
	)
	logger.Info("app.handleFriendRequest handler")

	if err := c.BodyParser(&req); err != nil {
		logger.Error("app.handleFriendRequest c.BodyParser: ", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	usernameVal := c.Locals("username")
	username, ok := usernameVal.(string)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": ErrBadRequest + ", incorrect token payload(missing username)"})
	}
	req.RequesterUsername = username

	if req.RequesterUsername == "" || req.RecipientUsername == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": ErrBadRequest + ", username required"})
	}
	if req.Status == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": ErrBadRequest + ", status required"})
	}

	usersFriendship, err := app.friendshipService.HandleFriendRequest(ctx, req)
	if err != nil {
		logger.Error("app.handleFriendRequest friendshipService.HandleFriendRequest: ", err)
		switch err {
		case friendship.ErrSameUser:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": friendship.ErrSameUser.Error()})
		case friendship.ErrRequesterNotFound:
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": friendship.ErrRequesterNotFound.Error()})
		case friendship.ErrRecipientNotFound:
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": friendship.ErrRecipientNotFound.Error()})
		case friendship.ErrFriendshipNotFound:
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": friendship.ErrFriendshipNotFound.Error()})
		case friendship.ErrInvalidStatus:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": friendship.ErrInvalidStatus.Error()})
		case friendship.ErrNotRecipient:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": friendship.ErrNotRecipient.Error()})
		default:
			return fiberInternalServerError(c)
		}
	}

	return c.Status(fiber.StatusOK).JSON(usersFriendship)
}

func (app *App) getFriendList(c *fiber.Ctx) error {
	var (
		ctx = c.Context()
		req = c.Query("username")
	)

	if req == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": ErrBadRequest + ", username required"})
	}

	friendList, err := app.friendshipService.GetFriendList(ctx, req)
	if err != nil {
		logger.Error("app.getFriendList friendshipService.GetFriendList: ", err)
		return fiberInternalServerError(c)
	}

	return c.Status(fiber.StatusOK).JSON(friendList)
}

func (app *App) getRequestList(c *fiber.Ctx) error {
	var (
		ctx = c.Context()
	)

	usernameVal := c.Locals("username")
	username, ok := usernameVal.(string)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": ErrBadRequest + ", incorrect token payload(missing username)"})
	}
	req := username

	if req == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": ErrBadRequest + ", username required"})
	}

	requestList, err := app.friendshipService.GetRequestList(ctx, req)
	if err != nil {
		logger.Error("app.getRequestList friendshipService.GetRequestList: ", err)
		return fiberInternalServerError(c)
	}

	return c.Status(fiber.StatusOK).JSON(requestList)
}
