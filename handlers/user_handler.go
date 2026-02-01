package handlers

import (
	"encoding/json"
	"net/http"

	"clothesstore/models"
	"clothesstore/repository"
)

func RegisterUserHandler(store *repository.MemoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodPost {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		var u models.User
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, `{"error":"invalid json body"}`, http.StatusBadRequest)
			return
		}

		if u.Name == "" || u.Email == "" || u.Password == "" {
			http.Error(w, `{"error":"name, email and password are required"}`, http.StatusBadRequest)
			return
		}

		created, _ := store.RegisterUser(u)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(created)
	}
}
