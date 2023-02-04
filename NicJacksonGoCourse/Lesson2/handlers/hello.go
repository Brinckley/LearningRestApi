package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

type Hello struct {
	l *log.Logger
}

func NewHello(l *log.Logger) *Hello {
	return &Hello{l: l}
}

func (h *Hello) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.l.Println("Hello World")

	d, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "error reading the body", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "Value from body : %s", string(d))
}
