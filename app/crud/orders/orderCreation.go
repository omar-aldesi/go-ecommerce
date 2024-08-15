package orders

import (
	"ecommerce/app/core"
	"ecommerce/app/models"
	"ecommerce/app/schemas"
	"fmt"
	"gorm.io/gorm"
	"log"
	"net/http"
)

// helper functions

// Generic function to check if an item exists in a slice
func contains[T any](slice []T, item T, comparer func(T, T) bool) bool {
	for _, element := range slice {
		if comparer(element, item) {
			return true
		}
	}
	return false
}

// Comparer functions
func addonComparer(a, b models.Addon) bool {
	return a.ID == b.ID
}

func variationComparer(a, b models.ProductVariation) bool {
	return a.ID == b.ID
}

func contain(slice []uint, item uint) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func checkOrderType(orderType string) error {
	switch orderType {
	case "pickup", "shipping":
		return nil
	}
	return fmt.Errorf("order Type Not Found")
}

// Order Creation Functions

func CreateOrder(tx *gorm.DB, user models.User, orderData schemas.OrderCreationSchema) error {
	log.Printf("Starting order creation for user ID: %d", user.ID)

	if err := checkOrderType(orderData.OrderType); err != nil {
		log.Printf("Invalid order type: %s", err)
		return err
	}

	newOrder, err := createInitialOrder(tx, user, orderData)
	if err != nil {
		return err
	}

	totalPrice, err := processOrderItems(tx, newOrder, orderData)
	if err != nil {
		return err
	}

	if err := createNewPayment(tx, user.ID, newOrder, orderData.Payment); err != nil {
		return err
	}

	if err := finalizeOrder(tx, newOrder, totalPrice, orderData); err != nil {
		return err
	}

	log.Printf("Order created successfully. Order ID: %d, Total Price: %.2f", newOrder.ID, totalPrice)
	return nil
}

func createInitialOrder(tx *gorm.DB, user models.User, orderData schemas.OrderCreationSchema) (*models.Order, error) {
	newOrder := models.Order{
		UserID:      user.ID,
		Type:        orderData.OrderType,
		IsScheduled: orderData.IsScheduled,
		BranchID:    orderData.BranchID,
		Status:      "pending",
	}

	if orderData.IsScheduled {
		newOrder.ScheduleTime = orderData.ScheduleAt
	}

	if err := tx.Create(&newOrder).Error; err != nil {
		log.Printf("Error creating order: %s", err)
		return nil, &core.HTTPError{
			StatusCode: http.StatusInternalServerError,
			Message:    fmt.Sprintf("Error creating order: %s", err),
		}
	}

	log.Printf("Initial order created. Order ID: %d", newOrder.ID)
	return &newOrder, nil
}

func processOrderItems(tx *gorm.DB, newOrder *models.Order, orderData schemas.OrderCreationSchema) (float64, error) {
	var totalPrice float64

	for _, product := range orderData.Products {
		itemTotalPrice, err := processOrderItem(tx, newOrder, product, orderData.BranchID)
		if err != nil {
			return 0, err
		}
		totalPrice += itemTotalPrice
	}

	log.Printf("All order items processed. Total Price: %.2f", totalPrice)
	return totalPrice, nil
}

func processOrderItem(tx *gorm.DB, newOrder *models.Order, product schemas.OrderItemSchema, branchID uint) (float64, error) {
	var dbProduct models.Product
	if err := tx.Preload("Addons").Preload("Variations").First(&dbProduct, product.ProductID).Error; err != nil {
		log.Printf("Product not found. ID: %d", product.ProductID)
		return 0, &core.HTTPError{
			StatusCode: http.StatusNotFound,
			Message:    fmt.Sprintf("Product %d not found", product.ProductID),
		}
	}
	if !checkProductStocks(dbProduct, product.Quantity) {
		return 0, &core.HTTPError{
			StatusCode: http.StatusNotFound,
			Message:    fmt.Sprintf("Product %d incefficient stocks", product.ProductID),
		}
	}
	if dbProduct.BranchID != branchID {
		log.Printf("Product branch mismatch. Product ID: %d, Branch ID: %d", product.ProductID, branchID)
		return 0, &core.HTTPError{
			StatusCode: http.StatusBadRequest,
			Message:    fmt.Sprintf("Product %d Branch ID does not match", product.ProductID),
		}
	}

	newOrderItem := models.OrderItem{
		OrderID:   newOrder.ID,
		ProductID: product.ProductID,
		Quantity:  product.Quantity,
	}

	if err := tx.Create(&newOrderItem).Error; err != nil {
		log.Printf("Error creating order item. Order ID: %d, Error: %s", newOrder.ID, err)
		return 0, &core.HTTPError{
			StatusCode: http.StatusBadRequest,
			Message:    fmt.Sprintf("Error creating order item %d: %s", newOrder.ID, err),
		}
	}

	itemTotalPrice, err := processAddonsAndVariations(tx, &newOrderItem, product, dbProduct)

	if err != nil {
		return 0, err
	}

	if err := tx.Save(&newOrderItem).Error; err != nil {
		log.Printf("Error updating order item. ID: %d, Error: %s", newOrderItem.ID, err)
		return 0, &core.HTTPError{
			StatusCode: http.StatusInternalServerError,
			Message:    fmt.Sprintf("Error updating order item : %s", err),
		}
	}

	log.Printf("Order item processed. Item ID: %d, Total Price: %.2f", newOrderItem.ID, itemTotalPrice)
	return itemTotalPrice, nil
}

func processAddonsAndVariations(tx *gorm.DB, newOrderItem *models.OrderItem, product schemas.OrderItemSchema, dbProduct models.Product) (float64, error) {
	var itemTotalPrice float64

	// Process addons
	for _, addon := range product.Addons {
		addonPrice, err := processAddon(tx, newOrderItem, addon, dbProduct)
		if err != nil {
			return 0, err
		}
		itemTotalPrice += addonPrice
	}

	// Process variations
	variationPrice, err := processVariations(tx, newOrderItem, product.Variations, dbProduct)
	if err != nil {
		return 0, err
	}
	itemTotalPrice += variationPrice

	return itemTotalPrice, nil
}

func processAddon(tx *gorm.DB, newOrderItem *models.OrderItem, addon schemas.AddonSchema, dbProduct models.Product) (float64, error) {
	var dbAddon models.Addon
	if err := tx.First(&dbAddon, addon.AddonID).Error; err != nil {
		log.Printf("Addon not found. ID: %d", addon.AddonID)
		return 0, &core.HTTPError{
			StatusCode: http.StatusNotFound,
			Message:    fmt.Sprintf("Addon %d not found", addon.AddonID),
		}
	}

	if !contains(dbProduct.Addons, dbAddon, addonComparer) {
		log.Printf("Addon not available for product. Addon ID: %d, Product ID: %d", addon.AddonID, dbProduct.ID)
		return 0, &core.HTTPError{
			StatusCode: http.StatusBadRequest,
			Message:    fmt.Sprintf("Addon %d not found in product available addons", addon.AddonID),
		}
	}

	newOrderItemAddon := models.OrderItemAddon{
		OrderItemID: newOrderItem.ID,
		AddonID:     addon.AddonID,
		Quantity:    addon.Quantity,
	}
	newOrderItem.SelectedAddons = append(newOrderItem.SelectedAddons, newOrderItemAddon)

	addonPrice := dbAddon.Price + dbAddon.Tax
	log.Printf("Addon processed. Addon ID: %d, Price: %.2f", addon.AddonID, addonPrice)
	return addonPrice, nil
}

func processVariations(tx *gorm.DB, newOrderItem *models.OrderItem, variations []schemas.ProductVariationSchema, dbProduct models.Product) (float64, error) {
	var totalVariationPrice float64
	var processedVariations []uint

	for _, variation := range variations {
		variationPrice, err := processVariation(tx, newOrderItem, variation, dbProduct)
		if err != nil {
			return 0, err
		}
		totalVariationPrice += variationPrice
		processedVariations = append(processedVariations, variation.ProductVariationID)
	}

	if err := checkRequiredVariations(dbProduct, processedVariations); err != nil {
		return 0, err
	}

	return totalVariationPrice, nil
}

func processVariation(tx *gorm.DB, newOrderItem *models.OrderItem, variation schemas.ProductVariationSchema, dbProduct models.Product) (float64, error) {
	var dbVariation models.ProductVariation
	if err := tx.First(&dbVariation, variation.ProductVariationID).Error; err != nil {
		log.Printf("Variation not found. ID: %d", variation.ProductVariationID)
		return 0, &core.HTTPError{
			StatusCode: http.StatusNotFound,
			Message:    fmt.Sprintf("Variation %d not found", variation.ProductVariationID),
		}
	}

	if !contains(dbProduct.Variations, dbVariation, variationComparer) {
		log.Printf("Variation not available for product. Variation ID: %d, Product ID: %d", variation.ProductVariationID, dbProduct.ID)
		return 0, &core.HTTPError{
			StatusCode: http.StatusBadRequest,
			Message:    fmt.Sprintf("Variation %d not found in product available variations", variation.ProductVariationID),
		}
	}

	newOrderItemVariation := models.OrderItemVariation{
		OrderItemID:        newOrderItem.ID,
		ProductVariationID: variation.ProductVariationID,
	}

	var variationPrice float64
	for _, option := range variation.Options {
		optionPrice, err := processVariationOption(tx, &newOrderItemVariation, option)
		if err != nil {
			return 0, err
		}
		variationPrice += optionPrice
	}

	newOrderItem.SelectedVariations = append(newOrderItem.SelectedVariations, newOrderItemVariation)
	log.Printf("Variation processed. Variation ID: %d, Price: %.2f", variation.ProductVariationID, variationPrice)
	return variationPrice, nil
}

func processVariationOption(tx *gorm.DB, newOrderItemVariation *models.OrderItemVariation, option schemas.VariationOptionSchema) (float64, error) {
	var dbVariationOption models.VariationOption
	if err := tx.First(&dbVariationOption, option.VariationOptionID).Error; err != nil {
		log.Printf("VariationOption not found. ID: %d", option.VariationOptionID)
		return 0, &core.HTTPError{
			StatusCode: http.StatusNotFound,
			Message:    fmt.Sprintf("VariationOption %d not found", option.VariationOptionID),
		}
	}

	newOrderItemVariation.SelectedOptions = append(newOrderItemVariation.SelectedOptions, dbVariationOption)
	log.Printf("VariationOption processed. Option ID: %d, Price: %.2f", option.VariationOptionID, dbVariationOption.Price)
	return dbVariationOption.Price, nil
}

func checkRequiredVariations(dbProduct models.Product, processedVariations []uint) error {
	for _, productVariation := range dbProduct.Variations {
		if productVariation.Required && !contain(processedVariations, productVariation.ID) {
			log.Printf("Required variation missing. Variation ID: %d", productVariation.ID)
			return &core.HTTPError{
				StatusCode: http.StatusBadRequest,
				Message:    fmt.Sprintf("Variation %d is required", productVariation.ID),
			}
		}
	}
	return nil
}

func finalizeOrder(tx *gorm.DB, newOrder *models.Order, totalPrice float64, orderData schemas.OrderCreationSchema) error {
	newOrder.Total = totalPrice

	if orderData.OrderType == "shipping" {
		if err := processShippingAddress(tx, newOrder, orderData.ShippingAddress); err != nil {
			return err
		}
	}
	if err := tx.Save(newOrder).Error; err != nil {
		log.Printf("Error updating order total price. Order ID: %d, Error: %s", newOrder.ID, err)
		return &core.HTTPError{
			StatusCode: http.StatusInternalServerError,
			Message:    fmt.Sprintf("Error updating order total price: %s", err),
		}
	}

	log.Printf("Order finalized. Order ID: %d, Total Price: %.2f", newOrder.ID, totalPrice)
	return nil
}

func processShippingAddress(tx *gorm.DB, newOrder *models.Order, shippingAddress schemas.ShippingAddressSchema) error {
	var newShippingAddress models.ShippingAddress

	result := tx.Where(models.ShippingAddress{
		AddressLine1: shippingAddress.AddressLine1,
		AddressLine2: shippingAddress.AddressLine2,
		City:         shippingAddress.City,
		Country:      shippingAddress.Country,
		Postcode:     shippingAddress.Postcode,
		State:        shippingAddress.State,
		OrderID:      newOrder.ID,
	}).FirstOrCreate(&newShippingAddress)

	if result.Error != nil {
		log.Printf("Error getting or creating shipping address. Order ID: %d, Error: %s", newOrder.ID, result.Error)
		return &core.HTTPError{
			StatusCode: http.StatusInternalServerError,
			Message:    fmt.Sprintf("Error getting or creating shipping address: %s", result.Error),
		}
	}

	newOrder.ShippingAddress = newShippingAddress
	log.Printf("Shipping address processed. Address ID: %d, Order ID: %d", newShippingAddress.ID, newOrder.ID)
	return nil
}

func checkProductStocks(product models.Product, quantity uint) bool {
	// TODO : Check for last daily stock update
	switch product.StockType {
	case "UNLIMITED":
		return true
	case "FIXED", "DAILY":
		if product.Stock < quantity {
			return false
		} else {
			return true
		}
	default:
		return false
	}
}

func createNewPayment(tx *gorm.DB, userID uint, newOrder *models.Order, paymentData schemas.NewPaymentSchema) error {
	newPayment := models.Payment{
		UserID:              userID,
		Amount:              paymentData.Amount,
		Currency:            paymentData.Currency,
		Status:              paymentData.Status,
		Gateway:             paymentData.Gateway,
		PaymentIntentID:     paymentData.PaymentIntentID,
		PaymentClientSecret: paymentData.PaymentClientSecret,
		ReceiptEmail:        paymentData.ReceiptEmail,
		OrderID:             newOrder.ID,
	}
	if err := tx.Create(&newPayment).Error; err != nil {
		return &core.HTTPError{
			StatusCode: http.StatusInternalServerError,
			Message:    fmt.Sprintf("Error creating new payment: %s", err),
		}
	}
	newOrder.Payment = newPayment
	return nil
}
