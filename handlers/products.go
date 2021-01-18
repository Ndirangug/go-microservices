package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-microservices/data"
	"github.com/gorilla/mux"
)

type Products struct {
	logger *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (products *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	products.logger.Println("handle GET products")
	productsList := data.GetProducts()
	err := productsList.ToJSON(rw)

	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}

}

func (products *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
	products.logger.Println("handle POST product")

	product := r.Context().Value(KeyProduct{}).(*data.Product)

	data.AddProduct(product)
	rw.WriteHeader(http.StatusAccepted)
}

func (products *Products) UpdateProduct(rw http.ResponseWriter, r *http.Request) {
	products.logger.Println("handle PUT product")
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	product := r.Context().Value(KeyProduct{}).(*data.Product)

	err := data.UpdateProduct(id, product)

	if err == data.ErrorProductNotFound {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "An error occurred", http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusAccepted)
}

type KeyProduct struct{}

func (products *Products) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		product := &data.Product{}
		err := product.FromJSON(req.Body)

		if err != nil {
			products.logger.Println("[Error] unmarshalling json", err)
			http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
			return
		}

		err = product.Validate()

		if err != nil {
			products.logger.Println("[Error] validating product", err)
			http.Error(rw, fmt.Sprintf("Error validating product: %s", err), http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(req.Context(), KeyProduct{}, product)
		req = req.WithContext(ctx)

		next.ServeHTTP(rw, req)
	})
}
