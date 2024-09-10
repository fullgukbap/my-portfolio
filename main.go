package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

const secretKey = "6327080c49cef882ecf3b8273a960e05688843a0a16a242822c960eb4e2a5a3955e4079073e5327e9f073cd93505605b0a5d9e3bfc3d19c0695a43597041003a"
const filePath = "./files/portfolio.pdf"

var pdfBytes []byte

func init() {
	pdfData, err := os.ReadFile(filePath)
	if err != nil {
		log.Panicf("failed to read file: %v", err)
	}

	pdfBytes = pdfData
}

func main() {
	// Fiber 앱 생성
	app := fiber.New(fiber.Config{})

	app.Use(recover.New())

	// 압축 미들웨어 추가
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed, // 압축 속도를 최적화하려면 BestSpeed 사용

	}))

	app.Use(helmet.New())

	app.All("/*", func(c *fiber.Ctx) error {

		key := c.Query("key", "")
		if key == "" || key != secretKey {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		// 캐시 헤더 설정
		c.Set("Cache-Control", "public, max-age=86400") // 1일 동안 캐싱

		c.Set("Content-Type", "application/pdf")
		c.Set("Content-Disposition", "inline; filename=한준범 포토폴리오.pdf")

		return c.Send(pdfBytes)
	})

	if err := app.Listen(":8080"); err != nil {
		log.Fatal(err)
	}
}
