package product_sale_yearly_report

import (
	"fmt"
	"hacktiv8-techrawih-go-product-sale/internal/module/product"
	"log"

	"gorm.io/gorm"
)

// Repository Interface
type Repository interface {
	DeleteAll() error
}

// NewRepository will implement ProductRepository Interface
func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

type repository struct {
	db *gorm.DB
}

func (r *repository) DeleteAll() (err error) {
	var yearlyReports []ProductSaleYearlyReport
	if err = r.db.Find(&yearlyReports).Error; err != nil {
		return fmt.Errorf("failed to load product sale yearly reports: %w", err)
	}

	productUpdates := make(map[uint]int)
	for _, report := range yearlyReports {
		productUpdates[report.ProductID] += report.CountSales
	}

	for productID, countSales := range productUpdates {
		if err = r.db.Model(&product.Product{}).Where("id = ?", productID).UpdateColumn("stock", gorm.Expr("stock + ?", countSales)).Error; err != nil {
			return fmt.Errorf("failed to update stock for product ID %d: %w", productID, err)
		}
	}

	if err = r.db.Where("1=1").Unscoped().Delete(&ProductSaleYearlyReport{}).Error; err != nil {
		log.Fatalf("failed to delete records: %v", err)
		return err
	}

	return nil
}
