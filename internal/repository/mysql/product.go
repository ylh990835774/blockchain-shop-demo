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

func (r *ProductRepository) Update(id int64, updates map[string]interface{}) error {
	result := r.db.Model(&model.Product{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *ProductRepository) UpdateStock(id int64, quantity int) error {
	result := r.db.Model(&model.Product{}).Where("id = ? AND stock >= ?", id, -quantity).
		UpdateColumn("stock", gorm.Expr("stock + ?", quantity))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *ProductRepository) UpdateStockWithTx(tx *gorm.DB, id int64, quantity int) error {
	result := tx.Model(&model.Product{}).Where("id = ? AND stock >= ?", id, -quantity).
		UpdateColumn("stock", gorm.Expr("stock + ?", quantity))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *ProductRepository) Delete(id int64) error {
	result := r.db.Delete(&model.Product{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
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
