package mysql

import "gorm.io/gorm"

type Transaction interface {
    Begin() *gorm.DB
    Commit() error
    Rollback() error
}

type GormTransaction struct {
    tx *gorm.DB
}

func NewTransaction(db *gorm.DB) *GormTransaction {
    return &GormTransaction{
        tx: db.Begin(),
    }
}

func (t *GormTransaction) Begin() *gorm.DB {
    return t.tx
}

func (t *GormTransaction) Commit() error {
    return t.tx.Commit().Error
}

func (t *GormTransaction) Rollback() error {
    return t.tx.Rollback().Error
}