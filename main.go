package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
)

// Item representa um item de inventario
type Item struct {
	SKU      string `json:"sku"`
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}

type InventoryStore struct {
	mu    sync.RWMutex
	items map[string]Item
}

func NewInventoryStore() *InventoryStore {
	return &InventoryStore{items: make(map[string]Item)}
}

func (s *InventoryStore) Upsert(item Item) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[item.SKU] = item
}

func (s *InventoryStore) Get(sku string) (Item, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	item, ok := s.items[sku]
	return item, ok
}

func (s *InventoryStore) List() []Item {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]Item, 0, len(s.items))
	for _, item := range s.items {
		result = append(result, item)
	}
	return result
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func itemsHandler(store *InventoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(store.List())
		case http.MethodPost:
			var item Item
			if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
				http.Error(w, "payload invalido", http.StatusBadRequest)
				return
			}
			if item.SKU == "" {
				http.Error(w, "sku e obrigatorio", http.StatusBadRequest)
				return
			}
			store.Upsert(item)
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(item)
		default:
			http.Error(w, "metodo nao suportado", http.StatusMethodNotAllowed)
		}
	}
}

func main() {
	store := NewInventoryStore()

	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/items", itemsHandler(store))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("gcp-cloudrun-inventory ouvindo na porta %s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}
