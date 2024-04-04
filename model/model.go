package model

import "github.com/dgrijalva/jwt-go"

type User struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Age       int    `json:"age"`
	Address   string `json:"address"`
	Passwords string `json:"Passwords"`
	Email     string `json:"email"`
	UserType  int    `json:"user_type"`
}

type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

type Transaction struct {
	ID        int `json:"id"`
	UserID    int `json:"userId"`
	ProductID int `json:"productId"`
	Quantity  int `json:"quantity"`
}

type UserResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    User   `json:"data"`
}

type UsersResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    []User `json:"data"`
}

type ProductResponse struct {
	Status  int     `json:"status"`
	Message string  `json:"message"`
	Data    Product `json:"data"`
}

type ProductsResponse struct {
	Status  int       `json:"status"`
	Message string    `json:"message"`
	Data    []Product `json:"data"`
}

type TransactionResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    Transaction `json:"data"`
}

type TransactionsResponse struct {
	Status  int           `json:"status"`
	Message string        `json:"message"`
	Data    []Transaction `json:"data"`
}

type TransactionsDetail struct {
	ID       int     `json:"ID"`
	User     User    `json:"user"`
	Product  Product `json:"product"`
	Quantity int     `json:"quantity"`
}

type TransactionsDetailResponse struct {
	Status  int                        `json:"status"`
	Message string                     `json:"message"`
	Data    TransactionDetailResponses `json:"Data"`
}

type TransactionDetailResponses struct {
	Transaction []TransactionsDetail `json:"Transactions"`
}

type Claims struct {
	ID       int    `json: "id"`
	Name     string `json:"name"`
	UserType int    `json:"user_type"`
	jwt.StandardClaims
}
