package handlers

import (
	"e-shop-modal/internal/dto"
	"e-shop-modal/internal/server"
	"e-shop-modal/internal/services"
	"net/http"
)

type PaymentHandler struct {
	service *services.PaymentService
}

func NewPaymentHandler(s *services.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		service: s,
	}
}

func (h *PaymentHandler) CreateCheckout(c *server.Context) {

	var req *dto.CheckoutItem

	// Leer JSON usando tu Context
	if err := c.BindJSON(&req); err != nil {
		c.JSONResponse(http.StatusBadRequest, "JSON inválido")
		return
	}

	// Usar el context real de la request
	resp, err := h.service.CreatePreference(c.Context(), req)
	if err != nil {
		c.JSONResponse(http.StatusInternalServerError, "Error creando preferencia")
		return
	}

	// Devolver init_point
	c.JSONResponse(http.StatusOK, map[string]interface{}{
		"checkout_url": resp.InitPoint,
	})
}

func (h *PaymentHandler) ConfirmWebhook(c *server.Context) {

	var body map[string]interface{}

	err := c.BindJSON(&body)
	if err != nil {
		JSONError(c, http.StatusBadRequest, "JSON invalido")
		return
	}

	// Obtiene el payment_id y lo procesa
	data := body["data"].(map[string]interface{})
	paymentID := int(data["id"].(float64))

	err = h.service.ProcessWebhook(paymentID)
	if err != nil {
		JSONError(c, http.StatusInternalServerError, "Error al aceptar pago")
		return
	}

	c.JSONResponse(http.StatusOK, "Pago exitoso")
}
