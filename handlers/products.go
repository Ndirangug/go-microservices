package handlers

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/go-microservices/data"
)

type Products struct {
	logger *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (products *Products) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		products.getProducts(rw, req)
	case http.MethodPost:
		products.addProduct(rw, req)
	case http.MethodPut:
		id, err := parseID(req)

		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
		}

		products.updateProduct(id, rw, req)

	default:
		rw.WriteHeader(http.StatusMethodNotAllowed)
	}

}

func (products *Products) getProducts(rw http.ResponseWriter, r *http.Request) {
	products.logger.Println("handle GET products")
	productsList := data.GetProducts()
	err := productsList.ToJSON(rw)

	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}

}

func (products *Products) addProduct(rw http.ResponseWriter, r *http.Request) {
	products.logger.Println("handle POST product")

	product := &data.Product{}
	err := product.FromJSON(r.Body)

	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
	}

	data.AddProduct(product)
	rw.WriteHeader(http.StatusAccepted)
}

func (products *Products) updateProduct(id int, rw http.ResponseWriter, r *http.Request) {
	products.logger.Println("handle PUT product")

	product := &data.Product{}
	err := product.FromJSON(r.Body)

	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
	}

	err = data.UpdateProduct(id, product)

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

func parseID(req *http.Request) (int, error) {
	regex := regexp.MustCompile(`/([0-9]+)`)
	g := regex.FindAllStringSubmatch(req.URL.Path, -1)
	fmt.Printf("%#v\n", g)

	if len(g) != 1 {
		return 0, fmt.Errorf("Invalid URI - More or less than one id")
	}
	// if len(g[0]) != 1 {
	// 	http.Error(rw, "Invalid URI", http.StatusBadRequest)
	// 	return
	// }

	idString := g[0][1]
	id, err := strconv.Atoi(idString)

	if err != nil {
		return 0, fmt.Errorf("Couldnt parse idString into an int")
	}

	return id, nil
}
