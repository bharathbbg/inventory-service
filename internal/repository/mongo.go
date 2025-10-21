package repository

import (
	"context"
	"time"

	"github.com/bharathbbg/inventory-service/internal/config"
	"github.com/bharathbbg/inventory-service/internal/model"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepository struct {
	client       *mongo.Client
	db           *mongo.Database
	products     *mongo.Collection
	reservations *mongo.Collection
}

func (r *MongoRepository) ReleaseStock(ctx context.Context, req *model.ReleaseStockRequest) (bool, error) {
	panic("unimplemented")
}

func (r *MongoRepository) ReserveStock(ctx context.Context, req *model.ReserveStockRequest) (bool, string, any, error) {
	panic("unimplemented")
}

func (r *MongoRepository) CheckStock(ctx context.Context, items []model.StockItem) (bool, any, error) {
	panic("unimplemented")
}

func (r *MongoRepository) DeleteProduct(ctx context.Context, id string) (any, any) {
	panic("unimplemented")
}

func (r *MongoRepository) ListProducts(ctx context.Context, category string, page int, pageSize int) ([]*model.Product, int, error) {
	panic("unimplemented")
}

func (r *MongoRepository) UpdateProduct(ctx context.Context, product *model.Product) (any, error) {
	panic("unimplemented")
}

func NewMongoRepository(config config.MongoDBConfig) (*MongoRepository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.URI))
	if err != nil {
		return nil, err
	}

	// Ping the database to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	db := client.Database(config.Database)
	products := db.Collection("products")
	reservations := db.Collection("reservations")

	return &MongoRepository{
		client:       client,
		db:           db,
		products:     products,
		reservations: reservations,
	}, nil
}

func (r *MongoRepository) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return r.client.Disconnect(ctx)
}

func (r *MongoRepository) CreateProduct(ctx context.Context, product *model.Product) (*model.Product, error) {
	product.ID = uuid.New().String()
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()

	_, err := r.products.InsertOne(ctx, product)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (r *MongoRepository) GetProduct(ctx context.Context, id string) (*model.Product, error) {
	var product model.Product
	err := r.products.FindOne(ctx, bson.M{"_id": id}).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &product, nil
}

// Additional repository methods would be implemented here
