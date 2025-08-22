package http

import (
	"net/http"
	"strconv"

	"github.com/David-Alejandro-Jimenez/sale-watches/internal/core/ports/input"
	"github.com/David-Alejandro-Jimenez/sale-watches/pkg/errors"
	httpUtil "github.com/David-Alejandro-Jimenez/sale-watches/pkg/http"
	"github.com/gorilla/mux"
)

type ProductsHandler struct {
	productsService input.ProductsGetService
}

func NewProductsHandler(productsService input.ProductsGetService) *ProductsHandler {
	return &ProductsHandler{
		productsService: productsService,
	}
}

func (h *ProductsHandler) Handle(w http.ResponseWriter, r *http.Request) {
	products, err := h.productsService.GetProducts()
	if err != nil {
		httpUtil.HandleError(w, errors.NewInternalError("Error getting products"))
		return
	}

	httpUtil.SendJSONResponse(w, http.StatusOK, products)
}

func (h *ProductsHandler) HandleGetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, exists := vars["id"]
	if !exists {
		httpUtil.HandleError(w, errors.NewBadRequestError("ID parameter is required"))
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
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
