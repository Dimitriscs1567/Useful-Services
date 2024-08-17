package compare_consumption

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"slices"
	"strconv"
	"strings"
)

func getGasPrices() map[string]float64 {
	prices := make(map[string]float64)
	res, err := http.Get("https://tradingeconomics.com/country-list/gasoline-prices?continent=world")
	if err != nil {
		log.Println(err)
		return prices
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return prices
	}

	res2, err := http.Get("https://www.ecb.europa.eu/stats/policy_and_exchange_rates/euro_reference_exchange_rates/html/eurofxref-graph-usd.en.html")
	if err != nil {
		log.Println(err)
		return prices
	}
	body2, err := io.ReadAll(res2.Body)
	if err != nil {
		log.Println(err)
		return prices
	}
	sb2 := string(body2)
	slpitted1 := strings.Split(sb2, "rateLatestInverse='")[1]
	slpitted2 := strings.Split(slpitted1, "'")[0]
	toEuroConv, err := strconv.ParseFloat(slpitted2, 64)
	if err != nil {
		log.Println(err)
		return prices
	}

	sb := string(body)
	country_tables := strings.Split(sb, "<tr class='datatable-row")
	country_tables = country_tables[1:]
	for _, table := range country_tables {
		first_split := strings.Split(table, "<a href")[1]
		second_split := strings.Split(first_split, ">")[1]
		country := strings.TrimSpace(strings.Split(second_split, "<")[0])

		first_split = strings.Split(table, "<td data-heatmap-value=")[1]
		second_split = strings.Split(first_split, ">")[1]
		price := strings.Split(second_split, "</td")[0]

		prices[country], err = strconv.ParseFloat(price, 64)
		if err != nil {
			fmt.Println(err)
			prices[country] = 0
		}
		prices[country] = math.Round((prices[country]*toEuroConv)*100) / 100
	}

	return prices
}

func getElectricityPrices() map[string][]float64 {
	interest_titles := []string{"Energy price for AC recharging (€/kWh)", "Energy price for DC recharging (€/kWh)"}

	prices := make(map[string][]float64)
	res, err := http.Get("https://alternative-fuels-observatory.ec.europa.eu/consumer-portal/electric-vehicle-recharging-prices")

	if err != nil {
		log.Println(err)
		return prices
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return prices
	}

	sb := string(body)
	chart_titles := strings.Split(sb, "chart-title\">")
	chart_titles = chart_titles[1:]
	for _, data := range chart_titles {
		title := strings.Split(data, "</h2>")[0]
		if slices.Contains(interest_titles, title) {
			first_split := strings.Split(data, "type=\"application/json\">")[1]
			json_string := strings.Split(first_split, "</script>")[0]

			var json_data map[string]any
			json.Unmarshal([]byte(json_string), &json_data)

			json1, ok := json_data["data"].(map[string]any)
			if !ok {
				return prices
			}

			json2, ok := json1["xAxis"].(map[string]any)
			if !ok {
				return prices
			}

			json3 := json2["categories"]
			countries_list, ok := json3.([]any)
			if !ok {
				return prices
			}

			for _, country := range countries_list {
				country_name, ok := country.(string)
				if !ok {
					return prices
				}

				if prices[country_name] == nil {
					prices[country_name] = make([]float64, 2)
				}
			}

			json4, ok := json1["series"].([]any)
			if !ok {
				return prices
			}

			json5, ok := json4[0].(map[string]any)
			if !ok {
				return prices
			}

			data_list, ok := json5["data"].([]any)
			if !ok {
				return prices
			}

			for index, value_any := range data_list {
				value, ok := value_any.(float64)
				if !ok {
					return prices
				}

				if prices[countries_list[index].(string)][0] == 0 {
					prices[countries_list[index].(string)][0] = value
				} else {
					prices[countries_list[index].(string)][1] = value
				}

			}
		}
	}

	return prices
}

func getFuelPrices() []*FuelPrices {
	gas_prices := getGasPrices()
	electric_prices := getElectricityPrices()

	fuel_prices := []*FuelPrices{}

	for country, gas_price := range gas_prices {
		values, ok := electric_prices[country]
		if ok {
			fuel_prices = append(fuel_prices,
				NewFuelPrices(country, gas_price, values[0], values[1]))
		} else {
			fuel_prices = append(fuel_prices,
				NewFuelPrices(country, gas_price, 0, 0))
		}
	}

	for country, electric_price := range electric_prices {
		_, ok := gas_prices[country]
		if !ok {
			fuel_prices = append(fuel_prices,
				NewFuelPrices(country, 0, electric_price[0], electric_price[1]))
		}
	}

	return fuel_prices
}

func GetFuelComparison(data *CompareConsumptionBody) *CompareConsumptionResponse {
	var responseConsumption float64

	if data.ElectricityPrice == 0 || data.GasPrice == 0 {
		return &CompareConsumptionResponse{
			L100kmTokwh100km: 0,
		}
	}

	gasMoney := data.GasPrice
	responseConsumption = math.Round((gasMoney/data.ElectricityPrice)*1000) / 1000

	return &CompareConsumptionResponse{
		L100kmTokwh100km: responseConsumption,
	}
}
