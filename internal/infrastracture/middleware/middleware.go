package middleware

import (
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing token"})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token format"})
		}

		tokenString := parts[1]
		secretKey := os.Getenv("JWT_SECRET")

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid claims"})
		}

		if exp, ok := claims["exp"].(float64); ok {
			if int64(exp) < time.Now().Unix() {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "token expired"})
			}
		}

		usernameVal, exists := claims["username"]
		if !exists {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token claims"})
		}
		username, ok := usernameVal.(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token claims"})
		}
		c.Locals("username", username)

		idVal, exists := claims["id"]
		if !exists {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token claims"})
		}
		id, ok := idVal.(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token claims"})
		}
		c.Locals("id", id)

		isAdmin := false
		if isAdminVal, exists := claims["isAdmin"]; exists {
			if isAdminBool, ok := isAdminVal.(bool); ok {
				isAdmin = isAdminBool
			} else if isAdminFloat, ok := isAdminVal.(float64); ok {
				isAdmin = isAdminFloat == 1
			} else {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token claims"})
			}
		}
		c.Locals("isAdmin", isAdmin)
		return c.Next()
	}
}

func AdminMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		isAdmin, ok := c.Locals("isAdmin").(bool)
		if !ok || !isAdmin {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "access denied"})
		}

		return c.Next()
	}
}
