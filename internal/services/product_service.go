package services

import (
	"e-shop-modal/internal/models"
	"e-shop-modal/internal/repositories"
	"errors"
)

//type Logger interface {
//Log(msg, error string)
//}

type ProductService struct {
	repository *repositories.ProductRepository
	//logger Logger
}

func NewProductService(r *repositories.ProductRepository) *ProductService {
	return &ProductService{
		repository: r,
		//logger: nil,
	}
}

func (s *ProductService) ObtenerTodosLosProducts() ([]*models.Product, error) {
	//s.logger.Log("Estamos obteniendo productos", "")

	products, err := s.repository.GetAll()
	if err != nil {
		//s.logger.Log("Error: %s", err.Error())
		return nil, err
	}
	return products, nil
}

func (s *ProductService) ObtenerProductPorID(id int) (*models.Product, error) {
	return s.repository.GetByID(id)
}

func (s *ProductService) SubirProduct(product models.Product) (*models.Product, error) {
	if product.Tipo == "" {
		return nil, errors.New("El producto no tiene nombre")
	}
	return s.repository.Create(&product)
}

func (s *ProductService) ModificarProduct(id int, product models.Product) (*models.Product, error) {
	if product.Tipo == "" {
		return nil, errors.New("El producto no tiene nombre")
	}
	return s.repository.Update(id, &product)
}

func (s *ProductService) QuitarProduct(id int) error {
	return s.repository.Delete(id)

}
