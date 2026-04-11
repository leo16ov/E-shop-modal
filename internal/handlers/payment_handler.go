package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"e-shop-modal/internal/dto"
	"e-shop-modal/internal/server"
	"e-shop-modal/internal/services"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type PaymentHandler struct {
	service       *services.PaymentService
	webhookSecret string
}

func NewPaymentHandler(s *services.PaymentService, secret string) *PaymentHandler {
	return &PaymentHandler{
		service:       s,
		webhookSecret: secret,
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
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		JSONError(c, http.StatusBadRequest, "Error leyendo body")
		return
	}

	// Validar la firma real de MP (formato ts=...,v1=...)
	xSignature := c.GetHeader("x-signature")
	xRequestID := c.GetHeader("x-request-id")
	queryID := c.Request.URL.Query().Get("data.id")
	if !validateMPSignature(xSignature, xRequestID, queryID, h.webhookSecret) {
		JSONError(c, http.StatusUnauthorized, "No autorizado")
		return
	}

	//Usa Unmarshal sobre los bytes ya leídos (no BindJSON)
	var body map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &body); err != nil {
		JSONError(c, http.StatusBadRequest, "JSON inválido")
		return
	}

	// Obtiene el payment_id y lo procesa
	data, ok := body["data"].(map[string]interface{})
	if !ok {
		JSONError(c, http.StatusBadRequest, "Payload inválido")
		return
	}
	paymentID := int(data["id"].(float64))

	err = h.service.ProcessWebhook(paymentID)
	if err != nil {
		JSONError(c, http.StatusInternalServerError, "Error al procesar pago")
		return
	}
	c.JSONResponse(http.StatusOK, "Pago exitoso")
}

// validateMPSignature implementa la validación real de MP.
// Docs: https://www.mercadopago.com.ar/developers/es/docs/your-integrations/notifications/webhooks
func validateMPSignature(xSignature, xRequestID, dataID, secret string) bool {
	if xSignature == "" {
		return false
	}

	// Parsear ts y v1 del header "ts=1234,v1=abc..."
	var ts, v1 string
	for _, part := range strings.Split(xSignature, ",") {
		kv := strings.SplitN(strings.TrimSpace(part), "=", 2)
		if len(kv) != 2 {
			continue
		}
		switch kv[0] {
		case "ts":
			ts = kv[1]
		case "v1":
			v1 = kv[1]
		}
	}
	if ts == "" || v1 == "" {
		return false
	}

	// El template que firma MP: "id:{dataID};request-id:{xRequestID};ts:{ts};"
	manifest := fmt.Sprintf("id:%s;request-id:%s;ts:%s;", dataID, xRequestID, ts)

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(manifest))
	expected := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(v1), []byte(expected))
}
