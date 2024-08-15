package crud

import (
	"ecommerce/app/core"
	"ecommerce/app/models"
	"gorm.io/gorm"
	"net/http"
)

func ListBranches(db *gorm.DB) ([]models.Branch, error) {
	var dbBranches []models.Branch
	if err := db.Find(&dbBranches).Error; err != nil {
		return nil, &core.HTTPError{
			Message:    "Error fetching branches",
			StatusCode: http.StatusInternalServerError,
		}
	}
	return dbBranches, nil
}
func GetBranchByID(db *gorm.DB, id uint) (models.Branch, error) {
	var dbBranch models.Branch
	if err := db.First(&dbBranch, id).Error; err != nil {
		return dbBranch, &core.HTTPError{
			Message:    "Error fetching branch",
			StatusCode: http.StatusNotFound,
		}
	}
	return dbBranch, nil
}
