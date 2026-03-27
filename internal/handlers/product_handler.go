package handlers

import (
	"e-shop-modal/internal/models"
	"e-shop-modal/internal/services"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type ProductHandler struct {
	services *services.ProductService
}

func NewProductHandler(s *services.ProductService) *ProductHandler {
	return &ProductHandler{
		services: s,
	}
}

func (h *ProductHandler) HandleProducts(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		products, err := h.services.ObtenerTodosLosProducts()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(products)

	case http.MethodPost:
		var product models.Product
		err := json.NewDecoder(r.Body).Decode(&product)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		created, err := h.services.SubirProduct(product)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(created)
	default:
		http.Error(w, "Metodo no disponible", http.StatusMethodNotAllowed)
	}
}

func (h *ProductHandler) HandleProductByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/producto/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Algo ha fallado", http.StatusBadRequest)
	}
	switch r.Method {
	case http.MethodGet:
		product, err := h.services.ObtenerProductPorID(id)
		if err != nil {
			http.Error(w, "Producto no encontrado", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(product)

	case http.MethodPut:
		var product models.Product
		err := json.NewDecoder(r.Body).Decode(&product)
		if err != nil {
			http.Error(w, "Input invalido", http.StatusBadRequest)
			return
		}
		updated, err := h.services.ModificarProduct(id, product)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(updated)

	case http.MethodDelete:
		if err := h.services.QuitarProduct(id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	default:
		http.Error(w, "Metodo no disponible", http.StatusMethodNotAllowed)
	}
}
