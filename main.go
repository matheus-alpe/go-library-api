package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type book struct {
	Id       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
}

var books = []book{
	{Id: "1", Title: "War and Peace", Author: "Leo Tolstoy", Quantity: 9},
	{Id: "2", Title: "The Great Gatsby", Author: "F. Scott Fitzgerald", Quantity: 2},
	{Id: "3", Title: "In Search of Lost Time", Author: "Marcel Proust", Quantity: 6},
}

func getBooks(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, books)
}

func getBookById(ctx *gin.Context) {
	id, found := ctx.Params.Get("id")
	if !found {
		return
	}

	book, err := findBookById(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, book)
}

func findBookById(id string) (*book, error) {
	for i, b := range books {
		if b.Id == id {
			return &books[i], nil
		}
	}

	return nil, errors.New("book not found")
}

func createBook(ctx *gin.Context) {
	var newBook book

	if err := ctx.BindJSON(&newBook); err != nil {
		// BindJSON will handle sending the error response
		return
	}

	books = append(books, newBook)
	ctx.JSON(http.StatusCreated, newBook)
}

func checkoutBook(ctx *gin.Context) {
	id, found := ctx.GetQuery("id")
	if !found {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "missing `id` query parameter"})
		return
	}
	
	book, err := findBookById(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	if book.Quantity  <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "book not available" })
		return
	}

	book.Quantity -= 1
	ctx.JSON(http.StatusOK, book)
}

func returnBook(ctx *gin.Context) {
	id, found := ctx.GetQuery("id")
	if !found {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "missing `id` query parameter"})
		return
	}
	
	book, err := findBookById(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	book.Quantity += 1
	ctx.JSON(http.StatusOK, book)
}


func main() {
	r := gin.Default()

	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Welcome to Library API",
		})
	})
	bookRoute := r.Group("/books")
	{
		bookRoute.GET("/", getBooks)
		bookRoute.GET("/:id", getBookById)
		bookRoute.POST("/", createBook)
	}
	r.PATCH("/checkout", checkoutBook)
	r.PATCH("/return", returnBook)

	r.Run("localhost:3232")
}
