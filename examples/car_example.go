// car_example.go
package main

import (
	"fmt"
	"github.com/xvertile/matrixsearch"
	"math/rand"
	"strconv"
	"time"
)

type Engine struct {
	Type       string  `text:"type"`
	Horsepower int     `text:"horsepower"`
	Capacity   float64 `text:"capacity"`
}

type Manufacturer struct {
	Name    string `text:"name"`
	Country string `text:"country"`
}

type Car struct {
	Model        string       `text:"model"`
	Brand        string       `text:"brand"`
	Year         int          `text:"year"`
	Color        string       `text:"color"`
	Engine       Engine       `text:"-"`
	Manufacturer Manufacturer `text:"manufacturer"`
	Price        float64      `text:"price"`
}

// Strict indexer: manually defined keys.
func strictCarIndexer(c Car) []string {
	return []string{
		"model:" + c.Model,
		"brand:" + c.Brand,
		"year:" + strconv.Itoa(c.Year),
		"color:" + c.Color,
		"manufacturer:" + c.Manufacturer.Name,
		"country:" + c.Manufacturer.Country,
	}
}

func strictCarExample() {
	fmt.Println("=== Car Strict Example ===")
	ds := matrixsearch.NewDataStore(func(c Car) string { return c.Model + "-" + c.Manufacturer.Name }, strictCarIndexer)
	car := Car{
		Model: "ModelS",
		Brand: "Tesla",
		Year:  2020,
		Color: "Red",
		Engine: Engine{
			Type:       "Electric",
			Horsepower: 500,
			Capacity:   0,
		},
		Manufacturer: Manufacturer{
			Name:    "Tesla Inc.",
			Country: "US",
		},
		Price: 80000,
	}
	ds.Insert(car)
	query := "model:" + car.Model + ":brand:" + car.Brand + ":year:" + strconv.Itoa(car.Year) +
		":color:" + car.Color + ":manufacturer:" + car.Manufacturer.Name + ":country:" + car.Manufacturer.Country
	fmt.Println("Strict Query:", query)
	results := ds.Search(query)
	if len(results) > 0 {
		for _, r := range results {
			fmt.Printf("Found Car: %+v\n", r)
		}
	} else {
		fmt.Println("No car found")
	}
}

func autoCarExample() {
	fmt.Println("=== Car Auto Example ===")
	ds := matrixsearch.NewDataStore(func(c Car) string { return c.Model + "-" + c.Manufacturer.Name }, matrixsearch.AutoIndexer[Car])
	car := Car{
		Model: "Model3",
		Brand: "Tesla",
		Year:  2021,
		Color: "Blue",
		Engine: Engine{
			Type:       "Electric",
			Horsepower: 450,
			Capacity:   0,
		},
		Manufacturer: Manufacturer{
			Name:    "Tesla Inc.",
			Country: "US",
		},
		Price: 50000,
	}
	ds.Insert(car)
	// Construct a query using a subset of auto-indexed keys.
	query := "model:" + car.Model +
		":brand:" + car.Brand +
		":year:" + strconv.Itoa(car.Year) +
		":color:" + car.Color +
		":country:" + car.Manufacturer.Country
	fmt.Println("Auto Query:", query)
	results := ds.Search(query)
	if len(results) > 0 {
		for _, r := range results {
			fmt.Printf("Found Car: %+v\n", r)
		}
	} else {
		fmt.Println("No car found")
	}
}

func bulkCarExample() {
	fmt.Println("=== Car Bulk Example ===")
	ds := matrixsearch.NewDataStore(func(c Car) string { return c.Model + "-" + c.Manufacturer.Name }, strictCarIndexer)
	models := []string{"ModelS", "Model3", "ModelX", "ModelY", "Roadster"}
	brands := []string{"Tesla", "BMW", "Audi", "Mercedes", "Ford"}
	colors := []string{"Red", "Blue", "Black", "White", "Silver"}
	manufacturers := []string{"Tesla Inc.", "BMW Group", "Audi AG", "Mercedes-Benz", "Ford Motor"}
	countries := []string{"US", "Germany", "UK", "Japan", "Canada"}
	for i := 0; i < 2000; i++ {
		car := Car{
			Model: models[rand.Intn(len(models))],
			Brand: brands[rand.Intn(len(brands))],
			Year:  1990 + rand.Intn(33),
			Color: colors[rand.Intn(len(colors))],
			Engine: Engine{
				Type:       "Type" + strconv.Itoa(rand.Intn(3)),
				Horsepower: 200 + rand.Intn(300),
				Capacity:   float64(rand.Intn(5)),
			},
			Manufacturer: Manufacturer{
				Name:    manufacturers[rand.Intn(len(manufacturers))],
				Country: countries[rand.Intn(len(countries))],
			},
			Price: float64(20000 + rand.Intn(130000)),
		}
		ds.Insert(car)
	}
	fmt.Println("Total Cars Inserted:", 2000)
	query := "brand:Tesla"
	fmt.Println("Bulk Query:", query)
	results := ds.Search(query)
	fmt.Println("Results found:", len(results))
	if len(results) > 0 {
		fmt.Printf("Sample Car: %+v\n", results[0])
	}
	r, ok := ds.SearchRandom(query)
	if ok {
		fmt.Printf("Random Car: %+v\n", r)
	} else {
		fmt.Println("SearchRandom found no car")
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	strictCarExample()
	autoCarExample()
	bulkCarExample()
}
