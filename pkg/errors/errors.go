package errors

import "errors"

var (
	ErrNotFound          = errors.New("记录不存在")
	ErrInvalidInput      = errors.New("无效的输入")
	ErrDuplicateEntry    = errors.New("记录已存在")
	ErrInvalidToken      = errors.New("无效的令牌")
	ErrUnauthorized      = errors.New("未授权的访问")
	ErrInsufficientStock = errors.New("库存不足")
	ErrForbidden         = errors.New("禁止访问")
	ErrBadRequest        = errors.New("请求参数错误")
	ErrNoFieldsToUpdate  = errors.New("至少需要更新一个字段")
)
