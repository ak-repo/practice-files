package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/ak-repo/microservice-demo/pkg/config"
	"github.com/ak-repo/microservice-demo/pkg/logger"
	"github.com/gofiber/fiber/v2"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
)

type OrderReq struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type Product struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

func main() {
	cfg := config.Load()
	_ = logger.Init()

	app := fiber.New()
	app.Use(fiberlogger.New())

	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// POST /orders → Create a demo order
	app.Post("/orders", func(c *fiber.Ctx) error {
		var req OrderReq
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
		}
		if req.ProductID == "" || req.Quantity <= 0 {
			return c.Status(400).JSON(fiber.Map{"error": "product_id & quantity required"})
		}

		// 1️⃣ Validate token via auth-service
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{"error": "missing Authorization header"})
		}

		authURL := cfg.AuthBaseURL + "/validate"
		fmt.Println("→ Calling Auth Service:", authURL)

		vreq, _ := http.NewRequest("GET", authURL, nil)
		vreq.Header.Set("Authorization", authHeader)

		client := http.Client{Timeout: 5 * time.Second}
		vresp, err := client.Do(vreq)
		if err != nil {
			log.Println("❌ Auth-service not reachable:", err)
			return c.Status(500).JSON(fiber.Map{"error": "auth-service unreachable"})
		}
		defer vresp.Body.Close()

		if vresp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(vresp.Body)
			return c.Status(401).JSON(fiber.Map{"error": fmt.Sprintf("invalid token: %s", string(body))})
		}

		var userInfo map[string]interface{}
		if err := json.NewDecoder(vresp.Body).Decode(&userInfo); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to decode auth response"})
		}

		// 2️⃣ Fetch product details from product-service
		productURL := cfg.ProductBaseURL + "/products/" + req.ProductID
		fmt.Println("→ Fetching product from:", productURL)

		prodResp, err := client.Get(productURL)
		if err != nil {
			log.Println("❌ Product-service not reachable:", err)
			return c.Status(500).JSON(fiber.Map{"error": "product-service unreachable"})
		}
		defer prodResp.Body.Close()

		if prodResp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(prodResp.Body)
			return c.Status(400).JSON(fiber.Map{"error": fmt.Sprintf("product not found: %s", string(body))})
		}

		var prod Product
		if err := json.NewDecoder(prodResp.Body).Decode(&prod); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to decode product response"})
		}

		// 3️⃣ Construct and return order response
		order := map[string]interface{}{
			"order_id":   fmt.Sprintf("order-%d", time.Now().UnixNano()),
			"product":    prod,
			"quantity":   req.Quantity,
			"ordered_by": userInfo["username"],
		}

		return c.Status(201).JSON(order)
	})

	addr := ":" + cfg.OrderPort
	fmt.Println("✅ order-service listening on", addr)
	if err := app.Listen(addr); err != nil {
		log.Fatal("❌ Failed to start order-service:", err)
	}
}
