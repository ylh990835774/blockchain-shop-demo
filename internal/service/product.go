package service

import (
	"github.com/ylh990835774/blockchain-shop-demo/internal/model"
	"github.com/ylh990835774/blockchain-shop-demo/internal/repository/mysql"
	"github.com/ylh990835774/blockchain-shop-demo/pkg/errors"

	"gorm.io/gorm"
)

type ProductService struct {
	repo *mysql.ProductRepository
}

func NewProductService(repo *mysql.ProductRepository) *ProductService {
	return &ProductService{
		repo: repo,
	}
}

func (s *ProductService) Create(product *model.Product) error {
	if product.Name == "" || product.Price <= 0 {
		return errors.ErrInvalidInput
	}
	return s.repo.Create(product)
}

func (s *ProductService) GetByID(id int64) (*model.Product, error) {
	product, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.ErrNotFound
	}
	return product, nil
}

func (s *ProductService) List(page, pageSize int) ([]*model.Product, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	return s.repo.List(page, pageSize)
}

func (s *ProductService) Update(product *model.Product) error {
	if product.ID <= 0 {
		return errors.ErrInvalidInput
	}

	// 检查商品是否存在
	existing, err := s.repo.GetByID(product.ID)
	if err != nil {
		return errors.ErrNotFound
	}

	// 保留创建时间
	product.CreatedAt = existing.CreatedAt

	return s.repo.Update(product)
}

func (s *ProductService) Delete(id int64) error {
	if id <= 0 {
		return errors.ErrInvalidInput
	}

	// 检查商品是否存在
	if _, err := s.repo.GetByID(id); err != nil {
		return errors.ErrNotFound
	}

	return s.repo.Delete(id)
}

func (s *ProductService) UpdateWithTx(tx *gorm.DB, product *model.Product) error {
	if product.ID <= 0 {
		return errors.ErrInvalidInput
	}

	if _, err := s.GetByID(product.ID); err != nil {
		if err == mysql.ErrNotFound {
			return errors.ErrNotFound
		}
		return err
	}

	return s.repo.UpdateWithTx(tx, product)
}
