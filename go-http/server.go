package main

import (
	"github.com/mohammednumaan/mini.go/go-http/internal/books"
	"log"
	"net/http"
)

func main() {
	router := http.NewServeMux()

	router.HandleFunc("GET /books", books.GetBooks)
	router.HandleFunc("GET /books/{id}", books.GetBookByID)
	router.HandleFunc("POST /books", books.CreateBook)

	router.HandleFunc("PUT /books/{id}", books.UpdateBook)
	router.HandleFunc("DELETE /books/{id}", books.DeleteBook)

	log.Println("[Server]: Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
