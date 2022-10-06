package data

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"io"
	"regexp"
	"time"
)

//The product struct defines the structure of an API product
// Also validate the product mentioned below using struct tags and specify which fields are valid using go validator
// Add custom validation types also

type Product struct { // use struct tags here, add annotations to fields and then we can later pickup these annotations
	ID          int     `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float32 `json:"price" validate:"gt=0"`
	SKU         string  `json:"sku" validate:"required,sku"`
	CreatedOn   string  `json:"-"` // completely omit the created on field from the output
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

// add a method to product struct and call it fromJSON, source of this will be an io reader and will return error if there is an issue
func (p *Product) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r) // passing reader into new decoder
	return e.Decode(p)      // decodes the next input value and will store it in the reference pointed by p (destination structure)
	// the product struct whenever called from reader, it will be encoded
	// create the product in product handler
}

//adding a validator method on our data type
func (p *Product) Validate() error {
	validate := validator.New()                     // creating a validator
	validate.RegisterValidation("sku", validateSKU) // registering a new validation functions
	validate.Struct(p)                              //validating the struct p, pass as reference, struct returns an error message
	return nil
}

// creating a new validator function
func validateSKU(fl validator.FieldLevel) bool {
	//sku is in the format abc-xyz-pqr
	// to ensure the format, use regexp
	re := regexp.MustCompile(`[a-z]+-[a-z]+-[a-z]+`)
	matches := re.FindAllString(fl.Field().String(), -1) //fl is validator here //this will return a slice of string
	if len(matches) != 1 {
		return false
	}
	return true
}

// Products is a collection of product
type Products []*Product //rather than returning slice of product via method, create a slice of product only

// add a method to this slice same way we do it for a struct
// add another method called ToJSON, takes i/p as iowriter and returns err.
// Now put the whole encoded logic in a seperate method
func (p *Products) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w) //create encoder e and pass the writer w in it
	return e.Encode(p)      //encode will return an error, method has an error as a response
}

// GetProducts returns a list of product
func GetProducts() Products { // this method acts as a DAM
	return ProductList
}

func AddProduct(p *Product) { // pass in a product into new function called add product
	// autogenerate id for these in the next function
	(p).ID = getNextID()                 // get the next product ID
	ProductList = append(ProductList, p) // add the new product to our list
}

func UpdateProduct(id int, p *Product) error { // pass in a product into new function called update product
	// does the product exist?
	_, pos, err := findProduct(id)
	if err != nil {
		return err
	}
	// update the product
	(p).ID = id
	// update the product list
	ProductList[pos] = p // find product results the position so we can replace it easily
	return nil
}

var ErrProductNotFound = fmt.Errorf("Product not found") // create a structured error, its important

// need to find the product first
func findProduct(id int) (*Product, int, error) {
	for i, p := range ProductList {
		if p.ID == id {
			return p, i, nil
		}
	}
	return nil, -1, ErrProductNotFound
}

func getNextID() int { // returns an int
	lp := ProductList[len(ProductList)-1] // to get the id of the last product in the list of products
	return lp.ID + 1
}

var ProductList = []*Product{
	&Product{
		ID:          1,
		Name:        "latte",
		Description: "Frothy Milk Coffee",
		Price:       2.45,
		SKU:         "abc323", //Product ID kind of
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
	&Product{
		ID:          2,
		Name:        "Espresso",
		Description: "Short and strong coffee without milk",
		Price:       1.99,
		SKU:         "fjd34",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
}
