package main

import (
	"log"
	"slices"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type Book struct {
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

var books []Book

func main() {
	app := fiber.New()

	books = append(books, Book{Id: 1, Title: "Book 1", Author: "Author 1"})
	books = append(books, Book{Id: 2, Title: "Book 2", Author: "Author 2"})

	app.Get("/books", getBooks)
	app.Get("/books/:id", getBook)
	app.Post("/books", createBook)
	app.Put("/books/:id", updateBook)
	app.Delete("/books/:id", deleteBook)

	log.Fatal(app.Listen(":8080"))
}

func getBooks(c *fiber.Ctx) error {
	return c.JSON(books)
}

func getBook(c *fiber.Ctx) error {
	bookId, error := strconv.Atoi(c.Params("id"))
	if error != nil {
		return c.Status(fiber.StatusBadRequest).SendString(error.Error())
	}

	for _, book := range books {
		if book.Id == bookId {
			return c.JSON(book)
		}
	}
	return c.Status(fiber.StatusNotFound).SendString("Book not found")
}

func createBook(c *fiber.Ctx) error {
	book := new(Book)
	if err := c.BodyParser(book); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	for _, b := range books {
		if b.Id == book.Id {
			return c.Status(fiber.StatusBadRequest).SendString("Book with this ID already exists")
		}
	}
	books = append(books, *book)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Book created successfully",
		"book":    *book,
	})
}

func updateBook(c *fiber.Ctx) error {
	bookId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	bookUpdate := new(Book)
	if err := c.BodyParser(bookUpdate); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	for i, book := range books {
		if book.Id == bookId {
			books[i].Title = bookUpdate.Title
			books[i].Author = bookUpdate.Author
			return c.JSON(fiber.Map{
				"message": "Book updated successfully",
				"book":    books[i],
			})
		}
	}
	return c.Status(fiber.StatusNotFound).SendString("Book not found")
}

func deleteBook(c *fiber.Ctx) error {
	bookId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if len(books) == 0 {
		return c.Status(fiber.StatusNotFound).SendString("No books found")
	}
	for i, book := range books {
		if book.Id == bookId {
			books = slices.Delete(books, i, i+1)
			return c.JSON(fiber.Map{
				"message": "Book deleted successfully",
			})
		}
	}
	return c.Status(fiber.StatusNotFound).SendString("Book not found")
}
