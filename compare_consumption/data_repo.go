package compare_consumption

import (
	"tikasdimitrios.com/usefull_services/database"
)

func updateFuelPrices(fuel_prices *FuelPrices) error {
	err := database.UpdateOrAddItem(database.FUELPRICES, fuel_prices.toJson(), map[string]interface{}{"country": fuel_prices.Country})

	return err
}

func getCountries() ([]string, error) {
	res, err := database.GetMappedStringItems(database.FUELPRICES, "country")
	if err != nil {
		return nil, err
	}

	return res, nil
}
