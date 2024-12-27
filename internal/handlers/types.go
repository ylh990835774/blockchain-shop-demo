package handlers

// 用户相关请求结构体
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Password string `json:"password" binding:"required,min=6,max=32"`
	Phone    string `json:"phone" binding:"required,len=11"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 商品相关请求结构体
type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Stock       int     `json:"stock" binding:"required,gte=0"`
}

// 订单相关请求结构体
type CreateOrderRequest struct {
	ProductID int64 `json:"product_id" binding:"required"`
	Quantity  int   `json:"quantity" binding:"required,gt=0"`
}

// 用户信息更新请求结构体
type UpdateProfileRequest struct {
	Phone   string `json:"phone" binding:"required,len=11"`
	Address string `json:"address" binding:"required"`
}
