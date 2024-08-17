package compare_consumption

import "github.com/gin-gonic/gin"

func RegisterRoutes(server *gin.Engine) {
	server.GET("/fill_prices", fillPrices)
	server.POST("/compare_consumption", compareConsumption)
}
