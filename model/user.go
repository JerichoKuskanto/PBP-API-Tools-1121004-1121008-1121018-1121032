package model

type User struct {
	ID       int    `json:"id" redis:"id"`
	Username string `json:"username" redis:"username"`
	Email    string `json:"email" redis:"email"`
	Password string `json:"password" redis:"password"`
	Type     int    `json:"type" redis:"type"`
}

type UsersResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    []User `json:"data"`
}
