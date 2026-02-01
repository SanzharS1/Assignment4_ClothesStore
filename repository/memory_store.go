package repository

import (
	"errors"
	"sync"
	"time"

	"clothesstore/models"
)

var (
	ErrNotFound        = errors.New("not found")
	ErrBadRequest      = errors.New("bad request")
	ErrOutOfStock      = errors.New("out of stock")
	ErrUserNotFound    = errors.New("user not found")
	ErrProductNotFound = errors.New("product not found")
)

type MemoryStore struct {
	mu sync.RWMutex

	products map[int]models.Product
	users    map[int]models.User
	orders   map[int]models.Order

	nextProductID   int
	nextUserID      int
	nextOrderID     int
	nextOrderItemID int
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		products:        make(map[int]models.Product),
		users:           make(map[int]models.User),
		orders:          make(map[int]models.Order),
		nextProductID:   1,
		nextUserID:      1,
		nextOrderID:     1,
		nextOrderItemID: 1,
	}
}

// -------------------- PRODUCTS (CRUD) --------------------
func (s *MemoryStore) CreateProduct(p models.Product) (models.Product, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	p.ID = s.nextProductID
	s.nextProductID++
	s.products[p.ID] = p
	return p, nil
}

func (s *MemoryStore) GetAllProducts() ([]models.Product, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]models.Product, 0, len(s.products))
	for _, p := range s.products {
		out = append(out, p)
	}
	return out, nil
}

func (s *MemoryStore) GetProductByID(id int) (models.Product, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	p, ok := s.products[id]
	if !ok {
		return models.Product{}, ErrNotFound
	}
	return p, nil
}

func (s *MemoryStore) UpdateProduct(id int, update models.Product) (models.Product, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.products[id]; !ok {
		return models.Product{}, ErrNotFound
	}
	update.ID = id
	s.products[id] = update
	return update, nil
}

func (s *MemoryStore) DeleteProduct(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.products[id]; !ok {
		return ErrNotFound
	}
	delete(s.products, id)
	return nil
}

// -------------------- USERS --------------------
func (s *MemoryStore) RegisterUser(u models.User) (models.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	u.ID = s.nextUserID
	s.nextUserID++
	s.users[u.ID] = u
	return u, nil
}

// -------------------- ORDERS --------------------
func (s *MemoryStore) CreateOrder(userID int, items []models.OrderItem) (models.Order, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.users[userID]; !ok {
		return models.Order{}, ErrUserNotFound
	}

	// validate + stock check
	for _, it := range items {
		if it.Quantity <= 0 {
			return models.Order{}, ErrBadRequest
		}
		p, ok := s.products[it.ProductID]
		if !ok {
			return models.Order{}, ErrProductNotFound
		}
		if p.Stock < it.Quantity {
			return models.Order{}, ErrOutOfStock
		}
	}

	// decrease stock
	for _, it := range items {
		p := s.products[it.ProductID]
		p.Stock -= it.Quantity
		s.products[it.ProductID] = p
	}

	// create order
	order := models.Order{
		ID:        s.nextOrderID,
		UserID:    userID,
		Status:    "created",
		CreatedAt: time.Now(),
		Items:     make([]models.OrderItem, 0, len(items)),
	}
	s.nextOrderID++

	for _, it := range items {
		it.ID = s.nextOrderItemID
		s.nextOrderItemID++
		it.OrderID = order.ID
		order.Items = append(order.Items, it)
	}

	s.orders[order.ID] = order
	return order, nil
}

func (s *MemoryStore) GetAllOrders() ([]models.Order, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]models.Order, 0, len(s.orders))
	for _, o := range s.orders {
		out = append(out, o)
	}
	return out, nil
}

func (s *MemoryStore) UpdateOrderStatus(orderID int, status string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	o, ok := s.orders[orderID]
	if !ok {
		return ErrNotFound
	}
	o.Status = status
	s.orders[orderID] = o
	return nil
}
