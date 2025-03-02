package tests

import (
	"fmt"
	"github.com/xvertile/matrixsearch"
	"strconv"
	"testing"

	"github.com/bxcodec/faker/v3"
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
	Name      string  `faker:"oneof: Apple, Banana, Cherry, Durian, Elderberry" text:"name"`
	Color     string  `faker:"oneof: Red, Yellow, Green, Purple, Orange" text:"color"`
	Weight    float64 `faker:"boundary_start=0.1, boundary_end=2.5" text:"weight"`
	Taste     string  `faker:"oneof: Sweet, Sour, Bitter, Tangy" text:"taste"`
	Price     float64 `faker:"boundary_start=0.5, boundary_end=10" text:"price"`
	Origin    Origin
	Nutrition Nutrition
}

func fruitIndexer(f Fruit) []string {
	return []string{
		"name:" + f.Name,
		"color:" + f.Color,
		"country:" + f.Origin.Country,
		"harvestyear:" + strconv.Itoa(f.Origin.HarvestYear),
		"calories:" + strconv.Itoa(f.Nutrition.Calories),
	}
}

func logFruitKeys(f Fruit, t *testing.T) {
	keys := fruitIndexer(f)
	for _, k := range keys {
		t.Log(k)
	}
}

func TestFruitDataStore(t *testing.T) {
	ds := matrixsearch.NewDataStore(func(f Fruit) string { return f.Name + "-" + f.Origin.Country }, fruitIndexer)
	var f Fruit
	if err := faker.FakeData(&f); err != nil {
		t.Fatal(err)
	}
	ds.Insert(f)
	query := "name:" + f.Name + ":color:" + f.Color + ":country:" + f.Origin.Country + ":harvestyear:" + strconv.Itoa(f.Origin.HarvestYear) + ":calories:" + strconv.Itoa(f.Nutrition.Calories)
	results := ds.Search(query)
	if len(results) == 0 {
		t.Error("Expected to find inserted fruit but got 0 results")
	} else {
		for _, r := range results {
			t.Log("Fruit found with query:", query)
			logFruitKeys(r, t)
		}
	}
}

func TestFruitUpdate(t *testing.T) {
	ds := matrixsearch.NewDataStore(func(f Fruit) string { return f.Name + "-" + f.Origin.Country }, fruitIndexer)
	var f Fruit
	if err := faker.FakeData(&f); err != nil {
		t.Fatal(err)
	}
	ds.Insert(f)
	f.Color = "Cyan"
	ds.Update(f)
	query := "name:" + f.Name + ":color:" + f.Color + ":country:" + f.Origin.Country + ":harvestyear:" + strconv.Itoa(f.Origin.HarvestYear) + ":calories:" + strconv.Itoa(f.Nutrition.Calories)
	results := ds.Search(query)
	if len(results) == 0 {
		t.Error("Expected updated fruit but got 0 results")
	} else {
		for _, r := range results {
			t.Log("Updated fruit found:")
			logFruitKeys(r, t)
		}
	}
}

func TestFruitDelete(t *testing.T) {
	ds := matrixsearch.NewDataStore(func(f Fruit) string { return f.Name + "-" + f.Origin.Country }, fruitIndexer)
	var f Fruit
	if err := faker.FakeData(&f); err != nil {
		t.Fatal(err)
	}
	ds.Insert(f)
	ds.Delete(f)
	query := "name:" + f.Name + ":color:" + f.Color + ":country:" + f.Origin.Country + ":harvestyear:" + strconv.Itoa(f.Origin.HarvestYear) + ":calories:" + strconv.Itoa(f.Nutrition.Calories)
	results := ds.Search(query)
	if len(results) != 0 {
		t.Error("Expected no results after deletion")
	} else {
		t.Log("Deletion confirmed, no results found")
	}
}

func TestFruitSearchNonExistent(t *testing.T) {
	ds := matrixsearch.NewDataStore(func(f Fruit) string { return f.Name + "-" + f.Origin.Country }, fruitIndexer)
	results := ds.Search("name:NonExistent:color:Invisible:country:Nowhere:harvestyear:0:calories:0")
	if len(results) != 0 {
		t.Error("Expected no results for non-existent query")
	} else {
		t.Log("No results for non-existent query as expected")
	}
}

func TestFruitMultipleInsert(t *testing.T) {
	ds := matrixsearch.NewDataStore(func(f Fruit) string { return f.Name + "-" + f.Origin.Country }, fruitIndexer)
	for i := 0; i < 5; i++ {
		var f Fruit
		if err := faker.FakeData(&f); err != nil {
			t.Fatal(err)
		}
		ds.Insert(f)
	}
	query := "name:" + "Apple"
	results := ds.Search(query)
	t.Log("Multiple insert test, results:")
	for _, r := range results {
		logFruitKeys(r, t)
	}
}

func TestFruitComplexQuery(t *testing.T) {
	ds := matrixsearch.NewDataStore(func(f Fruit) string { return f.Name + "-" + f.Origin.Country }, fruitIndexer)
	var f Fruit
	if err := faker.FakeData(&f); err != nil {
		t.Fatal(err)
	}
	f.Origin.Country = "ComplexCountry"
	f.Origin.HarvestYear = 2022
	ds.Insert(f)
	query := "name:" + f.Name + ":color:" + f.Color + ":country:" + f.Origin.Country + ":harvestyear:" + strconv.Itoa(f.Origin.HarvestYear) + ":calories:" + strconv.Itoa(f.Nutrition.Calories)
	results := ds.Search(query)
	t.Log("Complex query test, results:")
	for _, r := range results {
		logFruitKeys(r, t)
	}
}

func BenchmarkFruitSearchRandom(b *testing.B) {
	sizes := []int{10000, 100000, 1000000}
	for _, size := range sizes {
		b.Run(fmt.Sprintf("Size_%d", size), func(b *testing.B) {
			ds := matrixsearch.NewDataStore(func(f Fruit) string { return f.Name + "-" + f.Origin.Country }, fruitIndexer)
			var known Fruit
			for i := 0; i < size; i++ {
				var f Fruit
				if err := faker.FakeData(&f); err != nil {
					b.Fatal(err)
				}
				ds.Insert(f)
				if i == size/2 {
					known = f
				}
			}
			query := "name:" + known.Name + ":color:" + known.Color + ":country:" + known.Origin.Country + ":harvestyear:" + strconv.Itoa(known.Origin.HarvestYear) + ":calories:" + strconv.Itoa(known.Nutrition.Calories)
			b.Log("Fruit SearchRandom Benchmark - Size:", size, "Query:", query)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				res, ok := ds.SearchRandom(query)
				if ok {
					_ = res
				}
			}
		})
	}
}
