package mysql

import (
	"blockchain-shop/internal/model"

	"gorm.io/gorm"
)

type OrderRepository struct {
	BaseRepository
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *OrderRepository) Create(order *model.Order) error {
	return r.db.Create(order).Error
}

func (r *OrderRepository) CreateWithTx(tx *gorm.DB, order *model.Order) error {
	return tx.Create(order).Error
}

func (r *OrderRepository) GetByID(id int64) (*model.Order, error) {
	var order model.Order
	err := r.db.First(&order, id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) ListByUserID(userID int64, page, pageSize int) ([]*model.Order, int64, error) {
	var orders []*model.Order
	var total int64

	offset := (page - 1) * pageSize

	err := r.db.Model(&model.Order{}).Where("user_id = ?", userID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Where("user_id = ?", userID).Offset(offset).Limit(pageSize).Find(&orders).Error
	return orders, total, err
}

func (r *OrderRepository) BeginTx() *gorm.DB {
	return r.db.Begin()
}

func (r *OrderRepository) Update(order *model.Order) error {
	return r.db.Save(order).Error
}
