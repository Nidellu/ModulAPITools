package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	m "modulgo/model"
)

func GetAllTransactions(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	query := "SELECT * FROM transactions"

	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		return
	}

	var transaction m.Transaction
	var transactions []m.Transaction
	for rows.Next() {
		if err := rows.Scan(&transaction.ID, &transaction.UserID, &transaction.ProductID, &transaction.Quantity); err != nil {
			log.Println(err)
			return
		} else {
			transactions = append(transactions, transaction)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	var response m.TransactionsResponse
	response.Status = 200
	response.Message = "Success"
	response.Data = transactions
	json.NewEncoder(w).Encode(response)
}

func GetDetailUserTransactions(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	query := "SELECT t.id, u.id, u.name, u.age, u.address, p.id, p.name, p.price, t.quantity FROM transactions t JOIN users u ON t.user_id = u.id JOIN products p ON t.product_id = p.id"

	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		sendErrorResponse(w)
		return
	}

	var detailTransaction m.TransactionsDetail
	var detailTransactions []m.TransactionsDetail
	for rows.Next() {
		if err := rows.Scan(&detailTransaction.ID, &detailTransaction.User.ID, &detailTransaction.User.Name, &detailTransaction.User.Age, &detailTransaction.User.Address,
			&detailTransaction.Product.ID, &detailTransaction.Product.Name, &detailTransaction.Product.Price, &detailTransaction.Quantity); err != nil {
			log.Println(err)
			sendErrorResponse(w)
			return
		} else {
			detailTransactions = append(detailTransactions, detailTransaction)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	var response m.TransactionsDetailResponse
	response.Status = 200
	response.Data.Transaction = detailTransactions
	response.Message = "Success"
	json.NewEncoder(w).Encode(response)
}

func GetDetailUserTransactionsByID(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	userID := r.URL.Query().Get("id")
	query := "SELECT t.id, u.id, u.name, u.age, u.address, p.id, p.name, p.price, t.quantity FROM transactions t JOIN users u ON t.user_id = u.id JOIN products p ON t.product_id = p.id WHERE u.id = ?"

	rows, err := db.Query(query, userID)
	if err != nil {
		log.Println(err)
		sendErrorResponse(w)
		return
	}

	var detailTransaction m.TransactionsDetail
	var detailTransactions []m.TransactionsDetail
	for rows.Next() {
		if err := rows.Scan(&detailTransaction.ID, &detailTransaction.User.ID, &detailTransaction.User.Name, &detailTransaction.User.Age, &detailTransaction.User.Address,
			&detailTransaction.Product.ID, &detailTransaction.Product.Name, &detailTransaction.Product.Price, &detailTransaction.Quantity); err != nil {
			log.Println(err)
			sendErrorResponse(w)
			return
		} else {
			detailTransactions = append(detailTransactions, detailTransaction)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	var response m.TransactionsDetailResponse
	response.Status = 200
	response.Data.Transaction = detailTransactions
	response.Message = "Success"
	json.NewEncoder(w).Encode(response)
}

func InsertNewTransaction(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		log.Println("Error parsing form data:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	userID, _ := strconv.Atoi(r.Form.Get("user_id"))
	productID, _ := strconv.Atoi(r.Form.Get("product_id"))
	quantity, _ := strconv.Atoi(r.Form.Get("quantity"))

	if userID == 0 || productID == 0 || quantity == 0 {
		log.Println("Error: Incomplete data provided")
		http.Error(w, "Bad Request: Incomplete data", http.StatusBadRequest)
		return
	} else {
		data, err := db.Begin()
		if err != nil {
			log.Println("Error database not found:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer data.Rollback()

		var none int
		err = db.QueryRow("SELECT EXISTS (SELECT 1 FROM products WHERE id = ?)", productID).Scan(&none)
		if err != nil {
			http.Error(w, "Data Missing", http.StatusBadRequest)
			return
		}

		if none == 0 {
			_, err = db.Exec("INSERT INTO products (id, name, price) VALUES (?, '', 0)", productID)
			if err != nil {
				http.Error(w, "Insert failed", http.StatusBadRequest)
				return
			}
		}

		_, err = db.Exec("INSERT INTO transactions (user_id, product_id, quantity) VALUES (?, ?, ?)", userID, productID, quantity)
		if err != nil {
			http.Error(w, "Insert failed", http.StatusBadRequest)
			return
		}

		sendSuccessResponse(w)
	}

}

func PutTransaction(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		return
	}

	transID := r.URL.Query().Get("id")

	if transID == "" {
		log.Println("Error: ID missing")
		http.Error(w, "Bad Request: ID missing", http.StatusBadRequest)
		return
	}

	userID, _ := strconv.Atoi(r.Form.Get("user_id"))
	productID, _ := strconv.Atoi(r.Form.Get("product_id"))
	quantity, _ := strconv.Atoi(r.Form.Get("quantity"))

	if userID == 0 || productID == 0 || quantity == 0 {
		log.Println("Error: Incomplete data provided")
		http.Error(w, "Bad Request: Incomplete data", http.StatusBadRequest)
		return
	}

	data, err := db.Begin()
	if err != nil {
		log.Println("Error database not found:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer data.Rollback()

	_, errQuery := db.Exec("UPDATE transactions SET user_id = ?, product_id = ?, quantity = ? WHERE id = ?", userID, productID, quantity, transID)
	if errQuery != nil {
		http.Error(w, "Update failed", http.StatusBadRequest)
		return
	}

	if errQuery == nil {
		sendSuccessResponse(w)
	} else {
		sendErrorResponse(w)
	}
}

func DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	transID := r.URL.Query().Get("id")

	if transID == "" {
		log.Println("Error: ID missing")
		http.Error(w, "Bad Request: ID missing", http.StatusBadRequest)
		return
	}

	data, err := db.Begin()
	if err != nil {
		log.Println("Error database not found:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer data.Rollback()

	_, errQuery := db.Exec("DELETE FROM transactions WHERE id = ?", transID)
	if errQuery != nil {
		http.Error(w, "Delete failed", http.StatusBadRequest)
		return
	}

	if errQuery == nil {
		sendSuccessResponse(w)
	} else {
		sendErrorResponse(w)
	}
}

func sendSuccessResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	var response m.UserResponse
	response.Status = 200
	response.Message = "message"
	json.NewEncoder(w).Encode(response)
}
func sendErrorResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	var response m.UserResponse
	response.Status = 400
	response.Message = "message"
	json.NewEncoder(w).Encode(response)
}
