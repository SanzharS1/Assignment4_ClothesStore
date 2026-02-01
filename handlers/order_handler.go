package handlers

import (
	"encoding/json"
	"net/http"

	"clothesstore/models"
	"clothesstore/repository"
	"clothesstore/services"
)

type CreateOrderRequest struct {
	UserID int                `json:"user_id"`
	Items  []models.OrderItem `json:"items"`
}

func OrdersHandler(store *repository.MemoryStore, worker *services.OrderWorker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.Method {
		case http.MethodGet:
			orders, _ := store.GetAllOrders()
			json.NewEncoder(w).Encode(orders)

		case http.MethodPost:
			var req CreateOrderRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
				return
			}

			order, err := store.CreateOrder(req.UserID, req.Items)
			if err != nil {
				http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
				return
			}

			worker.Enqueue(order.ID)

			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(order)

		default:
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		}
	}
}
