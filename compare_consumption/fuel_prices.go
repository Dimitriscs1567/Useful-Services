package compare_consumption

type FuelPrices struct {
	Country string
	Gas     float64
	AC      float64
	DC      float64
}

func NewFuelPrices(country string, gas float64, ac float64, dc float64) *FuelPrices {
	return &FuelPrices{
		Country: country,
		Gas:     gas,
		AC:      ac,
		DC:      dc,
	}
}

func (fuel_prices *FuelPrices) toJson() map[string]interface{} {
	return map[string]interface{}{"country": fuel_prices.Country, "gas": fuel_prices.Gas, "ac": fuel_prices.AC, "dc": fuel_prices.DC}
}
