package seeds

import (
	"context"
	"hacktiv8-techrawih-go-product-sale/internal/module/product"
)

type Seed struct {
}

func Execute() {
	var seed Seed
	var ctx context.Context
	var repositoryProduct product.Repository

	seed.Product(ctx, repositoryProduct)
}
