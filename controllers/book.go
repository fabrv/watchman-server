package controllers

import (
	"github.com/fabrv/watchman-server/database"
	"github.com/fabrv/watchman-server/models"
	"github.com/fabrv/watchman-server/utils"
	"github.com/gofiber/fiber/v2"
)

func GetBooks(c *fiber.Ctx) error {
	db := database.DBConn
	var books []models.Book
	db.Find(&books)
	return c.JSON(books)
}

func GetBook(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn
	var book models.Book
	db.First(&book, id)
	return c.JSON(book)
}

func AddBook(c *fiber.Ctx) error {
	var book models.Book
	if err := c.BodyParser(&book); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := utils.ValidateStruct(book)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	db := database.DBConn
	err := db.Create(&book)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error.Error(),
		})
	}
	return c.JSON(book)
}

func UpdateBook(c *fiber.Ctx) error {
	id := c.Params("id")
	var book models.Book
	if err := c.BodyParser(&book); err != nil {
		return c.Status(503).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	db := database.DBConn
	db.Model(&book).Where("id = ?", id).Updates(book)
	return c.JSON(fiber.Map{
		"message": "Book updated",
	})
}

func DeleteBook(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn
	var book models.Book
	db.First(&book, id)
	if book.ID == 0 {
		return c.Status(404).JSON(fiber.Map{
			"error": "Book not found",
		})
	}
	db.Delete(&book)
	return c.JSON(fiber.Map{
		"message": "Book deleted",
	})
}
