package main

import (
	"LRU-cache-project/server/internal"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

var cache *internal.LRUCache

func main() {
	cache = internal.NewLRUCache(5)
	cache.StartCleanupTask()

	router := fiber.New()
	router.Use(cors.New())

	// router.Use(logger.New(logger.Config{
	// 	Format: "[${time}] ${status} - ${latency} ${method} ${path}\n",
	// }))

	router.Get("/cache/:key", getHandler)
	router.Post("/cache", setHandler)
	router.Delete("/cache/:key", deleteHandler)
	router.Get("/cache", getAllHandler)

	// log.SetOutput(os.Stdout)
	log.Println("Server is starting...")

	// Start server
	log.Fatal(router.Listen(":8080"))

	router.Listen(":8080")
}

func getHandler(c *fiber.Ctx) error {
	key := c.Params("key")
	fmt.Println("Cache items:", key)
	value, found := cache.Get(key)
	if !found {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Key not found"})
	}
	return c.JSON(fiber.Map{"key": key, "value": value})
}

func setHandler(c *fiber.Ctx) error {
	var input struct {
		Key        string      `json:"key"`
		Value      interface{} `json:"value"`
		Expiration int         `json:"expiration"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	cache.Set(input.Key, input.Value, time.Duration(input.Expiration)*time.Second)
	return c.JSON(fiber.Map{
		"message": "Key set successfully",
		"key":     input.Key,
		"value":   input.Value,
	})
}

func deleteHandler(c *fiber.Ctx) error {
	key := c.Params("key")
	cache.Delete(key)
	return c.JSON(fiber.Map{"message": "Key deleted successfully"})
}

func getAllHandler(c *fiber.Ctx) error {
	items := cache.GetAll()
	fmt.Println("Cache items:", items)
	return c.JSON(items)
}
