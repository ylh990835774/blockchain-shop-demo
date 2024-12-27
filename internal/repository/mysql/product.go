package mysql

import (
	"github.com/ylh990835774/blockchain-shop-demo/internal/model"

	"gorm.io/gorm"
)

type ProductRepository struct {
	BaseRepository
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *ProductRepository) Create(product *model.Product) error {
	return r.db.Create(product).Error
}

func (r *ProductRepository) GetByID(id int64) (*model.Product, error) {
	var product model.Product
	err := r.db.First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) List(page, pageSize int) ([]*model.Product, int64, error) {
	var products []*model.Product
	var total int64

	offset := (page - 1) * pageSize

	err := r.db.Model(&model.Product{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Offset(offset).Limit(pageSize).Find(&products).Error
	return products, total, err
}

func (r *ProductRepository) Update(product *model.Product) error {
	return r.db.Save(product).Error
}

func (r *ProductRepository) UpdateWithTx(tx *gorm.DB, product *model.Product) error {
	return tx.Save(product).Error
}

func (r *ProductRepository) Delete(id int64) error {
	return r.db.Delete(&model.Product{}, id).Error
}
