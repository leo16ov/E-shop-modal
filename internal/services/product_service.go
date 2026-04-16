package services

import (
	"e-shop-modal/internal/models"
	"e-shop-modal/internal/repositories"
	"e-shop-modal/internal/server"
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

func (s *ProductService) ObtenerTodosLosProducts(c *server.Context) ([]*models.Product, error) {
	//s.logger.Log("Estamos obteniendo productos", "")

	products, err := s.repository.GetAll(c)
	if err != nil {
		//s.logger.Log("Error: %s", err.Error())
		return nil, err
	}
	return products, nil
}

func (s *ProductService) ObtenerProductPorID(c *server.Context, id int) (*models.Product, error) {
	return s.repository.GetByID(c, id)
}

func (s *ProductService) SubirProduct(c *server.Context, product models.Product) (*models.Product, error) {
	if product.Tipo == "" {
		return nil, errors.New("El producto no tiene nombre")
	}
	return s.repository.Create(c, &product)
}

func (s *ProductService) ModificarProduct(c *server.Context, id int, product models.Product) (*models.Product, error) {
	if product.Tipo == "" {
		return nil, errors.New("El producto no tiene nombre")
	}
	return s.repository.Update(c, id, &product)
}

func (s *ProductService) QuitarProduct(c *server.Context, id int) error {
	return s.repository.Delete(c, id)

}
