package handlers

import (
	"Lesson3/data"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{
		l: l,
	}
}

func (p *Products) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		p.getProducts(w, r)
		return
	}

	if r.Method == http.MethodPost {
		p.addProduct(w, r)
		return
	}

	if r.Method == http.MethodPut {
		p.l.Println("PUT", r.URL.Path)
		// expect the id in the URI
		reg := regexp.MustCompile(`/([0-9]+)`)
		g := reg.FindAllStringSubmatch(r.URL.Path, -1)

		if len(g) != 1 {
			p.l.Println("Invalid URI more than one id")
			http.Error(w, "Invalid URI", http.StatusBadRequest)
			return
		}

		if len(g[0]) != 2 {
			p.l.Println("Invalid URI more than one capture group")
			http.Error(w, "Invalid URI", http.StatusBadRequest)
			return
		}

		idString := g[0][1]
		id, err := strconv.Atoi(idString)
		if err != nil {
			p.l.Println("Invalid URI unable to convert to numer", idString)
			http.Error(w, "Invalid URI", http.StatusBadRequest)
			return
		}

		p.updateProduct(id, w, r)
		return
	}
	// catch all
	w.WriteHeader(http.StatusNotImplemented)
}

func (p *Products) getProducts(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handling product POST method")
	pl := data.GetProducts()
	err := pl.ToJson(w)
	if err != nil {
		http.Error(w, "error marshalling data", http.StatusInternalServerError)
	}
}

func (p *Products) addProduct(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handling product GET method")

	prod := &data.Product{}
	err := prod.FromJson(r.Body)
	if err != nil {
		http.Error(w, "error unmarshalling data", http.StatusBadRequest)
	}

	data.AddProduct(prod)
}

func (p *Products) updateProduct(id int, w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handling product PUT method")

	prod := &data.Product{}
	err := prod.FromJson(r.Body)
	if err != nil {
		http.Error(w, "error unmarshalling data", http.StatusBadRequest)
	}

	err = data.UpdateProduct(id, prod)
	if err == data.ErrProductNotFound {
		http.Error(w, "error product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, "error product not found", http.StatusInternalServerError)
		return
	}
}
