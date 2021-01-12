package handlers

import (
	"log"
	"net/http"

	"github.com/go-microservices/data"
)

type Products struct {
	logger *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (products *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		products.getProducts(rw, r)
	case http.MethodPost:
		products.addProduct(rw, r)
	default:
		rw.WriteHeader(http.StatusMethodNotAllowed)	
	}
	

	// if r.Method == http.MethodGet {
	// 	products.getProducts(rw, r)
	// 	return
	// }

	// if r.Method == http.MethodPost {
	// 	products.addProduct(rw, r)
	// 	return
	// }

	//catch all
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (products *Products) getProducts(rw http.ResponseWriter, r *http.Request) {
	productsList := data.GetProducts()
	err := productsList.ToJSON(rw)

	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}

}

func (products *Products) addProduct(rw http.ResponseWriter, r *http.Request) {
	product := &data.Product{}
	err := product.FromJSON(r.Body)

	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
	}

	products.logger.Printf("%#v", product)
	rw.WriteHeader(http.StatusAccepted)
}