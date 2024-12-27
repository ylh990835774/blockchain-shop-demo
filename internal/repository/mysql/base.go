package mysql

import "gorm.io/gorm"

type Repository interface {
	DB() *gorm.DB
}

type BaseRepository struct {
	db *gorm.DB
}

func (r *BaseRepository) DB() *gorm.DB {
	return r.db
}

func NewBaseRepository(db *gorm.DB) BaseRepository {
	return BaseRepository{db: db}
}
