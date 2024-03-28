package sales

import (
	"context"
	"hacktiv8-techrawih-go-product-sale/internal/module/product"
	"log"

	"gorm.io/gorm"
)

// Repository Interface
type Repository interface {
	Save(ctx context.Context, request Sales) (*Sales, error)
	SaveAll(ctx context.Context, request []*Sales) error
	DeleteAll() error
	GetAll() ([]*Sales, error)
	GetProductByName(name string) (product.Product, error)
}

// NewRepository will implement SalesRepository Interface
func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

type repository struct {
	db *gorm.DB
}

func (r *repository) Save(ctx context.Context, request Sales) (res *Sales, err error) {
	dtoSale := Sales{ProductID: uint(request.ProductID), QtySold: request.QtySold, SaleAt: request.SaleAt}
	err = r.db.Save(&dtoSale).Error
	if err != nil {
		return nil, err
	}

	return &dtoSale, nil
}

func (r *repository) GetAll() ([]*Sales, error) {
	var salesDatas []*Sales

	if err := r.db.Find(&salesDatas).Error; err != nil {
		log.Fatalf("failed to get all records: %v", err)
		return nil, err
	}
	return salesDatas, nil
}

func (r *repository) DeleteAll() (err error) {
	if err = r.db.Where("1=1").Unscoped().Delete(&Sales{}).Error; err != nil {
		log.Fatalf("failed to delete records: %v", err)
		return err
	}
	return nil
}

func (r *repository) SaveAll(ctx context.Context, salesDto []*Sales) (err error) {
	// Start a transaction
	tx := r.db.Begin()

	// Check if the transaction was successful
	if tx.Error != nil {
		return tx.Error
	}

	batchSize := 500
	for i := 0; i < len(salesDto); i += batchSize {
		end := i + batchSize
		if end > len(salesDto) {
			end = len(salesDto)
		}

		err = tx.CreateInBatches((salesDto)[i:end], batchSize).Error
		if err != nil {
			// If there is an error, rollback the transaction
			tx.Rollback()
			return err
		}
	}

	// Commit the transaction if the insert is successful
	return tx.Commit().Error
}

func (r *repository) GetProductByName(name string) (res product.Product, err error) {
	query := r.db

	if name != "" {
		query = query.Where("name ILIKE ?", "%"+name+"%")
	}

	result := query.First(&res)

	return res, result.Error
}
