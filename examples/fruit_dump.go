package main

import (
	"fmt"
	"github.com/xvertile/matrixsearch"
	"math/rand"
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
	Name      string  `text:"name"`
	Color     string  `text:"color"`
	Weight    float64 `text:"weight"`
	Taste     string  `text:"taste"`
	Price     float64 `text:"price"`
	Origin    Origin
	Nutrition Nutrition
}

func strictFruitIndexer(f Fruit) []string {
	return []string{
		"name:" + f.Name,
		"color:" + f.Color,
	}
}

func randomFruit(i int) Fruit {
	names := []string{"Apple", "Banana", "Cherry", "Durian", "Elderberry"}
	colors := []string{"Red", "Yellow", "Green", "Purple", "Orange"}
	countries := []string{"US", "Ecuador", "France", "Spain", "Italy"}
	regions := []string{"Region1", "Region2", "Region3"}
	return Fruit{
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
}

func bulkFruitDumpExample() {
	ds := matrixsearch.NewDataStore(func(f Fruit) string { return f.Name + "-" + f.Origin.Country }, strictFruitIndexer)
	for i := 0; i < 10; i++ {
		ds.Insert(randomFruit(i))
	}
	query := "name:Apple"
	results := ds.Search(query)
	fmt.Println("Total Fruits Inserted:", 10)
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
	filename := "dumps/fruits.svg"
	err := ds.Dump(filename)
	if err != nil {
		fmt.Println("Error dumping DataStore:", err)
	} else {
		fmt.Println("Dump successfully written to", filename)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	bulkFruitDumpExample()
}
