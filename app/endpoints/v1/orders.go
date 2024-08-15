package v1

import (
	"ecommerce/app/core"
	"ecommerce/app/core/middlewares"
	"ecommerce/app/crud/orders"
	"ecommerce/app/models"
	"ecommerce/app/schemas"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// CreateOrder
// @Summary Create a new order
// @Description Creates a new order for the specified products and branch
// @Tags orders
// @Accept json
// @Produce json
// @Param request body schemas.OrderCreationSchema true "Order creation details"
// @Success 200 {object} schemas.OrderResponseSchema
// @Failure 400 {object} map[string]interface{}
// @Security BearerAuth
// @Router /orders/create [post]
func CreateOrder(c *gin.Context) {
	db := core.GetDB()

	user := c.MustGet("user").(models.User)

	var request schemas.OrderCreationSchema
	if err := c.ShouldBindJSON(&request); err != nil {
		core.HandleValidationErrors(c, err)
		return
	}

	tx := db.Begin()
	if err := orders.CreateOrder(tx, user, request); err != nil {
		core.CustomErrorResponse(c, err)
		tx.Rollback()
		return
	}
	tx.Commit()
	c.JSON(http.StatusCreated, gin.H{"message": "Order created successfully"})
}

// ListOrders
// @Summary List user orders
// @Description Retrieves a list of orders for the authenticated user
// @Tags orders
// @Accept json
// @Produce json
// @Success 200 {object} []schemas.OrderResponseSchema
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Security BearerAuth
// @Router /orders/list [get]
func ListOrders(c *gin.Context) {
	db := core.GetDB()
	user := c.MustGet("user").(models.User)
	userOrders, err := orders.ListUserOrders(db, user)

	if err != nil {
		core.CustomErrorResponse(c, err)
		return
	}
	userOrdersResponse := make([]schemas.OrderResponseSchema, len(userOrders))

	for i, order := range userOrders {
		userOrdersResponse[i] = order.ToResponse()
	}

	c.JSON(http.StatusOK, gin.H{"orders": userOrdersResponse})
}

// GetOrder
// @Summary Get order details
// @Description Retrieves the details of an order by its ID
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} schemas.OrderResponseSchema
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Security BearerAuth
// @Router /orders/get/{id} [get]
func GetOrder(c *gin.Context) {
	db := core.GetDB()
	user := c.MustGet("user").(models.User)
	orderIDStr := c.Param("id")

	// Convert orderID from string to uint
	orderID, err := strconv.ParseUint(orderIDStr, 10, 64)
	if err != nil {
		core.CustomErrorResponse(c, &core.HTTPError{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid order ID",
		})
		return
	}

	order, err := orders.GetOrderByID(db, user, uint(orderID))

	if err != nil {
		core.CustomErrorResponse(c, err)
	}

	c.JSON(http.StatusOK, gin.H{"order": order.ToResponse()})
}

// UpdateOrder
// @Summary Update an order status
// @Description Updates an order status by its ID
// @Tags orders
// @Accept json
// @Produce json
// @Param request body schemas.UpdateOrderStatusSchema true "Order creation details"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Security BearerAuth
// @Router /orders/update-status [put]
func UpdateOrder(c *gin.Context) {
	db := core.GetDB()
	user := c.MustGet("user").(models.User)

	type requestSchema struct {
		OrderID uint   `json:"order_id" binding:"required"`
		Status  string `json:"status" binding:"required"`
	}
	var request requestSchema

	if err := c.ShouldBindJSON(&request); err != nil {
		core.HandleValidationErrors(c, err)
		return
	}

	if err := orders.UpdateOrderStatus(db, user, request.OrderID, request.Status); err != nil {
		core.CustomErrorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Order updated successfully"})
}

func OrdersRouter(router *gin.Engine) {
	protected := router.Group("/api/v1/orders")
	protected.Use(middlewares.AuthMiddleware())
	{
		protected.POST("/create", CreateOrder)
		protected.GET("/list", ListOrders)
		protected.GET("/get/:id", GetOrder)
		protected.PUT("/update-status", UpdateOrder)
	}
}
