package v1

import (
	"ecommerce/app/core"
	"ecommerce/app/crud"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// ListBranches
// @Summary List all branches
// @Description Get a list of all branches
// @Tags branches
// @Accept json
// @Produce json
// @Success 200 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /branches/list [get]
func ListBranches(c *gin.Context) {
	db := core.GetDB()

	branches, err := crud.ListBranches(db)
	if err != nil {
		core.CustomErrorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"branches": branches})

}

// GetBranch
// @Summary Get branch by ID
// @Description Get details of a branch by its ID
// @Tags branches
// @Accept json
// @Produce json
// @Param id path int true "Branch ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /branches/get/{id} [get]
func GetBranch(c *gin.Context) {
	db := core.GetDB()
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		core.CustomErrorResponse(c, &core.HTTPError{
			Message:    "Invalid id",
			StatusCode: http.StatusBadRequest,
		})
		return
	}
	branch, err := crud.GetBranchByID(db, uint(id))
	if err != nil {
		core.CustomErrorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"branch": branch})
}

func BranchesRouter(router *gin.Engine) {
	public := router.Group("/api/v1/branches")
	{
		public.GET("/list", ListBranches)
		public.GET("/get/:id", GetBranch)
	}
}
