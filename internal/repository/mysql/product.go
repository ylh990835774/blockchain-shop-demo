package mysql

import (
	"github.com/ylh990835774/blockchain-shop-demo/internal/model"
	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) Create(product *model.Product) error {
	return r.db.Create(product).Error
}

func (r *ProductRepository) Update(product *model.Product) error {
	return r.db.Save(product).Error
}

func (r *ProductRepository) Delete(id int64) error {
	return r.db.Delete(&model.Product{}, id).Error
}

func (r *ProductRepository) GetByID(id int64) (*model.Product, error) {
	var product model.Product
	err := r.db.First(&product, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) GetByIDWithTx(tx *gorm.DB, id int64) (*model.Product, error) {
	var product model.Product
	err := tx.First(&product, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) UpdateWithTx(tx *gorm.DB, product *model.Product) error {
	return tx.Save(product).Error
}

func (r *ProductRepository) List(offset, limit int) ([]*model.Product, int64, error) {
	var products []*model.Product
	var total int64

	err := r.db.Model(&model.Product{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Offset(offset).Limit(limit).Find(&products).Error
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}
