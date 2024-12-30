package service

import (
	"github.com/ylh990835774/blockchain-shop-demo/internal/model"
	"github.com/ylh990835774/blockchain-shop-demo/internal/repository/mysql"
)

type ProductService struct {
	repo *mysql.ProductRepository
}

func NewProductService(repo *mysql.ProductRepository) IProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) Create(product *model.Product) error {
	return s.repo.Create(product)
}

func (s *ProductService) Update(id int64, updates map[string]interface{}) error {
	product, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	for key, value := range updates {
		switch key {
		case "name":
			if name, ok := value.(string); ok {
				product.Name = name
			}
		case "description":
			if desc, ok := value.(string); ok {
				product.Description = desc
			}
		case "price":
			if price, ok := value.(float64); ok {
				product.Price = price
			}
		case "stock":
			if stock, ok := value.(float64); ok {
				product.Stock = int(stock)
			}
		}
	}

	return s.repo.Update(product)
}

func (s *ProductService) Delete(id int64) error {
	return s.repo.Delete(id)
}

func (s *ProductService) Get(id int64) (*model.Product, error) {
	return s.repo.GetByID(id)
}

func (s *ProductService) List(page, pageSize int) ([]*model.Product, int64, error) {
	offset := (page - 1) * pageSize
	return s.repo.List(offset, pageSize)
}
