package models

type Product struct {
	ProductID int `db:"product_id"`
	Product_Name string `db:"product_name"`
	Product_Description string `db:"product_description"`
	Product_Price float64 `db:"product_price"`
	Stock_Quantity int `db:"stock_quantity"`
	Brand string `db:"brand"`
	Image_Url string `db:"image_url"`
	Movement_Type string `db:"movement_type"`
}
