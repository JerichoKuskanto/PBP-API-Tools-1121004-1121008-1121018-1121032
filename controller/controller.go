package controller

import (
	"PBP-API-Tools-1121004-1121008-1121018-1121032/model"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/redis/go-redis/v9"
	"gopkg.in/gomail.v2"
)

func sendMail(receiver string, usertype int) {
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
		log.Println("Failed to send email to", receiver)
	} else {
		log.Println("Email sent to", receiver)
	}
}

var ctx = context.Background()

func mailSending(receiver []string, usertype int) {
	for i := range receiver {
		//time.Sleep(100 * time.Millisecond)
		sendMail(receiver[i], usertype)
	}
}

func SendNotificationEmail(w http.ResponseWriter, r *http.Request) {

}
func setAllUserRedis() {
	var users []model.User

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	users = getAllUser()
	rdb.Del(ctx, "users")
	for i, v := range users {
		if err := rdb.HSet(ctx, "user"+strconv.Itoa(i), v).Err(); err != nil {
			panic(err)
		}
		rdb.Expire(ctx, "user"+strconv.Itoa(i), 60*time.Minute)
		rdb.SAdd(ctx, "users", "user"+strconv.Itoa(i))
	}
}
func getAllUserRedis() []model.User {
	var user model.User
	var users []model.User
	//init redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	//iterator buat scan keys di set users dalam redis
	iter := rdb.SScan(ctx, "users", 0, "", 0).Iterator()
	//jika ada,gw scan email ama type nya aja trus gw masukin ke users
	//jika tidak ada, gw bikin panic , itu ntar jadi berenti
	for iter.Next(ctx) {
		err := rdb.HMGet(ctx, iter.Val(), "email", "type").Scan(&user)
		if err != nil {
			panic(err)
		}
		users = append(users, user)
	}
	return users
}
func getAllUser() []model.User {
	db := connect()
	defer db.Close()

	query := "SELECT * FROM users"
	rows, _ := db.Query(query)

	var user model.User
	var users []model.User

	for rows.Next() {
		if err := rows.Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.Type); err != nil {
			log.Println(err)
			return nil
		} else {
			users = append(users, user)
		}
	}
	return users
}

func sendUnauthorizedResponse(w http.ResponseWriter) {
	var response model.Response
	response.Status = 401
	response.Message = "Unauthorized Access"
	w.Header().Set("Content=Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func Scheduler(w http.ResponseWriter, r *http.Request) {
	s := gocron.NewScheduler(time.UTC) //00.00 GMT

	//Set redis saat pertama kali scheduler dijalankan
	setAllUserRedis()

	//Do redis setiap satu jam
	s.Every(1).Hours().Do(setAllUserRedis)
	//send email setiap jam 6 (UTC +7)
	s.Every(1).Day().At("6.00").Do(func() {
		fmt.Println("Schedule Started")
		userList := getAllUserRedis()

		var userPremium []string
		var userBiasa []string

		for i := range userList {
			if userList[i].Type == 1 {
				userBiasa = append(userBiasa, userList[i].Email)
			}
			if userList[i].Type == 2 {
				userPremium = append(userPremium, userList[i].Email)
			}
		}
		go mailSending(userPremium, 1)
		go mailSending(userBiasa, 2)
		//time.Sleep(100 * time.Millisecond)
	})
	//Kirim email penawaran premium membership setiap jam 12 siang GMT+7
	s.Every(1).Day().At("12.00").Do(func() {
		userList := getAllUserRedis()

		var userBiasa []string

		for i := range userList {
			if userList[i].Type == 1 {
				userBiasa = append(userBiasa, userList[i].Email)
			}
		}
		sendPremiumOfferMail(userBiasa)
	})
	s.StartAsync()

	curtime := time.Now()
	var response model.Response
	response.Status = 200
	response.Message = "Schedule started at " + curtime.String()
	w.Header().Set("Content=Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func sendPremiumOfferMail(receivers []string) {
	m := gomail.NewMessage()
	text := "<h2>Premium Membership</h2>"
	text += "<p>There are many benefits that comes with a premium membership including but not limited to:</p><br>"
	text += "<ol><li>Access to exclusive articles</li><li>No advertisement before reading article</li></ol><br>"
	text += "<p>So, what are you waiting for? Get yourself a premium membership now!</p>"
	//Buat header untuk email
	for i := range receivers {
		m.SetHeader("From", "articler8375@gmail.com")
		m.SetHeader("To", receivers[i])
		m.SetHeader("Subject", "Premium Membership Offer")
		m.SetBody("text/html", text)
		d := gomail.NewDialer("smtp.gmail.com", 587, "articler8375@gmail.com", "1234")
		if err := d.DialAndSend(m); err != nil {
			log.Println("Failed to send premium offer email to", receivers[i])
		} else {
			log.Println("Successfully sent premium offer email to", receivers[i])
		}
	}
}
