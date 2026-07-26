package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	healthHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("esperado status 200, obtido %d", rec.Code)
	}
}

func TestItemsHandlerCreateAndList(t *testing.T) {
	store := NewInventoryStore()
	handler := itemsHandler(store)

	item := Item{SKU: "ABC-1", Name: "Produto Teste", Quantity: 10}
	body, _ := json.Marshal(item)

	postReq := httptest.NewRequest(http.MethodPost, "/items", bytes.NewReader(body))
	postRec := httptest.NewRecorder()
	handler(postRec, postReq)

	if postRec.Code != http.StatusCreated {
		t.Fatalf("esperado status 201, obtido %d", postRec.Code)
	}

	getReq := httptest.NewRequest(http.MethodGet, "/items", nil)
	getRec := httptest.NewRecorder()
	handler(getRec, getReq)

	var items []Item
	if err := json.Unmarshal(getRec.Body.Bytes(), &items); err != nil {
		t.Fatalf("erro ao decodificar resposta: %v", err)
	}

	if len(items) != 1 || items[0].SKU != "ABC-1" {
		t.Fatalf("esperado 1 item com SKU ABC-1, obtido %+v", items)
	}
}

func TestItemsHandlerRejectsInvalidPayload(t *testing.T) {
	store := NewInventoryStore()
	handler := itemsHandler(store)

	req := httptest.NewRequest(http.MethodPost, "/items", bytes.NewReader([]byte("{}")))
	rec := httptest.NewRecorder()
	handler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("esperado status 400 para sku ausente, obtido %d", rec.Code)
	}
}
