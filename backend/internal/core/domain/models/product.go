package models

type Product struct {
	Product_ID int `db:"Product_ID"`
	Product_Name string `db:"Product_Name"`
	Product_Description string `db:"Product_Description"`
	Product_Price float64 `db:"Product_Price"`
	Stock_Quantity int `db:"Stock_Quantity"`
	Brand string `db:"Brand"`
	Image_URL string `db:"Image_URL"`
	Movement_Type string `db:"Movement_Type"`
}
