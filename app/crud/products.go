package crud

import (
	"ecommerce/app/core"
	"ecommerce/app/models"
	"ecommerce/app/schemas"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

func ListProducts(db *gorm.DB, limit, offset int, filters map[string]interface{}) ([]schemas.ProductResponseSchema, error) {
	var dbProducts []models.Product
	query := db.Model(&models.Product{})

	// Apply filters
	if categoryID, ok := filters["category_id"].(uint); ok {
		query = query.Where("category_id = ?", categoryID)
	}

	if subcategoryID, ok := filters["subcategory_id"].(uint); ok {
		query = query.Where("subcategory_id = ?", subcategoryID)
	}

	if branchID, ok := filters["branch_id"].(uint); ok {
		query = query.Where("branch_id = ?", branchID)
	}

	if minPrice, ok := filters["min_price"].(float64); ok {
		query = query.Where("price >= ?", minPrice)
	}

	if maxPrice, ok := filters["max_price"].(float64); ok {
		query = query.Where("price <= ?", maxPrice)
	}

	if searchTerm, ok := filters["search"].(string); ok {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+searchTerm+"%", "%"+searchTerm+"%")
	}

	if inStock, ok := filters["in_stock"].(bool); ok && inStock {
		query = query.Where("stock > 0")
	}

	// Apply sorting
	if sortBy, ok := filters["sort_by"].(string); ok {
		order := "ASC"
		if sortOrder, ok := filters["sort_order"].(string); ok && strings.ToUpper(sortOrder) == "DESC" {
			order = "DESC"
		}
		query = query.Order(fmt.Sprintf("%s %s", sortBy, order))
	}

	// Apply pagination
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	// Preload related data
	query = query.Preload("Addons").Preload("Variations")

	// Execute the query
	if err := query.Find(&dbProducts).Error; err != nil {
		return nil, &core.HTTPError{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	// Convert to response schema
	dbProductsResponse := make([]schemas.ProductResponseSchema, len(dbProducts))
	for i, product := range dbProducts {
		dbProductsResponse[i] = product.ToResponse()
	}

	return dbProductsResponse, nil
}

func GetProductByID(db *gorm.DB, productID uint) (schemas.ProductResponseSchema, error) {
	var dbProduct models.Product
	if err := db.Preload("Addons").Preload("Variations").First(&dbProduct, productID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return schemas.ProductResponseSchema{}, &core.HTTPError{
				Message:    "Product not found",
				StatusCode: http.StatusNotFound,
			}
		}
		return schemas.ProductResponseSchema{}, &core.HTTPError{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}
	return dbProduct.ToResponse(), nil
}
