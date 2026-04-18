package services

import (
	"e-shop-modal/internal/dto"
	"e-shop-modal/internal/models"
	"e-shop-modal/internal/repositories"
	"e-shop-modal/internal/server"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/mercadopago/sdk-go/pkg/config"
	"github.com/mercadopago/sdk-go/pkg/preference"
)

type PaymentService struct {
	accessToken string
	productRepo *repositories.ProductRepository
	orderRepo   *repositories.OrderRepository
	httpClient  *http.Client
}

func NewPaymentService(token string, p *repositories.ProductRepository, o *repositories.OrderRepository) *PaymentService {
	return &PaymentService{
		accessToken: token,
		productRepo: p,
		orderRepo:   o,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *PaymentService) CreatePreference(c *server.Context, item *dto.CheckoutItem) (*preference.Response, error) {
	// Config Mercado Pago
	cfg, err := config.New(s.accessToken)
	if err != nil {
		return nil, err
	}
	client := preference.NewClient(cfg)

	// Busca producto en la DB para saber el nombre y el precio
	dataProduct, err := s.productRepo.GetByID(c, item.ProductID)
	if err != nil {
		return nil, err
	}

	// Calcula el precio total
	total := float64(dataProduct.Precio) * float64(item.Quantity)

	// Crea orden en DB
	order, err := s.orderRepo.Create(c, total)
	if err != nil {
		return nil, err
	}

	// Usar ID como external_reference y guardarlo en DB
	externalRef := fmt.Sprintf("%d", order.ID)

	err = s.orderRepo.SetExternalReference(c, order.ID, externalRef)
	if err != nil {
		return nil, err
	}
	req := preference.Request{
		Items: []preference.ItemRequest{
			{
				Title:      dataProduct.Tipo,
				Quantity:   item.Quantity,
				UnitPrice:  float64(dataProduct.Precio),
				CurrencyID: "ARS",
			},
		},
		/*BackURLs: &preference.BackURLsRequest{
			Success: "http://localhost:3000/success",
			Failure: "http://localhost:3000/failure",
			Pending: "http://localhost:3000/pending",
		},
		AutoReturn:        "approved",
		NotificationURL: "https://e-shop-modal.onrender.com/webhook",*/
		ExternalReference: externalRef,
	}
	return client.Create(c.Context(), req)
}

func (s *PaymentService) ProcessWebhook(c *server.Context, paymentID int64) error {

	// Obtiene el pago real desde MP
	payment, err := s.GetPayment(c, paymentID)
	if err != nil {
		fmt.Printf("Error 1")
		return err
	}

	// Busca orden
	order, err := s.orderRepo.GetByExternalReference(c, payment.ExternalReference)
	if err != nil {
		fmt.Printf("Error 2")
		return err
	}

	// Valida el monto
	if payment.TransactionAmount != order.Total {
		fmt.Printf("Error 3")
		return fmt.Errorf("monto inválido")
	}

	// Actualiza el estado
	return s.orderRepo.UpdateStatus(c, order.ID, payment.Status)
}

func (s *PaymentService) GetPayment(c *server.Context, paymentID int64) (*models.PaymentInfo, error) {

	url := fmt.Sprintf("https://api.mercadopago.com/v1/payments/%d", paymentID)

	req, err := http.NewRequestWithContext(c.Context(), "GET", url, nil)
	if err != nil {
		fmt.Printf("Error GetPayment 1")
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+s.accessToken)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		fmt.Printf("Error GetPayment 2")
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Status            string  `json:"status"`
		TransactionAmount float64 `json:"transaction_amount"`
		ExternalReference string  `json:"external_reference"`
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error GetPayment 3")
		return nil, fmt.Errorf("MP respondió %d", resp.StatusCode)
	}
	fmt.Println(resp.StatusCode)
	fmt.Println(resp)

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		fmt.Printf("Error GetPayment 4")
		return nil, fmt.Errorf("error parseando respuesta MP: %w", err)
	}

	return &models.PaymentInfo{
		Status:            result.Status,
		TransactionAmount: result.TransactionAmount,
		ExternalReference: result.ExternalReference,
	}, nil
}
