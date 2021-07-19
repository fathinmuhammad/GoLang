package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
)

var db *gorm.DB
var err error

//Product is a representation of a product
type Product struct {
	ID    int             `json:"ID"`
	Code  string          `json:"Code"`
	Name  string          `json:"Name"`
	Price decimal.Decimal `json:"Price" sql:"type:decimal(16,2)"`
}

//Result an array of product
type Result struct {
	Code    int         `json:"Code"`
	Data    interface{} `json:"Data"`
	Message string      `json:"Message"`
}

func main() {
	db, err = gorm.Open("mysql", "root:@/go_rest_api_crud?charset=utf8&parseTime=true")

	if err != nil {
		log.Println("Connection Failed", err)
	} else {
		log.Println("Connection Succeesed")
	}

	db.AutoMigrate(&Product{})

	handleRequests()
}

func handleRequests() {
	log.Println("Start the development server at http://localhost:9000")

	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/api/products", getListProducts).Methods("GET")
	myRouter.HandleFunc("/api/product/{id}", getProduct).Methods("GET")
	myRouter.HandleFunc("/api/product/create", createdProduct).Methods("POST")
	myRouter.HandleFunc("/api/product/{id}", updateProduct).Methods("PUT")
	myRouter.HandleFunc("/api/product/{id}", deleteProduct).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":9000", myRouter))
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome !")
}

//Create a Product
func createdProduct(w http.ResponseWriter, r *http.Request) {
	payload, _ := ioutil.ReadAll(r.Body)

	var product Product
	json.Unmarshal(payload, &product)

	db.Create(&product)

	res := Result{Code: 200, Data: product, Message: "Succesed Created Product"}
	result, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

//Get List Products
func getListProducts(w http.ResponseWriter, r *http.Request) {
	products := []Product{}

	db.Find(&products)

	res := Result{Code: 200, Data: products, Message: "Succesed Get Product"}
	results, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(results)
}

// Get Product by ID
func getProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["id"]

	var product Product
	db.First(&product, productID)

	res := Result{Code: 200, Data: product, Message: "Succesed Get Product"}
	result, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

//Updated Product
func updateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["id"]

	payload, _ := ioutil.ReadAll(r.Body)

	var productUpdates Product
	json.Unmarshal(payload, &productUpdates)

	var product Product
	db.First(&product, productID)
	db.Model(&product).Updates(productUpdates)

	res := Result{Code: 200, Data: product, Message: "Succesed Updated Product"}
	result, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

//Delete Product
func deleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["id"]

	var product Product
	db.First(&product, productID)
	db.Delete(&product)

	res := Result{Code: 200, Message: "Succesed Delete Product"}
	result, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}
