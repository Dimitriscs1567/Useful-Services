package compare_consumption

import "github.com/gin-gonic/gin"

func RegisterRoutes(server *gin.Engine) {
	server.GET("compare_consumption/fill_prices", fillPrices)
	server.GET("compare_consumption/get_supported_countries", getSupportedCountries)
	server.POST("compare_consumption/compare_consumption", compareConsumption)
}
