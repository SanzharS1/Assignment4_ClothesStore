package router

import (
	"net/http"

	"clothesstore/handlers"
	"clothesstore/repository"
	"clothesstore/services"
)

func SetupRouter(store *repository.MemoryStore, worker *services.OrderWorker) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", handlers.HealthHandler)

	mux.HandleFunc("/products", handlers.ProductsHandler(store))
	mux.HandleFunc("/products/", handlers.ProductByIDHandler(store))

	mux.HandleFunc("/users/register", handlers.RegisterUserHandler(store))
	mux.HandleFunc("/orders", handlers.OrdersHandler(store, worker))

	return mux
}
