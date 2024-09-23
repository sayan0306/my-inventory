package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (app *App) Initialize() error {
	connectionStr := fmt.Sprintf("%v:%v@tcp(127.0.0.1:3306)/%v", DBUser, DBPassword, DBName)
	var err error
	app.DB, err = sql.Open("mysql", connectionStr)
	checkErr(err)

	app.Router = mux.NewRouter().StrictSlash(true)
	app.handleRoutes()
	return nil
}

func (app *App) Run(address string) {
	log.Fatal(http.ListenAndServe(address, app.Router))
}

func (app *App) handleRoutes() {
	app.Router.HandleFunc("/products", app.getProducts).Methods(http.MethodGet)
	app.Router.HandleFunc("/products/{id}", app.getProductById).Methods(http.MethodGet)
	app.Router.HandleFunc("/products", app.createProduct).Methods(http.MethodPost)
	app.Router.HandleFunc("/products/{id}", app.updateProductById).Methods(http.MethodPut)
	app.Router.HandleFunc("/products/{id}", app.deleteProductById).Methods(http.MethodDelete)
}

func (app *App) getProducts(w http.ResponseWriter, r *http.Request) {
	products, err := getProducts(app.DB)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendResponse(w, http.StatusOK, products)
}

func (app *App) getProductById(w http.ResponseWriter, r *http.Request) {
	var err error
	var id int
	id, err = strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid Product ID")
		return
	}

	p := Product{Id: id}
	err = p.getProductById(app.DB)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			sendError(w, http.StatusNotFound, "Product not found")
		default:
			sendError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	sendResponse(w, http.StatusOK, p)

}

func (app *App) createProduct(w http.ResponseWriter, r *http.Request) {
	var p Product
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	err = p.createProduct(app.DB)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
	}
	sendResponse(w, http.StatusCreated, p)
}

func (app *App) updateProductById(w http.ResponseWriter, r *http.Request) {
	var err error
	var id int
	id, err = strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid Product ID")
		return
	}

	var p Product
	err = json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	p.Id = id
	err = p.updateProduct(app.DB)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendResponse(w, http.StatusOK, p)
}

func (app *App) deleteProductById(w http.ResponseWriter, r *http.Request) {
	var err error
	var id int
	id, err = strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid Product ID")
		return
	}
	p := Product{Id: id}
	err = p.deleteProduct(app.DB)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendResponse(w, http.StatusOK, map[string]string{"result": "successfully deleted product"})
}

func sendResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	response, _ := json.Marshal(data)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(response)
}

func sendError(w http.ResponseWriter, statusCode int, err string) {
	errorMessage := map[string]string{"error": err}
	sendResponse(w, statusCode, errorMessage)
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
}
