package crud

import (
	"ecommerce/app/core"
	"ecommerce/app/models"
	"errors"
	"gorm.io/gorm"
	"net/http"
)

func ListCategories(db *gorm.DB) ([]models.Category, error) {
	var dbCategories []models.Category

	if err := db.Preload("SubCategories").Find(&dbCategories).Error; err != nil {
		return nil, &core.HTTPError{
			Message:    "Error getting categories",
			StatusCode: http.StatusInternalServerError,
		}
	}
	return dbCategories, nil
}

func ListSubCategories(db *gorm.DB) ([]models.SubCategory, error) {
	var dbSubCategories []models.SubCategory

	if err := db.Find(&dbSubCategories).Error; err != nil {
		return nil, &core.HTTPError{
			Message:    "Error getting subCategories",
			StatusCode: http.StatusInternalServerError,
		}
	}
	return dbSubCategories, nil
}

func GetCategoryByID(db *gorm.DB, id uint) (models.Category, error) {
	var dbCategory models.Category
	if err := db.Preload("SubCategories").First(&dbCategory, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Category{}, &core.HTTPError{
				Message:    "Category not found",
				StatusCode: http.StatusNotFound,
			}
		}
		return models.Category{}, &core.HTTPError{
			Message:    "Error getting category",
			StatusCode: http.StatusInternalServerError,
		}
	}
	return dbCategory, nil
}

func GetSubCategoryByID(db *gorm.DB, id uint) (models.SubCategory, error) {
	var dbSubCategory models.SubCategory
	if err := db.First(&dbSubCategory, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.SubCategory{}, &core.HTTPError{
				Message:    "SubCategory not found",
				StatusCode: http.StatusNotFound,
			}
		}
		return models.SubCategory{}, &core.HTTPError{
			Message:    "Error getting subCategory",
			StatusCode: http.StatusInternalServerError,
		}
	}
	return dbSubCategory, nil
}
