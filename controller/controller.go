package controller

import (
	"PBP-API-Tools-1121004-1121008-1121018-1121032/model"
	"encoding/json"
	"log"
	"net/http"

	"gopkg.in/gomail.v2"
)

func sendMail(receiver string, usertype int) bool {
	m := gomail.NewMessage()
	//Buat header untuk email
	m.SetHeader("From", "articler8375@gmail.com")
	m.SetHeader("To", receiver)
	m.SetHeader("Subject", "New Article Published")
	//Cek tipe user yang dikirim
	//Jika 1(normal), mereka diberikan article biasa
	//Jika 2(premium), mereka diberikan article eksklusif
	if usertype == 1 {
		m.SetBody("text/html", "<h2>Title Goes Here</h2><a>Normal link goes here</a><p>Summary Here</p>")
	} else {
		m.SetBody("text/html", "<h2>Title Goes Here</h2><a>Premium link goes here</a><p>Summary Here</p>")
	}
	d := gomail.NewDialer("smtp.gmail.com", 587, "articler8375@gmail.com", "1234")
	if err := d.DialAndSend(m); err != nil {
		log.Println(err)
		return false
	} else {
		log.Println("Email sent to ", receiver)
		return true
	}
}

func SendNotificationEmail(w http.ResponseWriter, r *http.Request) {
	sendSuccess := sendMail("exampleuser1@gmail.com", 1)

	if sendSuccess {
		var response model.Response
		response.Status = 200
		response.Message = "Email Sent"
		w.Header().Set("Content=Type", "application/json")
		json.NewEncoder(w).Encode(response)

	} else {
		var response model.Response
		response.Status = 400
		response.Message = "Fail to send email"
		w.Header().Set("Content=Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
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
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	query := "SELECT id,name,age,address,email,password,usertype FROM users"
	rows, err := db.Query(query)

	var user model.User

	for rows.Next() {
		var usertype int
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &usertype); err != nil {
			log.Println(err)
			sendErrorResponse(w, "Something went wrong, please try again")
			return
		} else {
			if email == user.Email && password == user.Password {
				loginSuccess = true
				generateToken(w, user.ID, user.Username, usertype)
				break
			}
		}
	}

	var response model.UsersResponse
	if err == nil && loginSuccess {
		response.Status = 200
		response.Message = "Success login"
		var users []model.User
		users = append(users, model.User{ID: user.ID, Username: user.Username, Email: user.Email, Password: user.Password})
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
	var response model.Response
	response.Status = 401
	response.Message = "Unauthorized Access"
	w.Header().Set("Content=Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
