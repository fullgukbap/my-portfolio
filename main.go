package main

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

var pdfBytes []byte
var config *Config

func init() {
	c, err := NewConfig("./configs/.toml")
	if err != nil {
		log.Panicf("failed to open config: %v", err)
	}

	config = c

	pdfData, err := os.ReadFile(c.File.PortfolioPath)
	if err != nil {
		log.Panicf("failed to read file: %v", err)
	}

	pdfBytes = pdfData
}

func main() {

	// Fiber 앱 생성
	app := fiber.New(fiber.Config{
		GETOnly: true,
		Prefork: true,
	})

	app.Use(recover.New())

	// 압축 미들웨어 추가
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed, // 압축 속도를 최적화하려면 BestSpeed 사용
	}))

	app.Use(helmet.New())

	app.Use(cache.New(cache.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.Query("noCache") == "true"
		},
		Expiration:   12 * time.Hour, // 12 h
		CacheControl: true,
	}))

	// 만약 나중에 요청이 많이 올 시  아래 미들웨어 추가
	// Etga, Limiter

	app.All("/*", func(c *fiber.Ctx) error {
		key := c.Query("key", "")
		if key == "" || key != config.Auth.Key {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		c.Set("Content-Type", "application/pdf")
		c.Set("Content-Disposition", "inline; filename=한준범 포토폴리오.pdf")

		return c.Send(pdfBytes)
	})

	if err := app.Listen(config.Http.Port); err != nil {
		log.Fatal(err)
	}
}
