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
