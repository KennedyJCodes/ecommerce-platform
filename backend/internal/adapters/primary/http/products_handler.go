package http

import (
	"log"
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
		log.Println(err)
		httpUtil.HandleError(w, errors.NewInternalError("Error getting products"))
		return
	}

	httpUtil.SendJSONResponse(w, http.StatusOK, products)
}

func (h *ProductsHandler) HandleGetByID(w http.ResponseWriter, r *http.Request) {
	// Extraer el ID de los parámetros de la URL
	vars := mux.Vars(r)
	idStr, exists := vars["id"]
	if !exists {
		httpUtil.HandleError(w, errors.NewBadRequestError("ID parameter is required"))
		return
	}

	// Convertir el ID a entero
	id, err := strconv.Atoi(idStr)
	if err != nil {
		httpUtil.HandleError(w, errors.NewBadRequestError("Invalid ID format"))
		return
	}

	// Obtener el producto por ID
	product, err := h.productsService.GetProductByID(id)
	if err != nil {
		// Verificar si es un error de "no encontrado"
		if errors.IsNotFound(err) {
			httpUtil.HandleError(w, errors.NewNotFoundError("Product not found"))
			return
		}
		httpUtil.HandleError(w, errors.NewInternalError("Error getting product"))
		return
	}

	httpUtil.SendJSONResponse(w, http.StatusOK, product)
}
