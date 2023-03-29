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

func sendUnauthorizedResponse(w http.ResponseWriter) {
	var response model.Response
	response.Status = 401
	response.Message = "Unauthorized Access"
	w.Header().Set("Content=Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
