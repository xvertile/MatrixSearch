// fruit_example.go
package main

import (
	"fmt"
	"github.com/xvertile/matrixsearch"
	"math/rand"
	"strconv"
	"time"
)

type Origin struct {
	Country     string `text:"country"`
	Region      string `text:"region"`
	HarvestYear int    `text:"harvestyear"`
}

type Nutrition struct {
	Calories int     `text:"calories"`
	Sugar    float64 `text:"sugar"`
	Fiber    float64 `text:"fiber"`
}

type Fruit struct {
	Name      string    `text:"name"`
	Color     string    `text:"color"`
	Weight    float64   `text:"weight"`
	Taste     string    `text:"taste"`
	Price     float64   `text:"price"`
	Origin    Origin    // AutoIndexer recurses into Origin.
	Nutrition Nutrition // AutoIndexer recurses into Nutrition.
}

// Strict indexer: manually defined keys.
func strictFruitIndexer(f Fruit) []string {
	return []string{
		"name:" + f.Name,
		"color:" + f.Color,
		"harvestyear:" + strconv.Itoa(f.Origin.HarvestYear),
		"calories:" + strconv.Itoa(f.Nutrition.Calories),
	}
}

func strictFruitExample() {
	fmt.Println("=== Fruit Strict Example ===")
	ds := matrixsearch.NewDataStore(func(f Fruit) string { return f.Name + "-" + f.Origin.Country }, strictFruitIndexer)
	fruit := Fruit{
		Name:   "Apple",
		Color:  "Red",
		Weight: 0.2,
		Taste:  "Sweet",
		Price:  1.5,
		Origin: Origin{
			Country:     "US",
			Region:      "Washington",
			HarvestYear: 2022,
		},
		Nutrition: Nutrition{
			Calories: 95,
			Sugar:    19,
			Fiber:    4,
		},
	}
	ds.Insert(fruit)
	query := "name:" + fruit.Name +
		":color:" + fruit.Color +
		":harvestyear:" + strconv.Itoa(fruit.Origin.HarvestYear) +
		":calories:" + strconv.Itoa(fruit.Nutrition.Calories)
	fmt.Println("Strict Query:", query)
	results := ds.Search(query)
	if len(results) > 0 {
		for _, r := range results {
			fmt.Printf("Found Fruit: %+v\n", r)
		}
	} else {
		fmt.Println("No fruit found")
	}
}

func autoFruitExample() {
	fmt.Println("=== Fruit Auto Example ===")
	ds := matrixsearch.NewDataStore(func(f Fruit) string { return f.Name + "-" + f.Origin.Country }, matrixsearch.AutoIndexer[Fruit])
	fruit := Fruit{
		Name:   "Banana",
		Color:  "Yellow",
		Weight: 0.15,
		Taste:  "Sweet",
		Price:  0.5,
		Origin: Origin{
			Country:     "Ecuador",
			Region:      "Guayas",
			HarvestYear: 2021,
		},
		Nutrition: Nutrition{
			Calories: 105,
			Sugar:    14,
			Fiber:    3,
		},
	}
	ds.Insert(fruit)
	query := "name:" + fruit.Name +
		":color:" + fruit.Color +
		":country:" + fruit.Origin.Country +
		":harvestyear:" + strconv.Itoa(fruit.Origin.HarvestYear)
	fmt.Println("Auto Query:", query)
	results := ds.Search(query)
	if len(results) > 0 {
		for _, r := range results {
			fmt.Printf("Found Fruit: %+v\n", r)
		}
	} else {
		fmt.Println("No fruit found")
	}
}

func bulkFruitExample() {
	fmt.Println("=== Fruit Bulk Example ===")
	ds := matrixsearch.NewDataStore(func(f Fruit) string { return f.Name + "-" + f.Origin.Country }, strictFruitIndexer)
	names := []string{"Apple", "Banana", "Cherry", "Durian", "Elderberry"}
	colors := []string{"Red", "Yellow", "Green", "Purple", "Orange"}
	countries := []string{"US", "Ecuador", "France", "Spain", "Italy"}
	regions := []string{"Region1", "Region2", "Region3"}
	for i := 0; i < 2000; i++ {
		fruit := Fruit{
			Name:   names[rand.Intn(len(names))],
			Color:  colors[rand.Intn(len(colors))],
			Weight: 0.1 + rand.Float64()*2.4,
			Taste:  "Sweet",
			Price:  0.5 + rand.Float64()*9.5,
			Origin: Origin{
				Country:     countries[rand.Intn(len(countries))],
				Region:      regions[rand.Intn(len(regions))],
				HarvestYear: 2015 + rand.Intn(10),
			},
			Nutrition: Nutrition{
				Calories: 50 + rand.Intn(150),
				Sugar:    rand.Float64() * 30,
				Fiber:    rand.Float64() * 10,
			},
		}
		ds.Insert(fruit)
	}
	fmt.Println("Total Fruits Inserted:", 2000)
	query := "name:Apple"
	fmt.Println("Bulk Query:", query)
	results := ds.Search(query)
	fmt.Println("Results found:", len(results))
	if len(results) > 0 {
		fmt.Printf("Sample Fruit: %+v\n", results[0])
	}
	r, ok := ds.SearchRandom(query)
	if ok {
		fmt.Printf("Random Fruit: %+v\n", r)
	} else {
		fmt.Println("SearchRandom found no fruit")
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	strictFruitExample()
	autoFruitExample()
	bulkFruitExample()
}
