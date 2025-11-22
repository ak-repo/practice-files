package main

import (
	"fmt"
	"sync"

	"github.com/ak-repo/microservice-demo/pkg/config"
	"github.com/ak-repo/microservice-demo/pkg/logger"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger" // ✅ correct import


	"github.com/gofiber/fiber/v2"
)

type Product struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

var (
	mu       sync.Mutex
	products = map[string]*Product{}
)

func main() {
	cfg := config.Load()
	_ = logger.Init()

	app := fiber.New()
	app.Use(fiberlogger.New())

	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	app.Get("/products", func(c *fiber.Ctx) error {
		mu.Lock()
		defer mu.Unlock()
		list := make([]*Product, 0, len(products))
		for _, p := range products {
			list = append(list, p)
		}
		return c.JSON(list)
	})

	app.Post("/products", func(c *fiber.Ctx) error {
		var p Product
		if err := c.BodyParser(&p); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
		}
		if p.ID == "" {
			return c.Status(400).JSON(fiber.Map{"error": "id required"})
		}
		mu.Lock()
		products[p.ID] = &p
		mu.Unlock()
		return c.Status(201).JSON(p)
	})

	app.Get("/products/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		mu.Lock()
		p, ok := products[id]
		mu.Unlock()
		if !ok {
			return c.Status(404).JSON(fiber.Map{"error": "not found"})
		}
		return c.JSON(p)
	})

	addr := ":" + cfg.ProductPort
	fmt.Println("product-service listening on", addr)
	app.Listen(addr)

}
