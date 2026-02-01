package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"clothesstore/models"
	"clothesstore/repository"
)

func ProductsHandler(store *repository.MemoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.Method {
		case http.MethodGet:
			products, _ := store.GetAllProducts()
			json.NewEncoder(w).Encode(products)

		case http.MethodPost:
			var p models.Product
			if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
				writeJSON(w, http.StatusBadRequest, map[string]string{
					"error": "invalid json",
				})
				return
			}
			created, _ := store.CreateProduct(p)
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(created)

		default:
			writeJSON(w, http.StatusMethodNotAllowed, map[string]string{
				"error": "method not allowed",
			})
		}
	}
}

func ProductByIDHandler(store *repository.MemoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		idStr := strings.TrimPrefix(r.URL.Path, "/products/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{
				"error": "invalid id",
			})
			return
		}

		switch r.Method {
		case http.MethodGet:
			p, err := store.GetProductByID(id)
			if err != nil {
				writeJSON(w, http.StatusNotFound, map[string]string{
					"error": "not found",
				})
				return
			}
			json.NewEncoder(w).Encode(p)

		case http.MethodPut:
			var p models.Product
			if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
				writeJSON(w, http.StatusBadRequest, map[string]string{
					"error": "invalid json",
				})
				return
			}
			updated, err := store.UpdateProduct(id, p)
			if err != nil {
				writeJSON(w, http.StatusNotFound, map[string]string{
					"error": "not found",
				})
				return
			}
			json.NewEncoder(w).Encode(updated)

		case http.MethodDelete:
			if err := store.DeleteProduct(id); err != nil {
				writeJSON(w, http.StatusNotFound, map[string]string{
					"error": "not found",
				})
				return
			}
			json.NewEncoder(w).Encode(map[string]string{"result": "deleted"})

		default:
			writeJSON(w, http.StatusMethodNotAllowed, map[string]string{
				"error": "method not allowed"})
		}
	}
}
