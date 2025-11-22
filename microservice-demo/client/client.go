package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func main() {
	baseAuth := "http://localhost:8001"
	baseProduct := "http://localhost:8002"
	baseOrder := "http://localhost:8003"

	fmt.Println("🔹 Testing microservices...")

	// 1️⃣ Call AUTH-SERVICE (validate token)
	token := "Bearer demo-token-123"
	authURL := baseAuth + "/validate"
	req, _ := http.NewRequest("GET", authURL, nil)
	req.Header.Set("Authorization", token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("❌ Auth service error:", err)
	} else {
		body, _ := io.ReadAll(resp.Body)
		fmt.Println("✅ Auth-service response:", string(body))
		resp.Body.Close()
	}

	// 2️⃣ Call PRODUCT-SERVICE (get product)
	productURL := baseProduct + "/products/1"
	resp2, err := http.Get(productURL)
	if err != nil {
		fmt.Println("❌ Product service error:", err)
	} else {
		body, _ := io.ReadAll(resp2.Body)
		fmt.Println("✅ Product-service response:", string(body))
		resp2.Body.Close()
	}

	// 3️⃣ Call ORDER-SERVICE (create order)
	orderURL := baseOrder + "/orders"
	orderBody := map[string]interface{}{
		"product_id": "1",
		"quantity":   2,
	}
	jsonData, _ := json.Marshal(orderBody)
	req3, _ := http.NewRequest("POST", orderURL, bytes.NewBuffer(jsonData))
	req3.Header.Set("Content-Type", "application/json")
	req3.Header.Set("Authorization", token)

	resp3, err := http.DefaultClient.Do(req3)
	if err != nil {
		fmt.Println("❌ Order service error:", err)
	} else {
		body, _ := io.ReadAll(resp3.Body)
		fmt.Println("✅ Order-service response:", string(body))
		resp3.Body.Close()
	}

	fmt.Println("\n✅ All requests completed.")
}
