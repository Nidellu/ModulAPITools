package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	m "modulgo/model"
)

func GetAllProducts(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	query := "SELECT * FROM products"
	name := r.URL.Query()["name"]
	price := r.URL.Query()["price"]

	if name != nil {
		fmt.Println(name[0])
		query += " WHERE name= '" + name[0] + "'"
	}

	if price != nil {
		if name[0] != "" {
			query += "AND"
		} else {
			query += "WHERE"
		}
		query += " price= '" + price[0] + "'"
	}

	rows, err := db.Query(query)
	if err != nil {
		sendProductErrorResponse(w, "Error: Database not found")
		return
	}

	var product m.Product
	var products []m.Product
	for rows.Next() {
		if err := rows.Scan(&product.ID, &product.Name, &product.Price); err != nil {
			sendProductErrorResponse(w, "Data not found")
			return
		} else {
			products = append(products, product)
		}
	}

	if len(products) == 0 {
		sendProductErrorResponse(w, "Data not found")
		return
	}

	sendProductSuccessResponse(w, "Success get products", products)
}

func GetAllProductssGorm(w http.ResponseWriter, r *http.Request) {
	db := connectGorm()
	defer db.Close()

	var product []m.Product
	result := db.Find(&product)

	if result.Error != nil {
		sendProductErrorResponse(w, "Error")
	} else {
		sendProductSuccessResponse(w, "Success", product)
	}
}

func InsertNewProduct(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		sendProductErrorResponse(w, "Error: Parsing data")
		return
	}

	name := r.Form.Get("name")
	price, _ := strconv.Atoi(r.Form.Get("price"))

	if name == "" || price == 0 {
		sendProductErrorResponse(w, "Bad Request: Incomplete input data")
		return
	}

	data, err := db.Begin()
	if err != nil {
		sendProductErrorResponse(w, "Error: Database not found")
		return
	}
	defer data.Rollback()

	_, errQuery := db.Exec("INSERT INTO products (name, price) VALUES (?, ?)", name, price)
	if errQuery == nil {
		sendProductSuccessResponse(w, "Success insert data", nil)
	} else {
		sendProductErrorResponse(w, "Error insert data")
	}
}

func InsertNewProductGorm(w http.ResponseWriter, r *http.Request) {
	db := connectGorm()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		log.Println("Error parsing form data:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	name := r.Form.Get("name")
	price, _ := strconv.Atoi(r.Form.Get("price"))

	if name == "" || price == 0 {
		log.Println("Error: Incomplete data provided")
		http.Error(w, "Bad Request: Incomplete data", http.StatusBadRequest)
		return
	}

	var user m.User
	result := db.Raw("INSERT INTO products (name, price) VALUES (?, ?)", name, price).Scan(&user)

	if result.Error != nil {
		sendSuccessResponse(w)
	} else {
		sendErrorResponse(w)
	}
}

func PutProduct(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		return
	}

	prodID := r.URL.Query().Get("id")

	if prodID == "" {
		sendProductErrorResponse(w, "Bad Request: ID missing")
		return
	}

	name := r.Form.Get("name")
	price, _ := strconv.Atoi(r.Form.Get("price"))

	if name == "" || price == 0 {
		sendProductErrorResponse(w, "Bad Request: Incomplete data")
		return
	}

	data, err := db.Begin()
	if err != nil {
		sendProductErrorResponse(w, "Database not found")
		return
	}
	defer data.Rollback()

	_, errQuery := db.Exec("UPDATE products SET name = ?, price = ? WHERE id = ?", name, price, prodID)

	if errQuery == nil {
		sendProductSuccessResponse(w, "Success update data", nil)
	} else {
		sendProductErrorResponse(w, "Error update data")
	}
}

func PutProductGorm(w http.ResponseWriter, r *http.Request) {
	db := connectGorm()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		return
	}

	name := r.Form.Get("name")
	price, _ := strconv.Atoi(r.Form.Get("price"))

	if name == "" || price == 0 {
		sendProductErrorResponse(w, "Bad Request: Incomplete data")
		return
	}

	product := m.Product{Name: name, Price: price}
	result := db.Save(&product)

	if result.Error != nil {
		sendUserErrorResponse(w, "Error")
	} else {
		sendUserSuccessResponse(w, "Success", nil)
	}
}

func DeleteSingleProduct(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	prodID := r.URL.Query().Get("id")

	if prodID == "" {
		sendProductErrorResponse(w, "Bad Request: ID missing")
		return
	} else {
		data, err := db.Begin()
		if err != nil {
			sendProductErrorResponse(w, "Database not found")
			return
		}
		defer data.Rollback()

		_, err = db.Exec("DELETE FROM transactions WHERE product_id = ?", prodID)
		if err != nil {
			sendProductErrorResponse(w, "Delete failed")
			return
		}

		_, err = db.Exec("DELETE FROM products WHERE id = ?", prodID)
		if err != nil {
			sendProductErrorResponse(w, "Delete failed")
		}
		sendSuccessResponse(w)
	}
}

func DeleteSingleProductGorm(w http.ResponseWriter, r *http.Request) {
	db := connectGorm()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		return
	}

	prodID := r.URL.Query().Get("id")
	if prodID == "" {
		sendUserErrorResponse(w, "Bad request: Missing ID")
		return
	}

	userID, err := strconv.Atoi(prodID)
	if err != nil {
		sendUserErrorResponse(w, "Bad request: Invalid ID")
		return
	}

	user := m.User{ID: userID}
	result := db.Delete(user)

	if result.Error != nil {
		sendUserErrorResponse(w, "Error")
	} else {
		sendUserSuccessResponse(w, "Success", nil)
	}
}

func sendProductSuccessResponse(w http.ResponseWriter, message string, products []m.Product) {
	w.Header().Set("Content-Type", "application/json")
	var response m.ProductsResponse
	response.Status = 200
	response.Message = message
	response.Data = products
	json.NewEncoder(w).Encode(response)
}
func sendProductErrorResponse(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	var response m.UserResponse
	response.Status = 400
	response.Message = message
	json.NewEncoder(w).Encode(response)
}
