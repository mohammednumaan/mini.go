package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/mohammednumaan/mini.go/go-cache/internal/cache"
)

type response struct {
	Success bool   `json:"success"`
	Value   any    `json:"value"`
	Error   string `json:"error"`
}

type cacheRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type metricsResponse struct {
	Hits      int `json:"hits"`
	Misses    int `json:"misses"`
	Evictions int `json:"evictions"`
	Size      int `json:"size"`
	Capacity  int `json:"capacity"`
}

func main() {
	lru := cache.NewLRUCache[string](5, 10 * time.Minute)

	http.HandleFunc("/cache", handleCache(lru))
	http.HandleFunc("/cache/", handleCache(lru))
	http.HandleFunc("/metrics", handleMetrics(lru))

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleCache(lru *cache.LRUCache[string]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		path := r.URL.Path

		switch r.Method {
		case http.MethodGet:

			key := path[len("/cache/"):]
			if key == "" {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(response{Success: false, Error: "key required"})
				return
			}
			if val, ok := lru.Get(key); ok {

				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(response{Success: true, Value: val})
			} else {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(response{Success: false, Error: "not found"})
			}

		case http.MethodPut:
			var req cacheRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Key == "" {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(response{Success: false, Error: "key and value required"})
				return
			}
			lru.Put(req.Key, req.Value)
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response{Success: true, Value: req.Value})

		case http.MethodDelete:
			key := path[len("/cache/"):]
			if key == "" {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(response{Success: false, Error: "key required"})
				return

			}
			
			w.WriteHeader(http.StatusOK)
			lru.Delete(key)

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(response{Success: false, Error: "method not allowed"})
		}
	}
}

func handleMetrics(lru *cache.LRUCache[string]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		stats := lru.Stats()

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(metricsResponse{
			Hits:      stats.Hits,
			Misses:    stats.Misses,
			Evictions: stats.Evictions,
			Size:      lru.Len(),
			Capacity:  lru.Capacity(),
		})
	}
}
