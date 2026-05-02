package main

import (
	"encoding/json"
	"net/http"
	"log"
)

type User struct {
	ID int `json:"id"`
	Name string `json:"name"`
}

var (
	currId int
	users []User
)
func main() {
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/users", usersHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func healthHandler(w http.ResponseWriter, r *http.Request){
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK\n"))
}

func usersHandler(w http.ResponseWriter, r *http.Request){

	switch r.Method {
	case "GET":
		w.Header().Set("Content-Type", "application/json")
    	json.NewEncoder(w).Encode(users)

	case "POST":
		var user User
		err := json.NewDecoder(r.Body).Decode(&user)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		currId++
		user.ID = currId
		users = append(users, user)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}
