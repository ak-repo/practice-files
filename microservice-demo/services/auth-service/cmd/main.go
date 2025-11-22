package main

import (
	"fmt"
	"time"

	"github.com/ak-repo/microservice-demo/pkg/config"
	"github.com/ak-repo/microservice-demo/pkg/jwt"
	"github.com/ak-repo/microservice-demo/pkg/logger"
	"github.com/gofiber/fiber/v2"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger" // ✅ correct import
)

var users = map[string]string{} // username -> password (demo only)

func main() {
	cfg := config.Load()
	_ = logger.Init()

	app := fiber.New()
	app.Use(fiberlogger.New())

	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	app.Post("/signup", func(c *fiber.Ctx) error {
		var body struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
		}
		if body.Username == "" || body.Password == "" {
			return c.Status(400).JSON(fiber.Map{"error": "username/password required"})
		}
		users[body.Username] = body.Password
		return c.Status(201).JSON(fiber.Map{"message": "user created"})
	})

	app.Post("/login", func(c *fiber.Ctx) error {
		var body struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
		}
		pw, ok := users[body.Username]
		if !ok || pw != body.Password {
			return c.Status(401).JSON(fiber.Map{"error": "invalid credentials"})
		}
		token, err := jwt.GenerateToken(cfg.JwtSecret, body.Username, body.Username, time.Hour*24)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "could not create token"})
		}
		return c.JSON(fiber.Map{"token": token})
	})

	// validate token endpoint for other services
	app.Get("/validate", func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "" {
			return c.Status(401).JSON(fiber.Map{"error": "missing authorization"})
		}
		var tokenStr string
		fmt.Sscanf(auth, "Bearer %s", &tokenStr)
		cl, err := jwt.ParseToken(cfg.JwtSecret, tokenStr)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "invalid token"})
		}
		return c.JSON(fiber.Map{"user_id": cl.UserID, "username": cl.Username})
	})

	addr := ":" + cfg.AuthPort
	fmt.Println("auth-service listening on", addr)
	app.Listen(addr)
}
