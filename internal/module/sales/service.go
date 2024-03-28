package sales

import (
	"context"
	product_sale_yearly_report "hacktiv8-techrawih-go-product-sale/internal/module/product-sale-yearly-report"
	"hacktiv8-techrawih-go-product-sale/internal/pkg/http/request/sales"
	"hacktiv8-techrawih-go-product-sale/internal/pkg/utils"
	"log"
	"strconv"

	"gorm.io/gorm"
)

type Service interface {
	SaveAll(context context.Context, request sales.Import) ([]*product_sale_yearly_report.ProductSaleYearlyReport, error)
}

type service struct {
	repo           Repository
	productTrxRepo product_sale_yearly_report.Repository
	db             *gorm.DB
}

func NewService(repo Repository, productTrxRepo product_sale_yearly_report.Repository, db *gorm.DB) Service {
	return &service{
		repo:           repo,
		productTrxRepo: productTrxRepo,
		db:             db,
	}
}

func (us *service) SaveAll(context context.Context, request sales.Import) (res []*product_sale_yearly_report.ProductSaleYearlyReport, err error) {
	records, err := utils.ReadCSV(request.FilePath)
	if err != nil {
		return nil, err
	}

	if err = us.repo.DeleteAll(); err != nil {
		log.Fatalf("failed to delete all sales records: %v", err)
		return nil, err
	}

	if err = us.productTrxRepo.DeleteAll(); err != nil {
		log.Fatalf("failed to delete all product_trxes records: %v", err)
		return nil, err
	}

	var salesStoreDto []*Sales
	for index, record := range records {
		if index != 0 {
			productRes, errProduct := us.repo.GetProductByName(record[0])
			if errProduct != nil {
				return nil, err
			}

			qtySold, _ := strconv.Atoi(record[1])
			saleAt, _ := utils.ConvertStringToTime(record[2])
			salesDto := &Sales{
				ProductID: uint(productRes.ID),
				QtySold:   qtySold,
				SaleAt:    saleAt,
			}
			salesStoreDto = append(salesStoreDto, salesDto)

			// _, err = us.repo.Save(context, *salesDto)
			// if err != nil {
			// 	return nil, err
			// }
		}
	}

	err = us.repo.SaveAll(context, salesStoreDto)
	if err != nil {
		return nil, err
	}

	_, err = product_sale_yearly_report.AggregateSalesByProduct(us.db)
	product_sale_yearly_report.UpdateStockProduct(us.db)
	return res, err
}
