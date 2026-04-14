package services

import (
	"context"
	"e-shop-modal/internal/dto"
	"e-shop-modal/internal/models"
	"e-shop-modal/internal/repositories"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mercadopago/sdk-go/pkg/config"
	"github.com/mercadopago/sdk-go/pkg/preference"
)

type PaymentService struct {
	accessToken string
	productRepo *repositories.ProductRepository
	orderRepo   *repositories.OrderRepository
}

func NewPaymentService(token string, p *repositories.ProductRepository, o *repositories.OrderRepository) *PaymentService {
	return &PaymentService{
		accessToken: token,
		productRepo: p,
		orderRepo:   o,
	}
}

func (s *PaymentService) CreatePreference(ctx context.Context, item *dto.CheckoutItem) (*preference.Response, error) {
	// Config Mercado Pago
	cfg, err := config.New(s.accessToken)
	if err != nil {
		return nil, err
	}
	client := preference.NewClient(cfg)

	// Busca producto en la DB para saber el nombre y el precio
	dataProduct, err := s.productRepo.GetByID(item.ProductID)
	if err != nil {
		return nil, err
	}

	// Calcula el precio total
	total := float64(dataProduct.Precio) * float64(item.Quantity)

	// Crea orden en DB
	order, err := s.orderRepo.Create(total)
	if err != nil {
		return nil, err
	}

	// Usar ID como external_reference y guardarlo en DB
	externalRef := fmt.Sprintf("%d", order.ID)

	err = s.orderRepo.SetExternalReference(order.ID, externalRef)
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
		ExternalReference: externalRef,*/
	}
	return client.Create(ctx, req)
}

func (s *PaymentService) ProcessWebhook(paymentID int64) error {

	// Obtiene el pago real desde MP
	payment, err := s.GetPayment(paymentID)
	if err != nil {
		return err
	}

	// Busca orden
	order, err := s.orderRepo.GetByExternalReference(
		fmt.Sprintf("%s", payment.ExternalReference),
	)
	if err != nil {
		return err
	}

	// Valida el monto
	if payment.TransactionAmount != order.Total {
		return fmt.Errorf("monto inválido")
	}

	// Actualiza el estado
	return s.orderRepo.UpdateStatus(order.ID, payment.Status)
}

func (s *PaymentService) GetPayment(paymentID int64) (*models.PaymentInfo, error) {

	url := fmt.Sprintf("https://api.mercadopago.com/v1/payments/%d", paymentID)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Bearer "+s.accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Status            string  `json:"status"`
		TransactionAmount float64 `json:"transaction_amount"`
		ExternalReference string  `json:"external_reference"`
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("MP respondió %d", resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("error parseando respuesta MP: %w", err)
	}

	return &models.PaymentInfo{
		Status:            result.Status,
		TransactionAmount: result.TransactionAmount,
		ExternalReference: result.ExternalReference,
	}, nil
}
