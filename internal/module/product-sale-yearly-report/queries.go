package product_sale_yearly_report

import (
	"hacktiv8-techrawih-go-product-sale/internal/module/product"
	"time"

	"gorm.io/gorm"
)

func AggregateSalesByProduct(db *gorm.DB) (*ProductSaleYearlyReport, error) {
	var results []struct {
		ProductID    uint    `gorm:"column:product_id"`
		QuantitySold float64 `gorm:"column:total_quantity"`
		SellingPrice float64 `gorm:"column:selling_price"`
		BuyingPrice  float64 `gorm:"column:buying_price"`
		SaleAt       int     `gorm:"column:sale_year"`
	}

	var productSaleYearlyReport ProductSaleYearlyReport
	db.Table("sales").
		Select("sales.product_id, EXTRACT(YEAR FROM sales.sale_at) AS sale_year, SUM(sales.qty_sold) as total_quantity, products.selling_price, products.buying_price").
		Joins("left join products on products.id = sales.product_id").
		Group("sales.product_id, EXTRACT(YEAR FROM sales.sale_at), products.selling_price, products.buying_price").
		Scan(&results)

	for _, result := range results {

		totalGrossSales := result.QuantitySold * result.SellingPrice
		totalNettSales := totalGrossSales - (result.QuantitySold * result.BuyingPrice)

		productSaleYearlyReport = ProductSaleYearlyReport{
			ProductID:       result.ProductID,
			SellingPrice:    result.SellingPrice,
			BuyingPrice:     result.BuyingPrice,
			TotalGrossSales: totalGrossSales,
			TotalNettSales:  totalNettSales,
			CountSales:      int(result.QuantitySold),
			Year:            result.SaleAt,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}

		err := db.Create(&productSaleYearlyReport).Error
		if err != nil {
			return nil, err
		}
	}

	return &productSaleYearlyReport, nil
}

func UpdateStockProduct(db *gorm.DB) ([]*ProductSaleYearlyReport, error) {
	var productSaleYearlyReport []*ProductSaleYearlyReport
	err := db.Find(&productSaleYearlyReport).Error
	if err != nil {
		return nil, err
	}

	for _, yearlyReport := range productSaleYearlyReport {
		var product product.Product
		err := db.First(&product, yearlyReport.ProductID).Error
		if err != nil {
			return nil, err
		}

		finalStock := product.Stock - yearlyReport.CountSales
		err = db.Model(&product).Where("id = ?", yearlyReport.ProductID).Update("stock", finalStock).Error
		if err != nil {
			return nil, err
		}
	}

	return productSaleYearlyReport, nil
}
