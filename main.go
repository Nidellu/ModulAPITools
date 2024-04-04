package main

import (
	"fmt"
	"log"
	"modulgo/controller"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/v1/users/admin", controller.Authenticate(controller.GetAllUsers, 1)).Methods("GET")
	router.HandleFunc("/v1/users", controller.GetAllUsers).Methods("GET")
	router.HandleFunc("/v2/users", controller.GetAllUsersGorm).Methods("GET")
	router.HandleFunc("/v1/users", controller.InsertNewUser).Methods("POST")
	router.HandleFunc("/v2/users", controller.InsertNewUserGorm).Methods("POST")
	router.HandleFunc("/v1/users", controller.PutUser).Methods("PUT")
	router.HandleFunc("/v2/users", controller.PutUserGorm).Methods("PUT")
	router.HandleFunc("/v1/users", controller.DeleteUser).Methods("DELETE")
	router.HandleFunc("/v2/users", controller.DeleteUserGorm).Methods("DELETE")

	router.HandleFunc("/v1/products", controller.GetAllProducts).Methods("GET")
	router.HandleFunc("/v2/products", controller.GetAllProductssGorm).Methods("GET")
	router.HandleFunc("/v1/products", controller.InsertNewProduct).Methods("POST")
	router.HandleFunc("/v2/products", controller.InsertNewProductGorm).Methods("POST")
	router.HandleFunc("/v1/products", controller.PutProduct).Methods("PUT")
	router.HandleFunc("/v1/products", controller.DeleteSingleProduct).Methods("DELETE")

	router.HandleFunc("/v1/transactions", controller.GetAllTransactions).Methods("GET")
	router.HandleFunc("/v1/transactions", controller.InsertNewTransaction).Methods("POST")
	router.HandleFunc("/v1/transactions", controller.PutTransaction).Methods("PUT")
	router.HandleFunc("/v1/transactions", controller.DeleteTransaction).Methods("DELETE")

	router.HandleFunc("/v1/usertransactions", controller.GetDetailUserTransactions).Methods("GET")
	router.HandleFunc("/v1/usertransactionsID", controller.GetDetailUserTransactionsByID).Methods("GET")

	router.HandleFunc("/login", controller.CheckUserLogin).Methods("POST")
	router.HandleFunc("/reset", controller.Logout).Methods("DELETE")

	http.Handle("/", router)
	fmt.Println("Connected to port 8888")
	log.Println("Connected to port 8888")
	log.Fatal(http.ListenAndServe(":8888", router))
}
