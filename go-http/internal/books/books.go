package books

import (
	"encoding/json"
	"github.com/mohammednumaan/mini.go/go-http/internal/response"
	"net/http"
	"strconv"
	"sync"
)

type Book struct {
	ID     int `json:"id"`
	Isbn   string `json:"isbn"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

type InMemoryBookStorage struct {
	books []Book ``
	mu    sync.Mutex
}

func createBookStorage() InMemoryBookStorage {
	return InMemoryBookStorage{books: []Book{}}
}

var allBooks = createBookStorage()

func GetBooks(w http.ResponseWriter, r *http.Request) {
	allBooks.mu.Lock()
	defer allBooks.mu.Unlock()
	response.SendValidResponse(w, "Books retrieved successfully", allBooks.books)
}

func GetBookByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id, err := strconv.Atoi(r.PathValue("id"))

	if err != nil {
		response.SendErrorResponse(w, "Invalid request", http.StatusBadRequest)
		return
	}

	for _, book := range allBooks.books {
		if book.ID == id {
			allBooks.mu.Lock()
			response.SendValidResponse(w, "Book retrieved successfully", book)
			allBooks.mu.Unlock()
			return
		}
	}

	response.SendErrorResponse(w, "Book not found", http.StatusNotFound)
}

func CreateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book Book
	json.NewDecoder(r.Body).Decode(&book)

	allBooks.mu.Lock()
	defer allBooks.mu.Unlock()

	// ideally this should be uuid
	book.ID = len(allBooks.books) + 1
	allBooks.books = append(allBooks.books, book)
	response.SendValidResponse(w, "Book created successfully", book)
}

func UpdateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id, err := strconv.Atoi(r.PathValue("id"))

	if err != nil {
		response.SendErrorResponse(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var updatedBook Book
	json.NewDecoder(r.Body).Decode(&updatedBook)

	for i, book := range allBooks.books {
		if book.ID == id {
			allBooks.mu.Lock()
			allBooks.books[i] = updatedBook
			allBooks.mu.Unlock()

			response.SendValidResponse(w, "Book updated successfully", updatedBook)
			return
		}
	}

	response.SendErrorResponse(w, "Book not found", http.StatusNotFound)
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id, err := strconv.Atoi(r.PathValue("id"))

	if err != nil {
		response.SendErrorResponse(w, "Invalid request", http.StatusBadRequest)
		return
	}

	for i, book := range allBooks.books {
		if book.ID == id {
			deletedBook := book
			allBooks.mu.Lock()
			allBooks.books = append(allBooks.books[:i], allBooks.books[i+1:]...)
			allBooks.mu.Unlock()

			response.SendValidResponse(w, "Book deleted successfully", deletedBook)
			return
		}
	}
}
