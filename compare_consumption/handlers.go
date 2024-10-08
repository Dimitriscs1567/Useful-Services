package compare_consumption

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CompareConsumptionBody struct {
	GasPrice         float64 `json:"gas_price"`
	ElectricityPrice float64 `json:"electricity_price"`
}

type CompareConsumptionResponse struct {
	L100kmTokwh100km float64 `json:"lkm_to_kwhkm"`
}

type GetCountryPricesBody struct {
	Country string `json:"country"`
}

func fillPrices(ctx *gin.Context) {
	allFuelPrices := getFuelPrices()

	for _, fuelPrices := range allFuelPrices {
		err := updateFuelPrices(fuelPrices)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func compareConsumption(ctx *gin.Context) {
	var body CompareConsumptionBody

	if err := ctx.BindJSON(&body); err != nil {
		return
	}

	ctx.JSON(http.StatusOK, GetFuelComparison(&body))
}

func getSupportedCountries(ctx *gin.Context) {
	res, err := getCountries()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func getCountryPrices(ctx *gin.Context) {
	var body GetCountryPricesBody

	if err := ctx.BindJSON(&body); err != nil {
		return
	}

	res, err := getFuelPricesForCountry(body.Country)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, res)
}
