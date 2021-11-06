package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

var (
	dbFile string = "carpool.json"
)

type Carpool struct {
	Periods []Period
}

type Period struct {
	Literprice  float64              // price per liter gas in euro
	Consumption float64              // liters per 100km
	Distance    float64              // one way distance in km
	Passengers  map[string]Passenger // people that took a ride in this period
	Fixed       float64              // fixed rate per km in euro
}

type Passenger struct {
	Drivelog string
	Bill     float64
}

func main() {
	carpool, err := load()
	if err != nil {
		panic(err)
	}

	carpool.calculateAll()

	carpool.save()
}

func (per *Period) calculate() {
	for name, pas := range per.Passengers {
		pas.Bill = 0
		for _, b := range []byte(pas.Drivelog) {
			i, err := strconv.Atoi(string(b))
			if err == nil {
				pas.Bill += (((per.Consumption / 100.) * per.Distance * per.Literprice) + (per.Distance * per.Fixed)) / float64(i)
			}
		}
		per.Passengers[name] = pas
	}
}

func (carpool *Carpool) calculateAll() {
	for _, per := range carpool.Periods {
		per.calculate()
	}
}

func (carpool *Carpool) save() {
	file, err := os.Create(dbFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer file.Close()
	enc := json.NewEncoder(file)
	enc.SetIndent("", "\t")
	enc.Encode(carpool)

}

func load() (*Carpool, error) {
	file, err := os.Open(dbFile)
	if err != nil {
		fmt.Println(err)
		return &Carpool{}, err
	}

	defer file.Close()
	dec := json.NewDecoder(file)
	carpool := Carpool{}
	dec.Decode(&carpool)

	return &carpool, nil
}
