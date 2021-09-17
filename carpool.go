package main

import (
	"encoding/json"
	"os"
	"strconv"
)

type Carpool struct {
	Periods []Period
}

type Period struct {
	Literprice  float64              // price per liter gas in euro
	Consumption float64              // liters per 100km
	Distance    float64              // one way distance in km
	Passengers  map[string]Passenger // people that took a ride in this period
}

type Passenger struct {
	Drivelog string
	Bill     float64
}

func main() {
	carpool := Carpool{
		Periods: []Period{
			{
				Literprice:  1.319,
				Consumption: 8,
				Distance:    50,
				Passengers: map[string]Passenger{
					"Tick": {
						Drivelog: "5534",
					},
					"Trick": {
						Drivelog: "55",
					},
				},
			},
			{
				Literprice:  1.339,
				Consumption: 8,
				Distance:    50,
				Passengers: map[string]Passenger{
					"Trick": {
						Drivelog: "3322",
					},
					"Track": {
						Drivelog: "2",
					},
				},
			},
		},
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "\t")

	for _, per := range carpool.Periods {
		per.calculate()
	}

	enc.Encode(carpool)
}

func (per *Period) calculate() {
	for name, pas := range per.Passengers {
		pas.Bill = 0
		for _, b := range []byte(pas.Drivelog) {
			i, err := strconv.Atoi(string(b))
			if err == nil {
				pas.Bill += ((per.Consumption / 100.) * per.Distance * per.Literprice) / float64(i)
			}
		}
		per.Passengers[name] = pas
	}
}
