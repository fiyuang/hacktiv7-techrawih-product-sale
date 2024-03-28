package router

import (
	product_trx "hacktiv8-techrawih-go-product-sale/internal/module/product-sale-yearly-report"
	sales2 "hacktiv8-techrawih-go-product-sale/internal/module/sales"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterAPIService(e *gin.Engine, db *gorm.DB) {
	basePath := "api/v1"
	registerSalesAPIService(e, db, basePath)
}

func registerSalesAPIService(r *gin.Engine, db *gorm.DB, basePath string) {
	// Initialize Sales Service
	salesRepo := sales2.NewRepository(db)
	productTrxRepo := product_trx.NewRepository(db)
	salesService := sales2.NewService(salesRepo, productTrxRepo, db)
	salesController := sales2.NewHTTPController(salesService)

	// Start API
	sales2.SalesRoute(r, salesController, basePath)
}
