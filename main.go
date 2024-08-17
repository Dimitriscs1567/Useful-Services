package main

import (
	"github.com/gin-gonic/gin"
	"tikasdimitrios.com/usefull_services/compare_consumption"
	"tikasdimitrios.com/usefull_services/database"
)

func main() {
	database.InitDb()

	server := gin.Default()
	compare_consumption.RegisterRoutes(server)

	server.Run()
}
