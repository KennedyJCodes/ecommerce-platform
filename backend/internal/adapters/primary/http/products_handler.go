// Package http implements HTTP handlers and adapters for the ecommerce-platform application.
// This file contains the ProductsHandler, which manages product-related requests such as retrieving all products, filtering by brand, or fetching by specific ID.
package http

import (
	"net/http"
	"strconv"

	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/ports/input"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/pkg/errors"
	httpUtil "github.com/David-Alejandro-Jimenez/ecommerce-platform/pkg/http"
	"github.com/gorilla/mux"
)

// ProductsHandler handles HTTP requests related to the product catalog.
// It acts as a primary adapter, translating HTTP parameters and routes into call to the productsService domain port.
type ProductsHandler struct {
	// productsService is the input port used to interact with product business logic.
	productsService input.ProductsGetService
}

// NewProductsHandler creates and returns a new instance of ProductsHandler.
// Parameters:
//   - productsService: an implementation of the ProductsGetService port.
func NewProductsHandler(productsService input.ProductsGetService) *ProductsHandler {
	return &ProductsHandler{
		productsService: productsService,
	}
}

// Handle processes requests to retrieve the full list of products.
// It returns a JSON array of products with a 200 OK status on success, or a 500 Internal Server Error if the service fails.
func (h *ProductsHandler) Handle(w http.ResponseWriter, r *http.Request) {
	products, err := h.productsService.GetProducts()
	if err != nil {
		httpUtil.HandleError(w, errors.NewInternalError("Error getting products"))
		return
	}

	httpUtil.SendJSONResponse(w, http.StatusOK, products)
}

// HandleGetByID processes requests to retrieve a single product by its unique ID.
// Implementation details:
//  1. Extracts the "id" variable from the request URL using gorilla/mux.
//  2. Converts the ID from string to integer.
//  3. Validates that the ID is a non-negative value.
//  4. Calls the domain service to fetch the product.
//  5. Returns 404 if not found, 400 for invalid ID format, or 200 with the product data.
func (h *ProductsHandler) HandleGetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, exists := vars["id"]
	if !exists {
		httpUtil.HandleError(w, errors.NewBadRequestError("ID parameter is required"))
		return
	}

	// Conversion and validation of the ID parameter
	id, err := strconv.Atoi(idStr)
	if err != nil {
		httpUtil.HandleError(w, errors.NewBadRequestError("Invalid ID format"))
		return
	}

	if id < 0 {
		httpUtil.HandleError(w, errors.NewBadRequestError("Invalid ID format"))
		return
	}

	product, err := h.productsService.GetProductByID(id)
	if err != nil {
		if errors.IsNotFound(err) {
			httpUtil.HandleError(w, errors.NewNotFoundError("Product not found"))
			return
		}
		httpUtil.HandleError(w, errors.NewInternalError("Error getting product"))
		return
	}

	httpUtil.SendJSONResponse(w, http.StatusOK, product)
}

// HandleGetByBrand processes requests to filter products by their brand name.
// Steps:
//  1. Extracts the "brand" parameter from the URL variables.
//  2. Delegates the filtering logic to the domain service.
//  3. Responds with 200 OK and the list of matching products, or 404 if no products match.
func (h *ProductsHandler) HandleGetByBrand(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	brand, exists := vars["brand"]
	if !exists {
		httpUtil.HandleError(w, errors.NewBadRequestError("Brand parameter is required"))
		return
	}

	products, err := h.productsService.GetProductsByBrand(brand)
	if err != nil {
		if errors.IsNotFound(err) {
			httpUtil.HandleError(w, errors.NewNotFoundError("Product not found"))
			return
		}
		httpUtil.HandleError(w, errors.NewInternalError("Error getting product"))
		return
	}

	httpUtil.SendJSONResponse(w, http.StatusOK, products)
}
