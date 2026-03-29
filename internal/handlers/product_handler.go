package handlers

import (
	"e-shop-modal/internal/models"
	"e-shop-modal/internal/server"
	"e-shop-modal/internal/services"
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

func (h *ProductHandler) GetProducts(c *server.Context) {
	products, err := h.services.ObtenerTodosLosProducts()
	if err != nil {
		c.JSONResponse(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSONResponse(http.StatusOK, products)
}

func (h *ProductHandler) CreateProduct(c *server.Context) {
	var product models.Product
	err := c.BindJSON(&product)
	if err != nil {
		c.JSONResponse(http.StatusBadRequest, err.Error())
		return
	}
	created, err := h.services.SubirProduct(product)
	if err != nil {
		c.JSONResponse(http.StatusInternalServerError, err.Error())
	}
	c.JSONResponse(http.StatusCreated, created)
}

func (h *ProductHandler) GetProductByID(c *server.Context) {
	id, err := getIDFromPath(c)
	if err != nil {
		c.JSONResponse(http.StatusBadRequest, map[string]string{
			"error": "ID inválido",
		})
		return
	}
	product, err := h.services.ObtenerProductPorID(id)
	if err != nil {
		c.JSONResponse(http.StatusNotFound, map[string]string{
			"error": "Producto no encontrado",
		})
		return
	}

	c.JSONResponse(http.StatusOK, product)
}
func (h *ProductHandler) UpdateProduct(c *server.Context) {

	id, err := getIDFromPath(c)
	if err != nil {
		c.JSONResponse(http.StatusBadRequest, map[string]string{
			"error": "ID inválido",
		})
		return
	}

	var product models.Product
	if err := c.BindJSON(&product); err != nil {
		c.JSONResponse(http.StatusBadRequest, map[string]string{
			"error": "Datos inválidos",
		})
		return
	}

	updated, err := h.services.ModificarProduct(id, product)
	if err != nil {
		c.JSONResponse(http.StatusInternalServerError, map[string]string{
			"error": "No se pudo actualizar",
		})
		return
	}

	c.JSONResponse(http.StatusOK, updated)
}

func (h *ProductHandler) DeleteProduct(c *server.Context) {

	id, err := getIDFromPath(c)
	if err != nil {
		c.JSONResponse(http.StatusBadRequest, map[string]string{
			"error": "ID inválido",
		})
		return
	}

	err = h.services.QuitarProduct(id)
	if err != nil {
		c.JSONResponse(http.StatusInternalServerError, map[string]string{
			"error": "No se pudo eliminar",
		})
		return
	}

	c.Status(http.StatusNoContent)
}

func getIDFromPath(c *server.Context) (int, error) {
	path := c.Request.URL.Path

	// ejemplo: /products/5
	parts := strings.Split(path, "/")

	idStr := parts[len(parts)-1]

	return strconv.Atoi(idStr)
}
