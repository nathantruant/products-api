package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/nathantruant/products-api/data"
)

// Products is a http.Handler
type Products struct {
	l *log.Logger
}

// NewProducts creates a products handler with the given logger
func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

// ServeHTTP is the main entry point for the handler and satisfies the
// http.Handler interface
func (p *Products) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	// handle request for list of products
	case http.MethodGet:
		p.getProducts(w, r)
		return
	case http.MethodPost:
		p.addProduct(w, r)
		return
	case http.MethodPut:
		p.l.Println("PUT HIT")

		// expect the id in the URI
		reg := regexp.MustCompile(`/(\d+)`)
		g := reg.FindAllStringSubmatch(r.URL.Path, -1)

		// if one group containing one match isn't found, throw bad request error
		if len(g) != 1 || len(g[0]) != 2 {
			p.l.Println("Could not find ID in URI")
			http.Error(w, "Invalid URI", http.StatusBadRequest)
			return
		}

		idStr := g[0][1]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			p.l.Println("Invalid URI more than one id")
			http.Error(w, "Invalid URI", http.StatusBadRequest)
			return
		}

		p.updateProduct(id, w, r)
		return

	// catch all
	// if no method is satisfied, return an error
	default:
		http.Error(w, "Invalid endpoint", http.StatusMethodNotAllowed)
	}
}

// getProducts returns the products from the data store
func (p *Products) getProducts(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET Products")

	// get products from the data store
	lp := data.GetProducts()

	// serialize products to JSON
	err := lp.ToJSON(w)
	if err != nil {
		http.Error(w, "Unable to encode json", http.StatusInternalServerError)
		return
	}
}

func (p *Products) addProduct(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Product")

	prod := &data.Product{}
	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(w, "Unable to decode JSON", http.StatusBadRequest)
		return
	}

	data.AddProduct(prod)
}

func (p *Products) updateProduct(id int, w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle PUT Product")
	prod := &data.Product{}
	err := prod.FromJSON(r.Body)
	if err != nil {
		p.l.Println("FromJSON err: ", err)
		http.Error(w, "Unable to decode JSON", http.StatusBadRequest)
		return
	}

	err = data.UpdateProduct(id, prod)
	if err == data.ErrProductNotFound {
		p.l.Println(err)
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		p.l.Println(err)
		http.Error(w, "Unable to update product", http.StatusInternalServerError)
		return
	}
}
