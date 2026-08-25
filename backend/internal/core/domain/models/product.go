// Package models defines core domain entities for the ecommerce-platform application.
// This file declares the Product type, the central entity representing a watch in the catalog.
package models

// Product represents a watch available in the store's inventory.
// It encapsulates all physical and commercial attributes of a watch, including technical specifications like movement type and logistical data like stock quantity.
// The struct uses `db` tags for direct mapping from MySQL table columns using sqlx.
type Product struct {
	// Product_ID: unique identifier for the product (Primary Key in DB).
	Product_ID int `db:"product_id"`

	// Product_Name: the commercial name of the watch.
	Product_Name string `db:"product_name"`

	// Product_Description: a detailed text describing the watch features and style.
	Product_Description string `db:"product_description"`

	// Product_Price: the selling price. Using float64 allows for decimal precision in currency.
	Product_Price float64 `db:"product_price"`

	// Stock_Quantity: current number of units available in inventory.
	Stock_Quantity int `db:"stock_quantity"`

	// Brand: the manufacturer of the watch (e.g., Rolex, Casio, Seiko).
	Brand string `db:"brand"`

	// Image_URL: the path or link to the product's visual asset.
	Image_URL string `db:"image_url"`

	// Movement_Type: technical specification of the watch mechanism (e.g., Quartz, Automatic, Manual).
	Movement_Type string `db:"movement_type"`
}
