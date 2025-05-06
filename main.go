package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

const (
	host         = "localhost"
	port         = 5432
	databaseName = "mydatabase"
	username     = "myuser"
	password     = "mypassword"
)

var db *sql.DB

type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

func main() {
	// Connection string
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, username, password, databaseName)

	sdb, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	db = sdb
	// Clost connection
	defer db.Close()
	// Check the connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected!")

	app := fiber.New()

	app.Get("/hello", func(c *fiber.Ctx) error {
		return c.SendString("Hello World!")
	})
	app.Get("/products/:id", getProductHandler)
	app.Get("/products/", getProductsHandler)
	app.Post("/products/", createProductHandler)
	app.Put("/products/:id", updateProductHandler)
	app.Delete("/products/:id", deleteProductHandler)
	app.Listen(":8080")
}

func getProductHandler(c *fiber.Ctx) error {
	productId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid product ID")
	}
	product, err := getProduct(productId)
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Product Not Found")
	}
	return c.JSON(product)
}

func getProductsHandler(c *fiber.Ctx) error {
	products, err := getProducts()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to fetch products")
	}
	return c.JSON(products)
}

func createProductHandler(c *fiber.Ctx) error {
	p := new(Product)
	if err := c.BodyParser(p); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	if p.Name == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid product data")
	}

	err := createProduct(p)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to create product")
	}

	return c.Status(fiber.StatusCreated).JSON(p)
}

func updateProductHandler(c *fiber.Ctx) error {
	productId, err := strconv.Atoi(c.Params("id"))
	p := new(Product)
	if err := c.BodyParser(p); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	product, err := updateProduct(productId, p)
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Product Not Found")
	}
	return c.JSON(product)
}

func deleteProductHandler(c *fiber.Ctx) error {
	productId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid product ID")
	}
	err = deleteProduct(productId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to delete product")
	}
	return c.SendStatus(fiber.StatusNoContent)
}
