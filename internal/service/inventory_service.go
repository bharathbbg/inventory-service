package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bharathbbg/inventory-service/internal/model"
	"github.com/bharathbbg/inventory-service/internal/repository"
)

type InventoryService struct {
	repo  *repository.MongoRepository
	cache *repository.RedisCache
}

func NewInventoryService(repo *repository.MongoRepository, cache *repository.RedisCache) *InventoryService {
	return &InventoryService{
		repo:  repo,
		cache: cache,
	}
}

func (s *InventoryService) CreateProduct(ctx context.Context, req *model.CreateProductRequest) (*model.Product, error) {
	// Validate request
	if req.Name == "" {
		return nil, errors.New("product name is required")
	}
	if req.SKU == "" {
		return nil, errors.New("product SKU is required")
	}

	// Create product
	product := &model.Product{
		Name:          req.Name,
		Description:   req.Description,
		SKU:           req.SKU,
		Price:         req.Price,
		StockQuantity: req.StockQuantity,
		Category:      req.Category,
		Attributes:    req.Attributes,
	}

	// Save to database
	savedProduct, err := s.repo.CreateProduct(ctx, product)
	if err != nil {
		return nil, err
	}

	// Cache the result
	if err := s.cache.CacheProduct(ctx, savedProduct); err != nil {
		// Just log error, don't fail the request
		// log.Printf("Failed to cache product: %v", err)
	}

	return savedProduct, nil
}

func (s *InventoryService) GetProduct(ctx context.Context, id string) (*model.Product, error) {
	// Try to get from cache first
	cachedProduct, err := s.cache.GetCachedProduct(ctx, id)
	if err == nil && cachedProduct != nil {
		return cachedProduct, nil
	}

	// If not in cache, get from database
	product, err := s.repo.GetProduct(ctx, id)
	if err != nil {
		return nil, err
	}

	// Cache the result for future requests
	if product != nil {
		if err := s.cache.CacheProduct(ctx, product); err != nil {
			// log.Printf("Failed to cache product: %v", err)
		}
	}

	return product, nil
}

func (s *InventoryService) UpdateProduct(ctx context.Context, req *model.UpdateProductRequest) (*model.Product, error) {
	// Get existing product
	product, err := s.GetProduct(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, errors.New("product not found")
	}

	// Update fields if provided
	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.Price.Amount > 0 {
		product.Price = req.Price
	}
	if req.StockQuantity >= 0 {
		product.StockQuantity = req.StockQuantity
	}
	if req.Category != "" {
		product.Category = req.Category
	}
	if req.Attributes != nil {
		product.Attributes = req.Attributes
	}

	product.UpdatedAt = time.Now()

	// Save to database
	updatedAny, err := s.repo.UpdateProduct(ctx, product)
	if err != nil {
		return nil, err
	}

	updatedProduct, ok := updatedAny.(*model.Product)
	if !ok || updatedProduct == nil {
		return nil, errors.New("repository returned unexpected type for updated product")
	}

	// Update cache
	if err := s.cache.CacheProduct(ctx, updatedProduct); err != nil {
		// log.Printf("Failed to update product cache: %v", err)
	}

	return updatedProduct, nil
}

func (s *InventoryService) ListProducts(ctx context.Context, category string, page, pageSize int) ([]*model.Product, int, error) {
	// Ensure valid pagination
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// Get from database
	return s.repo.ListProducts(ctx, category, page, pageSize)
}
func (s *InventoryService) DeleteProduct(ctx context.Context, id string) (bool, error) {
	result, err := s.repo.DeleteProduct(ctx, id)
	if err != nil {
		// err may be of type any (interface{}); convert to error if possible
		if e, ok := err.(error); ok {
			return false, e
		}
		return false, fmt.Errorf("%v", err)
	}

	success, ok := result.(bool)
	if !ok {
		return false, errors.New("repository returned unexpected type for delete result")
	}

	// If successfully deleted, remove from cache
	if success {
		// We can ignore cache errors here
		s.cache.DeleteCachedProduct(ctx, id)
	}

	return success, nil
}

func (s *InventoryService) CheckStock(ctx context.Context, items []model.StockItem) (bool, []model.UnavailableItem, error) {
	ok, anyUnavailable, err := s.repo.CheckStock(ctx, items)
	if err != nil {
		return false, nil, err
	}

	if anyUnavailable == nil {
		return ok, nil, nil
	}

	unavailable, castOK := anyUnavailable.([]model.UnavailableItem)
	if !castOK {
		return false, nil, errors.New("repository returned unexpected type for unavailable items")
	}

	return ok, unavailable, nil
}

func (s *InventoryService) ReserveStock(ctx context.Context, req *model.ReserveStockRequest) (bool, string, []model.UnavailableItem, error) {
	ok, reservationID, anyUnavailable, err := s.repo.ReserveStock(ctx, req)
	if err != nil {
		return false, "", nil, err
	}

	if anyUnavailable == nil {
		return ok, reservationID, nil, nil
	}

	unavailable, castOK := anyUnavailable.([]model.UnavailableItem)
	if !castOK {
		return false, "", nil, errors.New("repository returned unexpected type for unavailable items")
	}

	return ok, reservationID, unavailable, nil
}

func (s *InventoryService) ReleaseStock(ctx context.Context, req *model.ReleaseStockRequest) (bool, error) {
	return s.repo.ReleaseStock(ctx, req)
}
