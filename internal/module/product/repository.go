package product

import (
	"context"
	"gorm.io/gorm"
)

// Repository Interface
type Repository interface {
	Insert(ctx context.Context, product Product) (Product, error)
	GetIdByName(name string) (Product, error)
}

// NewRepository will implement ProductRepository Interface
func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

type repository struct {
	db *gorm.DB
}

func (r *repository) Insert(ctx context.Context, product Product) (res Product, err error) {
	result := r.db.Save(&product)
	return product, result.Error
}

func (r *repository) GetIdByName(name string) (res Product, err error) {
	query := r.db

	if name != "" {
		query = query.Where("title ILIKE ?", "%"+name+"%")
	}

	result := query.First(&res)

	return res, result.Error
}
