package service

import "errors"

var (
	// ErrInsufficientStock 库存不足
	ErrInsufficientStock = errors.New("insufficient stock")
	// ErrTransactionNotFound 交易记录未找到
	ErrTransactionNotFound = errors.New("transaction not found")
)
