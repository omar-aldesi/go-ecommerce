package v1

import (
	"ecommerce/app/core"
	"ecommerce/app/crud"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// ListCategories
// @Summary List categories
// @Description Retrieves a list of all categories
// @Tags categories
// @Accept json
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /categories/list [get]
func ListCategories(c *gin.Context) {
	db := core.GetDB()

	categories, err := crud.ListCategories(db)
	if err != nil {
		core.CustomErrorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"categories": categories})
}

// ListSubCategories
// @Summary List subcategories
// @Description Retrieves a list of all subcategories
// @Tags categories
// @Accept json
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /categories/list-subcategories [get]
func ListSubCategories(c *gin.Context) {
	db := core.GetDB()

	subCategories, err := crud.ListSubCategories(db)
	if err != nil {
		core.CustomErrorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"sub_categories": subCategories})
}

// GetCategory
// @Summary Get category details
// @Description Retrieves the details of a category by its ID
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /categories/get/{id} [get]
func GetCategory(c *gin.Context) {
	db := core.GetDB()
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		core.CustomErrorResponse(c, &core.HTTPError{
			Message:    "Invalid id",
			StatusCode: http.StatusBadRequest,
		})
		return
	}
	category, err := crud.GetCategoryByID(db, uint(id))
	if err != nil {
		core.CustomErrorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"category": category})
}

// GetSubCategory
// @Summary Get subcategory details
// @Description Retrieves the details of a subcategory by its ID
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Subcategory ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /categories/get-subcategory/{id} [get]
func GetSubCategory(c *gin.Context) {
	db := core.GetDB()
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		core.CustomErrorResponse(c, &core.HTTPError{
			Message:    "Invalid id",
			StatusCode: http.StatusBadRequest,
		})
		return
	}
	category, err := crud.GetSubCategoryByID(db, uint(id))
	if err != nil {
		core.CustomErrorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"sub_category": category})
}

func CategoriesRouter(router *gin.Engine) {
	public := router.Group("/api/v1/categories")
	{
		public.GET("/list", ListCategories)
		public.GET("/list-subcategories", ListSubCategories)
		public.GET("/get/:id", GetCategory)
		public.GET("/get-subcategory/:id", GetSubCategory)
	}
}
