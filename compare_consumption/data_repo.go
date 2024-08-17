package compare_consumption

import (
	"tikasdimitrios.com/usefull_services/database"
)

func updateFuelPrices(fuel_prices *FuelPrices) error {
	err := database.GetAndUpdateOrAddItem(database.FUELPRICES, fuel_prices.toJson(), map[string]interface{}{"country": fuel_prices.Country})

	return err
}
