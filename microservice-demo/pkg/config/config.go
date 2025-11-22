package config

import "os"

type Config struct {
	AuthPort       string
	ProductPort    string
	OrderPort      string
	JwtSecret      string
	AuthBaseURL    string
	ProductBaseURL string
}

func Load() *Config {

	cfg := &Config{
		AuthPort:       get("AUTH_PORT", "8001"),
		ProductPort:    get("PRODUCT_PORT", "8002"),
		OrderPort:      get("ORDER_PORT", "8003"),
		JwtSecret:      get("JWT_SECRET", "devsecret"),
		AuthBaseURL:    get("AUTH_BASE_URL", " http://127.0.0.1:8001"),
		ProductBaseURL: get("PRODUCT_BASE_URL", " http://127.0.0.1:8002"),
	}

	return cfg
}

func get(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
