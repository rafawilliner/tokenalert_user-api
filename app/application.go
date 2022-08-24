package app

import (
	"tokenalert_user-api/datasources/mysql/users_db"

	"github.com/gin-gonic/gin"
)

var (
	router = gin.Default()
)

func StartApplication() {
	mapUrls()
	users_db.InitDataBase()
	router.Run(":8080")

}

