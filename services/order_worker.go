package services

import (
	"time"

	"clothesstore/repository"
)

type OrderWorker struct {
	store *repository.MemoryStore
	queue chan int
}

func NewOrderWorker(store *repository.MemoryStore) *OrderWorker {
	return &OrderWorker{
		store: store,
		queue: make(chan int, 100),
	}
}

func (w *OrderWorker) Start() {
	go func() {
		for orderID := range w.queue {
			time.Sleep(2 * time.Second)
			_ = w.store.UpdateOrderStatus(orderID, "processed")
		}
	}()
}

func (w *OrderWorker) Enqueue(orderID int) {
	w.queue <- orderID
}
