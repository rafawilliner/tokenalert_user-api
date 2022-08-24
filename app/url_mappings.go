package app

import (
	"tokenalert_user-api/controllers/ping"
	"tokenalert_user-api/controllers/users"
)


func mapUrls() {
	router.GET("/ping", ping.Ping)

	router.GET("/users/:user_id", users.Get)
	router.POST("/users", users.Create)
	router.POST("/users/login", users.Login)
}
