package main

import (
	"fmt"
	"log"
	"net/http"

	"clothesstore/models"
	"clothesstore/repository"
	"clothesstore/router"
	"clothesstore/services"
)

func main() {
	store := repository.NewMemoryStore()
	store.CreateProduct(models.Product{Name: "T-Shirt", Price: 4990, Size: "M", Category: "tops", Stock: 10})
	store.CreateProduct(models.Product{Name: "Jeans", Price: 12990, Size: "L", Category: "bottoms", Stock: 5})

	worker := services.NewOrderWorker(store)
	worker.Start()

	r := router.SetupRouter(store, worker)

	fmt.Println("ClothesStore server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
