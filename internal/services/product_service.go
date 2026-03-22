package services

import (
	"e-shop-modal/internal/models"
	"e-shop-modal/internal/store"
	"errors"
)

//type Logger interface {
//Log(msg, error string)
//}

type Service struct {
	store store.Store
	//logger Logger
}

func New(s store.Store) *Service {
	return &Service{
		store: s,
		//logger: nil,
	}
}

func (s *Service) ObtenerTodosLosProducts() ([]*models.Product, error) {
	//s.logger.Log("Estamos obteniendo productos", "")

	products, err := s.store.GetAll()
	if err != nil {
		//s.logger.Log("Error: %s", err.Error())
		return nil, err
	}
	return products, nil
}

func (s *Service) ObtenerProductPorID(id int) (*models.Product, error) {
	return s.store.GetByID(id)
}

func (s *Service) SubirProduct(product models.Product) (*models.Product, error) {
	if product.Tipo == "" {
		return nil, errors.New("El producto no tiene nombre")
	}
	return s.store.Create(&product)
}

func (s *Service) ModificarProduct(id int, product models.Product) (*models.Product, error) {
	if product.Tipo == "" {
		return nil, errors.New("El producto no tiene nombre")
	}
	return s.store.Update(id, &product)
}

func (s *Service) QuitarProduct(id int) error {
	return s.store.Delete(id)

}
