package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	m "modulgo/model"
)

func CheckUserLogin(w http.ResponseWriter, r *http.Request) {

	db := connect()
	defer db.Close()

	email := r.URL.Query().Get("email")
	password := r.URL.Query().Get("password")

	if email == "" || password == "" {
		sendUserErrorResponse(w, "Bad request: Incomplete input data")
		return
	}

	query := "SELECT * FROM users WHERE email = ?"

	var user m.User
	row := db.QueryRow(query, email)
	if err := row.Scan(&user.ID, &user.Name, &user.Age, &user.Address, &user.Passwords, &user.Email, &user.UserType); err != nil {
		sendUserErrorResponse(w, "Error: Data not found")
		return
	}

	if password != user.Passwords {
		sendUserErrorResponse(w, "Error: Incorrect password")
		return
	}

	startScheduler()
	generateToken(w, user.ID, user.Name, user.UserType)
	sendSuccessResponse(w)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	resetUserToken(w)

	var response m.UserResponse
	response.Status = 200
	response.Message = "Success"

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetAllUsersGorm(w http.ResponseWriter, r *http.Request) {
	db := connectGorm()
	defer db.Close()
	var users []m.User
	users = GetUsers()
	if users == nil {
		result := db.Find(&users)

		if result.Error != nil {
			sendUserErrorResponse(w, "Error")

		}
		SetUsers(users)
	}
	sendUserSuccessResponse(w, "Success", users)
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	query := "SELECT * FROM users"
	name := r.URL.Query()["name"]
	age := r.URL.Query()["age"]

	if name != nil {
		fmt.Println(name[0])
		query += " WHERE name= '" + name[0] + "'"
	}

	if age != nil {
		if name[0] != "" {
			query += "AND"
		} else {
			query += "WHERE"
		}
		query += " age= '" + age[0] + "'"
	}

	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		return
	}

	var user m.User
	var users []m.User
	for rows.Next() {
		if err := rows.Scan(&user.ID, &user.Name, &user.Age, &user.Address, &user.Passwords, &user.Email, &user.UserType); err != nil {
			log.Println(err)
			return
		} else {
			users = append(users, user)
		}
	}

	if len(users) == 0 {
		sendUserErrorResponse(w, "Data not Found")
		return
	}

	sendUserSuccessResponse(w, "Success", users)
}

func InsertNewUser(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		sendUserErrorResponse(w, "Error: error parsing data")
		return
	}

	name := r.Form.Get("name")
	age, _ := strconv.Atoi(r.Form.Get("age"))
	address := r.Form.Get("address")
	passwords := r.Form.Get("password")
	email := r.Form.Get("email")

	if name == "" || age == 0 || address == "" || passwords == "" || email == "" {
		sendUserErrorResponse(w, "Bad request: Incomplete Data!")
		return
	}

	data, err := db.Begin()
	if err != nil {
		sendUserErrorResponse(w, "Error: Database not Found!")
		return
	}
	defer data.Rollback()

	_, errQuery := db.Exec("INSERT INTO users (name, age, address, Passwords, email, usertype) VALUES (?, ?, ?, ?, ?, 2)", name, age, address, passwords, email)
	if errQuery != nil {
		sendUserErrorResponse(w, "Error: Failed to insert data")
		return
	}

	go func() {
		SendEmail("Berhasil Insert")
	}()

	if errQuery == nil {
		sendUserSuccessResponse(w, "Success", nil)
	} else {
		sendUserErrorResponse(w, "Failed")
	}
}

func InsertNewUserGorm(w http.ResponseWriter, r *http.Request) {
	db := connectGorm()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		sendUserErrorResponse(w, "Error: error parsing data")
		return
	}

	name := r.Form.Get("name")
	age, _ := strconv.Atoi(r.Form.Get("age"))
	address := r.Form.Get("address")
	passwords := r.Form.Get("password")
	email := r.Form.Get("email")

	if name == "" || age == 0 || address == "" || passwords == "" || email == "" {
		sendUserErrorResponse(w, "Bad request: Incomplete Data!")
		return
	}

	user := m.User{Name: name, Age: age, Address: address, Passwords: passwords, Email: email}
	result := db.Create(&user)

	if result.Error != nil {
		sendUserErrorResponse(w, "Error")
	} else {
		sendUserSuccessResponse(w, "Success", nil)
	}
}

func PutUser(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		return
	}

	userID := r.URL.Query().Get("id")

	if userID == "" {
		sendUserErrorResponse(w, "Bad request: Missing ID input")
		return
	}

	name := r.Form.Get("name")
	age, _ := strconv.Atoi(r.Form.Get("age"))
	address := r.Form.Get("address")

	if name == "" || age == 0 || address == "" {
		sendUserErrorResponse(w, "Bad request: Incomplete Data!")
		return
	}

	data, err := db.Begin()
	if err != nil {
		sendUserErrorResponse(w, "Error: Database not Found")
		return
	}
	defer data.Rollback()

	_, errQuery := db.Exec("UPDATE users SET name = ?, age = ?, address = ? WHERE id = ?", name, age, address, userID)
	if errQuery != nil {
		sendUserErrorResponse(w, "Error: Failed to update data")
		return
	}

	if errQuery == nil {
		sendUserSuccessResponse(w, "Success", nil)
	} else {
		sendUserErrorResponse(w, "Failed")
	}

}

func PutUserGorm(w http.ResponseWriter, r *http.Request) {
	db := connectGorm()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		return
	}

	getID := r.URL.Query().Get("id")
	if getID == "" {
		sendUserErrorResponse(w, "Bad request: Missing ID")
		return
	}

	userID, err := strconv.Atoi(getID)
	if err != nil {
		sendUserErrorResponse(w, "Bad request: Invalid ID")
		return
	}

	name := r.Form.Get("name")
	age, _ := strconv.Atoi(r.Form.Get("age"))
	address := r.Form.Get("address")
	passwords := r.Form.Get("password")
	email := r.Form.Get("email")

	if name == "" || age == 0 || address == "" || passwords == "" || email == "" {
		sendUserErrorResponse(w, "Bad request: Incomplete Data!")
		return
	}

	user := m.User{ID: userID, Name: name, Age: age, Address: address, Passwords: passwords, Email: email}
	result := db.Save(&user)

	if result.Error != nil {
		sendUserErrorResponse(w, "Error")
	} else {
		sendUserSuccessResponse(w, "Success", nil)
	}
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	userID := r.URL.Query().Get("id")

	if userID == "" {
		sendUserErrorResponse(w, "Bad request: Missing ID")
		return
	}

	data, err := db.Begin()
	if err != nil {
		sendUserErrorResponse(w, "Error: Database not found")
		return
	}
	defer data.Rollback()

	_, errQuery := db.Exec("DELETE FROM users WHERE id = ?", userID)
	if errQuery != nil {
		sendUserErrorResponse(w, "Error: Failed to delete")
		return
	}

	if errQuery == nil {
		sendUserSuccessResponse(w, "Success", nil)
	} else {
		sendUserErrorResponse(w, "Failed")
	}
}

func DeleteUserGorm(w http.ResponseWriter, r *http.Request) {
	db := connectGorm()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		return
	}

	getID := r.URL.Query().Get("id")
	if getID == "" {
		sendUserErrorResponse(w, "Bad request: Missing ID")
		return
	}

	userID, err := strconv.Atoi(getID)
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
