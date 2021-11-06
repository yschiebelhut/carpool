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
	Fixed       float64              // fixed rate per km in euro
	Distance    float64              // one way distance in km
	Notes       string               // if you have anything to say, do it here
	Passengers  map[string]Passenger // people that took a ride in this period
	Calculated  bool                 // true if period already has been calculated
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

// calculate how much each passenger has to pay
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

// calculate all periods that still have to
func (carpool *Carpool) calculateAll() {
	for i, per := range carpool.Periods {
		if !per.Calculated {
			per.calculate()
			per.Calculated = true
			carpool.Periods[i] = per
		}
	}
}

// calculate all periods even if they already have been calculated
func (carpool *Carpool) recalculateAll() {
	for _, per := range carpool.Periods {
		per.calculate()
	}
}

// save the current working data to disk
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

// load saved data from disk
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

// calculate the total bill of a passenger
// This method in its current form is likely to get deprecated.
// The database is designed to be longterm but passengers will meanwhile already pay parts of their bill.
func (carpool *Carpool) totalPerPassenger() {
	carpool.calculateAll()

	sum := make(map[string]float64)
	for _, per := range carpool.Periods {
		for name, pas := range per.Passengers {
			sum[name] += pas.Bill
		}
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "\t")
	enc.Encode(sum)
}
