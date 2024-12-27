package mysql

import (
	"fmt"

	"github.com/ylh990835774/blockchain-shop-demo/internal/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewDB(conf *Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.Username,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.Database,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// 自动迁移表结构
	err = db.AutoMigrate(
		&model.User{},
		&model.Product{},
		&model.Order{},
	)

	return db, err
}

type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
}
