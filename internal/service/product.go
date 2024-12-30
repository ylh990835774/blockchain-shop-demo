package service

import (
	"github.com/ylh990835774/blockchain-shop-demo/internal/model"
	"github.com/ylh990835774/blockchain-shop-demo/internal/repository/mysql"
	"github.com/ylh990835774/blockchain-shop-demo/pkg/errors"
)

type ProductService struct {
	repo *mysql.ProductRepository
}

func NewProductService(repo *mysql.ProductRepository) IProductService {
	return &ProductService{
		repo: repo,
	}
}

func (s *ProductService) Create(product *model.Product) error {
	if product == nil {
		return errors.ErrInvalidInput
	}
	return s.repo.Create(product)
}

func (s *ProductService) Update(id int64, updates map[string]interface{}) error {
	if id <= 0 {
		return errors.ErrInvalidInput
	}

	if err := s.repo.Update(id, updates); err != nil {
		if err == mysql.ErrNotFound {
			return errors.ErrNotFound
		}
		return err
	}

	return nil
}

func (s *ProductService) Delete(id int64) error {
	if id <= 0 {
		return errors.ErrInvalidInput
	}

	if err := s.repo.Delete(id); err != nil {
		if err == mysql.ErrNotFound {
			return errors.ErrNotFound
		}
		return err
	}

	return nil
}

func (s *ProductService) GetByID(id int64) (*model.Product, error) {
	if id <= 0 {
		return nil, errors.ErrInvalidInput
	}

	product, err := s.repo.GetByID(id)
	if err != nil {
		if err == mysql.ErrNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, err
	}

	return product, nil
}

func (s *ProductService) List(page, pageSize int) ([]*model.Product, int64, error) {
	if page <= 0 || pageSize <= 0 {
		return nil, 0, errors.ErrInvalidInput
	}

	offset := (page - 1) * pageSize
	return s.repo.List(offset, pageSize)
}
