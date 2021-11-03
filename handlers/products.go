package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/nathantruant/products-api/data"
)

// KeyProduct is a key used for the Product object in the context
type KeyProduct struct{}

// Products is a http.Handler
type Products struct {
	l *log.Logger
}

// NewProducts creates a products handler with the given logger
func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

// getProducts returns the products from the data store
func (p *Products) GetProducts(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET Products")

	// get products from the data store
	lp := data.GetProducts()

	// serialize products to JSON
	err := lp.ToJSON(w)
	if err != nil {
		http.Error(w, "Unable to serialize json", http.StatusInternalServerError)
		return
	}
}

func (p *Products) AddProduct(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Product")

	prod := r.Context().Value(KeyProduct{}).(data.Product)
	data.AddProduct(&prod)
}

func (p *Products) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusBadRequest)
		return
	}

	p.l.Println("Handle PUT Product: ", id)
	prod := r.Context().Value(KeyProduct{}).(data.Product)

	err = data.UpdateProduct(id, &prod)
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
