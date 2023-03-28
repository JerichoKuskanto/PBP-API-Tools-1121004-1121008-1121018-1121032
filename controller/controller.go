package controller

import (
	"PBP-API-Tools-1121004-1121008-1121018-1121032/model"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	query := "SELECT id,name,age,address,email FROM users"

	name := r.URL.Query()["name"]

	if name != nil {
		query += " WHERE name = '" + name[0] + "'"
	}

	rows, err := db.Query(query)

	if err != nil {
		log.Println(err)
		sendErrorResponse(w, "Something went wrong, please try again")
		return
	}

	var user model.User
	var users []model.User

	for rows.Next() {
		if err := rows.Scan(&user.ID, &user.Name, &user.Age, &user.Address, &user.Email); err != nil {
			log.Println(err)
			sendErrorResponse(w, "Something went wrong, please try again")
			return
		} else {
			user.Password = "********"
			users = append(users, user)
		}
	}

	var response model.UsersResponse
	response.Status = 200
	response.Message = "Success"
	response.Data = users
	w.Header().Set("Content=Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func InsertUser(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		sendErrorResponse(w, "Something went wrong, please try again")
		return
	}
	name := r.Form.Get("name")
	age, _ := strconv.Atoi(r.Form.Get("age"))
	address := r.Form.Get("address")

	res, errQuery := db.Exec("INSERT INTO users(name,age,address) VALUES (?,?,?)", name, age, address)
	id, _ := res.LastInsertId()

	var response model.UsersResponse
	if errQuery == nil {
		var user model.User
		var users []model.User
		response.Status = 200
		response.Message = "Insert Success"
		id := int(id)
		user.Name = name
		user.Age = age
		user.Address = address
		user.Password = "********"
		users = append(users, model.User{ID: id, Name: name, Age: age, Address: address, Email: "", Password: "********"})
		response.Data = users
	} else {
		response.Status = 400
		response.Message = "Insert Failed"
	}
	w.Header().Set("Content=Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		sendErrorResponse(w, "Something went wrong, please try again")
		return
	}
	vars := mux.Vars(r)
	userId := vars["user_id"]
	name := r.Form.Get("name")
	age, _ := strconv.Atoi(r.Form.Get("age"))
	address := r.Form.Get("address")
	_, errQuery := db.Exec("UPDATE users SET name = ?, age = ?, address = ? WHERE id = ?", name, age, address, userId)

	var response model.UsersResponse
	if errQuery == nil {
		var users []model.User
		response.Status = 200
		response.Message = "Update Success"
		id, _ := strconv.Atoi(userId)
		users = append(users, model.User{ID: id, Name: name, Age: age, Address: address, Email: "", Password: "********"})
		response.Data = users
	} else {
		response.Status = 400
		response.Message = "Update Failed"
	}
	w.Header().Set("Content=Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		sendErrorResponse(w, "Something went wrong, please try again")
		return
	}
	vars := mux.Vars(r)
	userId := vars["user_id"]
	_, errQuery := db.Exec("DELETE FROM users WHERE id = ?", userId)

	var response model.UsersResponse
	if errQuery == nil {
		response.Status = 200
		response.Message = "Success"
	} else {
		response.Status = 400
		response.Message = "Delete Failed"
	}
	w.Header().Set("Content=Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
func LoginUser(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		sendErrorResponse(w, "Something went wrong, please try again")
		return
	}
	loginSuccess := false
	platform := r.Header.Get("platform")
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	query := "SELECT id,name,age,address,email,password,usertype FROM users"
	rows, err := db.Query(query)

	var user model.User

	for rows.Next() {
		var usertype int
		if err := rows.Scan(&user.ID, &user.Name, &user.Age, &user.Address, &user.Email, &user.Password, &usertype); err != nil {
			log.Println(err)
			sendErrorResponse(w, "Something went wrong, please try again")
			return
		} else {
			if email == user.Email && password == user.Password {
				loginSuccess = true
				generateToken(w, user.ID, user.Name, usertype)
				break
			}
		}
	}

	var response model.UsersResponse
	if err == nil && loginSuccess {
		response.Status = 200
		response.Message = "Success login from " + platform
		var users []model.User
		users = append(users, model.User{ID: user.ID, Name: user.Name, Age: user.Age, Address: user.Password, Email: user.Email, Password: "********"})
		response.Data = users

	} else {
		response.Status = 400
		response.Message = "Login failed!"
	}
	w.Header().Set("Content=Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

func LogoutUser(w http.ResponseWriter, r *http.Request) {
	resetUserToken(w)
	var response model.UsersResponse
	response.Status = 200
	response.Message = "Success"
	w.Header().Set("Content=Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func sendUnauthorizedResponse(w http.ResponseWriter) {
	var response model.ErrorResponse
	response.Status = 401
	response.Message = "Unauthorized Access"
	w.Header().Set("Content=Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
