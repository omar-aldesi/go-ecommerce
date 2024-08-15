package v1

import (
	"ecommerce/app/core"
	"ecommerce/app/crud"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// ListProducts
// @Summary List products
// @Description Retrieves a list of products with optional filtering and pagination
// @Tags products
// @Accept json
// @Produce json
// @Param limit query int false "Number of results to return" default(10)
// @Param offset query int false "Number of results to skip" default(0)
// @Param category_id query int false "Filter by category ID"
// @Param subcategory_id query int false "Filter by subcategory ID"
// @Param branch_id query int false "Filter by branch ID"
// @Param min_price query float64 false "Filter by minimum price"
// @Param max_price query float64 false "Filter by maximum price"
// @Param search query string false "Search products by name or description"
// @Param in_stock query bool false "Filter by in-stock products"
// @Param sort_by query string false "Sort by field (e.g., price, name)"
// @Param sort_order query string false "Sort order (asc or desc)" default(asc)
// @Success 200 {array} schemas.ProductResponseSchema
// @Failure 400 {object} map[string]interface{}
// @Router /products/list [get]
func ListProducts(c *gin.Context) {
	db := core.GetDB()

	limit, _ := strconv.Atoi(c.Query("limit"))
	offset, _ := strconv.Atoi(c.Query("offset"))

	filters := make(map[string]interface{})

	// Parse and add filters
	if categoryID, err := strconv.ParseUint(c.Query("category_id"), 10, 32); err == nil {
		filters["category_id"] = uint(categoryID)
	}
	if subcategoryID, err := strconv.ParseUint(c.Query("subcategory_id"), 10, 32); err == nil {
		filters["subcategory_id"] = uint(subcategoryID)
	}
	if branchID, err := strconv.ParseUint(c.Query("branch_id"), 10, 32); err == nil {
		filters["branch_id"] = uint(branchID)
	}
	if minPrice, err := strconv.ParseFloat(c.Query("min_price"), 64); err == nil {
		filters["min_price"] = minPrice
	}
	if maxPrice, err := strconv.ParseFloat(c.Query("max_price"), 64); err == nil {
		filters["max_price"] = maxPrice
	}
	if search := c.Query("search"); search != "" {
		filters["search"] = search
	}
	if inStock, err := strconv.ParseBool(c.Query("in_stock")); err == nil {
		filters["in_stock"] = inStock
	}
	if sortBy := c.Query("sort_by"); sortBy != "" {
		filters["sort_by"] = sortBy
		filters["sort_order"] = c.Query("sort_order") // Default to ASC if not provided
	}
	products, err := crud.ListProducts(db, limit, offset, filters)
	if err != nil {
		core.CustomErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"products": products})
}

// GetProduct
// @Summary Get product details
// @Description Retrieves the details of a product by its ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} schemas.ProductResponseSchema
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /products/get/{id} [get]
func GetProduct(c *gin.Context) {
	db := core.GetDB()
	productID, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil {
		core.CustomErrorResponse(c, &core.HTTPError{
			Message:    "Product ID should be integer",
			StatusCode: http.StatusBadRequest,
		})
		return
	}
	product, err := crud.GetProductByID(db, uint(productID))
	if err != nil {
		core.CustomErrorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"product": product})
}

func ProductsRouter(router *gin.Engine) {
	public := router.Group("/api/v1/products")
	{
		public.GET("/list", ListProducts)
		public.GET("/get/:id", GetProduct)
	}
}
