package model

import (
	"time"
)

type Product struct {
	ID            string            `json:"id" bson:"_id,omitempty"`
	Name          string            `json:"name" bson:"name"`
	Description   string            `json:"description" bson:"description"`
	SKU           string            `json:"sku" bson:"sku"`
	Price         Money             `json:"price" bson:"price"`
	StockQuantity int32             `json:"stock_quantity" bson:"stock_quantity"`
	Category      string            `json:"category" bson:"category"`
	Attributes    map[string]string `json:"attributes" bson:"attributes"`
	CreatedAt     time.Time         `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at" bson:"updated_at"`
}

type Money struct {
	Currency string `json:"currency" bson:"currency"`
	Amount   int64  `json:"amount" bson:"amount"`
}

// Request/Response models
type CreateProductRequest struct {
	Name          string            `json:"name" binding:"required"`
	Description   string            `json:"description"`
	SKU           string            `json:"sku" binding:"required"`
	Price         Money             `json:"price" binding:"required"`
	StockQuantity int32             `json:"stock_quantity" binding:"required"`
	Category      string            `json:"category" binding:"required"`
	Attributes    map[string]string `json:"attributes"`
}

type UpdateProductRequest struct {
	ID            string            `json:"-"`
	Name          string            `json:"name"`
	Description   string            `json:"description"`
	Price         Money             `json:"price"`
	StockQuantity int32             `json:"stock_quantity"`
	Category      string            `json:"category"`
	Attributes    map[string]string `json:"attributes"`
}

type StockItem struct {
	ProductID string `json:"product_id" binding:"required"`
	Quantity  int32  `json:"quantity" binding:"required"`
}

type ReserveStockRequest struct {
	OrderID string      `json:"order_id" binding:"required"`
	Items   []StockItem `json:"items" binding:"required"`
}

type ReleaseStockRequest struct {
	OrderID string      `json:"order_id" binding:"required"`
	Items   []StockItem `json:"items" binding:"required"`
}

type UnavailableItem struct {
	OrderID string      `json:"order_id" binding:"required"`
	Items   []StockItem `json:"items" binding:"required"`
}
