package handlers

import (
	"Lesson3/data"
	"log"
	"net/http"
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

	// catch all
	w.WriteHeader(http.StatusNotImplemented)
}

func (p *Products) getProducts(w http.ResponseWriter, r *http.Request) {
	pl := data.GetProducts()
	err := pl.ToJson(w)
	if err != nil {
		http.Error(w, "error marshalling data", http.StatusInternalServerError)
	}
}
