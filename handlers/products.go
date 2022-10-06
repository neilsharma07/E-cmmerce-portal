// Package classification of product api
//
//Documentation for product API
//
//Schemes: http
// Base Path: /
//Version: 1.0.0
//
//Consumes:
// -application/json
//
//Produces:
// -application/json
// swagger:meta

package handlers

import (
	data "Practice/Product-API"
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type Products struct {
	l *log.Logger // Putting a log in the product handler
}

//New function using idiomatic go based approach
func NewProducts(l *log.Logger) *Products { //takes in logger as an input and returns the product handler
	return &Products{l}
}

// we need a serve http method
// modify handler to entertain diffferent request types
// func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
//	if r.Method == http.MethodGet { // if our request is a get request
//		p.getProducts(rw, r) // we will call our getProducts and forward the request to rw
//		return
//	}
//	// implement a post to add a new product
//	if r.Method == http.MethodPost { // new post method
//		p.addProduct(rw, r) // call a function addProduct
//		return
//	}
//	// implement a put to update an existing product
//	if r.Method == http.MethodPut { // new put method
//		//expect the id in the URI
//		p.l.Println("PUT", r.URL.Path)
//		// create a regexp
//		reg := regexp.MustCompile(`/([0-9]+)`) // the root path (/) and the ID // bracket signifies a capture group
//		// call the find all matches using find all string sub-match
//		g := reg.FindAllStringSubmatch(r.URL.Path, -1) // returns groups as a string array
//		// we should only have 1 group, so
//		if len(g) != 1 {
//			p.l.Println("Invalid URI, more than 1 ID")
//			http.Error(rw, "Invalid URI", http.StatusBadRequest)
//			return
//		}
//		if len(g[0]) != 2 {
//			p.l.Println("Invalid URI, more than 1 capture group")
//			http.Error(rw, "Invalid URI for element", http.StatusBadRequest)
//			return
//		}
//		// the capture group
//		// create the ID string which will be the second element of the capture group
//		idString := g[0][1]
//		id, err := strconv.Atoi(idString)
//		if err != nil {
//			p.l.Println("Unable to convert to integer")
//			http.Error(rw, "Invalid conversion", http.StatusBadRequest)
//			return
//		}
//		// call the put
//		p.updateProducts(id, rw, r)
//		return
//		p.l.Println("got id:", id)
//
//		// p := r.URL.Path // gives the / path of the URL, will return the /whatever part called by server
//		// we need just the ID from the entire path, how to extract?
//	}
//
//	// to handle an update, use the http verb called "put"
//	// CATCH ALL, if no method is  satisfied, then return an error
//	rw.WriteHeader(http.StatusMethodNotAllowed)
//  }

// convert the below getProducts to public, that is make it capital G, coz now we can use it in main
func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle Get Products")
	//when a req comes into server, we want to know where it is bound,
	// is it a get, retreiving a product
	// OR is it post, new product, update product
	//signature of serve http is an http resp writer (an interface) and an http request
	//we want to make a get request to serve.http and we want to return the product list
	//use encoding.json package
	// this package allows to convert a struct or a reference (to a pointer) into JSON file
	// to return the info, add a method to products package (main file)
	// this method will act as a data access model (DAM)
	lp := data.GetProducts() // gets us a list of products, lp now has a list of products
	// to return the list to user, convert lp to json array or json string using encoding
	// d, err := json.Marshal(lp) // do the usual go defencing programing
	// instead of the above statement, just call the ToJSON method and pass the rw here
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "unable to marshal json", http.StatusInternalServerError)
	}
	// rw.Write(d) // take resp writer and write the data out using the write method
}

// defining the addProduct function for the method and operation
func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) { // it is on p.Product, signature being same
	p.l.Println("Handle Post Products") //logging
	// create a new product object
	// prod := &data.Product{}      //JSOn we enter will be converted to data.Product struct
	// err := prod.FromJSON(r.Body) // pass the response body reader into this
	// if err != nil {
	//	http.Error(rw, "Unable to unmarhsal JSON", http.StatusBadRequest)
	// }
	// log the from JSON
	// p.l.Printf("Prod: %#v", prod) // printf gives nice o/p // %v is value, f also gives fields
	// add the product here
	prod := r.Context().Value(KeyProduct{}).(data.Product) // this will return an interface, so we have to cast it
	data.AddProduct(&prod)
}

// to manipulate data: take data from post and convert it into product object
// json encoder encodes json straight to rw, was better than marshal, this was product to json
// do the reverse to create json to product
// http request body is an io reader. decode body from request into an interface and pull that into a struct

// create the updateProducts function
func (p Products) UpdateProducts(rw http.ResponseWriter, r *http.Request) {
	// gorilla extracts varibaled into mux.Vars
	vars := mux.Vars(r)
	// vars will now have an id
	id, err := strconv.Atoi(vars["id"]) // coz vars is a map
	if err != nil {
		http.Error(rw, "Unable to convert ID string", http.StatusBadRequest)
		return
	}
	p.l.Println("Handle Update Products", id) //logging
	// replacing the above code with below code:
	// we will get our product from request
	prod := r.Context().Value(KeyProduct{}).(data.Product) // this will return an interface, so we have to cast it

	// create a new product object
	// prod := &data.Product{}     //JSOn we enter will be converted to data.Product struct
	// err = prod.FromJSON(r.Body) // pass the response body reader into this
	// if err != nil {
	// 	http.Error(rw, "Unable to unmarhsal JSON", http.StatusBadRequest)
	// }
	err = data.UpdateProduct(id, &prod) // method added in the API
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}
	// for normal error
	if err != nil {
		http.Error(rw, "Product not found (something else)", http.StatusInternalServerError)
		return
	}
}

// Using MIDDLEWARE
type KeyProduct struct{} // defining key for context
// create "next" http handler and return http handler
func (p Products) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// to validate request, take the product and de-serialize it
		// create a new product object
		prod := &data.Product{}      //JSOn we enter will be converted to data.Product struct
		err := prod.FromJSON(r.Body) // pass the response body reader into this
		if err != nil {
			http.Error(rw, "Unable to unmarhsal JSON", http.StatusBadRequest)
			return
		}
		//validate the product
		err = prod.Validate()
		if err != nil {
			p.l.Println("[ERROR]: deserializing product", err)
			http.Error(
				rw,
				fmt.Sprintf("error validating product: 5s ", err),
				http.StatusBadRequest,
			)
			return
		}
		// if no error, then put the product on the request coz the request can have a context
		ctx := context.WithValue(r.Context(), KeyProduct{}, prod) // create a new context from context that is already on the request.
		// create a copy of the request
		req := r.WithContext(ctx)
		// call the "next" handler
		next.ServeHTTP(rw, req) // pass in the rw and the req
	})
}
