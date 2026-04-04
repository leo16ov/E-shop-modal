package services

import (
	"context"
	"e-shop-modal/internal/dto"
	"e-shop-modal/internal/repositories"

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

func (s *PaymentService) CreatePreference(ctx context.Context, product *dto.CheckoutItem) (*preference.Response, error) {
	cfg, err := config.New(s.accessToken)
	if err != nil {
		return nil, err
	}

	client := preference.NewClient(cfg)

	dataProduct, err := s.productRepo.GetByID(product.ProductID)
	if err != nil {
		return nil, err
	}

	req := preference.Request{
		Items: []preference.ItemRequest{
			{
				Title:      dataProduct.Tipo,
				Quantity:   product.Quantity,
				UnitPrice:  float64(dataProduct.Precio),
				CurrencyID: "ARS",
			},
		},
		BackURLs: &preference.BackURLsRequest{
			Success: "http://localhost:3000/success",
			Failure: "http://localhost:3000/failure",
			Pending: "http://localhost:3000/pending",
		},
		AutoReturn: "approved",

		ExternalReference: "order_123",
	}

	return client.Create(ctx, req)
}
