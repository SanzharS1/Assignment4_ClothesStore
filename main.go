package main

import (
	"fmt"
	"log"
	"net/http"

	"clothesstore/repository"
	"clothesstore/router"
	"clothesstore/services"
)

func main() {
	store := repository.NewMemoryStore()

	worker := services.NewOrderWorker(store)
	worker.Start()

	r := router.SetupRouter(store, worker)

	fmt.Println("ClothesStore server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
