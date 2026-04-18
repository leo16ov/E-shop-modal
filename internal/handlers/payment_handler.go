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
	"log"
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
	resp, err := h.service.CreatePreference(c, req)
	fmt.Println(resp.ID)
	if err != nil {
		fmt.Printf("ERROR CreatePreference: %v\n", err.Error())
		c.JSONResponse(http.StatusInternalServerError, "Error creando preferencia")
		return
	}

	// Devolver init_point
	c.JSONResponse(http.StatusOK, map[string]interface{}{
		"prefence_id":  resp.ID,
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

	// Deserializa a struct tipado (sin castings manuales)
	var payload *dto.Payload
	if err := json.Unmarshal(bodyBytes, &payload); err != nil {
		JSONError(c, http.StatusBadRequest, "JSON inválido")
		fmt.Println("JSON ignorada")
		return
	}
	// Filtra eventos que no son de pago
	if payload.Type != "payment" {
		c.JSONResponse(http.StatusOK, "evento ignorado")
		fmt.Println("Evento ignorada")
		return
	}
	if payload.Action != "payment.created" && payload.Action != "payment.updated" {
		c.JSONResponse(http.StatusOK, "accion ignorada")
		fmt.Println("Accion ignorada")
		return
	}

	// Proceso de errores internos no los propagamos a MP
	if err := h.service.ProcessWebhook(c, payload.Data.ID); err != nil {
		log.Printf("error procesando webhook payment_id=%d: %v", payload.Data.ID, err)
		c.JSONResponse(http.StatusOK, "ok")
		fmt.Printf("\nError procesando webhook payment_id=%d: %v\n", payload.Data.ID, err)
		return
	}
	fmt.Println("Pago procesado")
	c.JSONResponse(http.StatusOK, "Pago procesado")
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
