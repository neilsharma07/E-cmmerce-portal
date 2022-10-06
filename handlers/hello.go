package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Hello struct {
	l *log.Logger
}

func NewHello(l *log.Logger) *Hello {
	return &Hello{l}
}

// defining a method on a struct in go
func (h *Hello) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	h.l.Println("Hello!")
	d, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, "error occured", http.StatusBadRequest)
		return
	}
	fmt.Fprintf(rw, "Hi %s", d)

}
