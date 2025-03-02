package tests

import (
	"fmt"
	"github.com/xvertile/matrixsearch"
	"strconv"
	"testing"

	"github.com/bxcodec/faker/v3"
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
	Model        string       `faker:"oneof: ModelS, Model3, ModelX, ModelY, Roadster" text:"model"`
	Brand        string       `faker:"oneof: Tesla, BMW, Audi, Mercedes, Ford" text:"brand"`
	Year         int          `faker:"boundary_start=1990, boundary_end=2023" text:"year"`
	Color        string       `faker:"oneof: Red, Blue, Black, White, Silver" text:"color"`
	Engine       Engine       `faker:"-"`
	Manufacturer Manufacturer `faker:"-"`
	Price        float64      `faker:"boundary_start=20000, boundary_end=150000" text:"price"`
}

func carIndexer(c Car) []string {
	return []string{
		"model:" + c.Model,
		"brand:" + c.Brand,
		"year:" + strconv.Itoa(c.Year),
		"color:" + c.Color,
		"manufacturer:" + c.Manufacturer.Name,
		"country:" + c.Manufacturer.Country,
	}
}

func logCarKeys(c Car, t *testing.T) {
	keys := carIndexer(c)
	for _, k := range keys {
		t.Log(k)
	}
}

func TestCarDataStore(t *testing.T) {
	ds := matrixsearch.NewDataStore(func(c Car) string { return c.Model + "-" + c.Manufacturer.Name }, carIndexer)
	var c Car
	if err := faker.FakeData(&c); err != nil {
		t.Fatal(err)
	}
	c.Manufacturer.Name = faker.Username()
	c.Manufacturer.Country = faker.Word()
	c.Engine.Type = "V8"
	c.Engine.Horsepower = 400
	c.Engine.Capacity = 4.0
	ds.Insert(c)
	query := "model:" + c.Model + ":brand:" + c.Brand + ":year:" + strconv.Itoa(c.Year) + ":color:" + c.Color + ":manufacturer:" + c.Manufacturer.Name + ":country:" + c.Manufacturer.Country
	results := ds.Search(query)
	if len(results) == 0 {
		t.Error("Expected to find inserted car but got 0 results")
	} else {
		for _, r := range results {
			t.Log("Car found:")
			logCarKeys(r, t)
		}
	}
}

func TestCarUpdate(t *testing.T) {
	ds := matrixsearch.NewDataStore(func(c Car) string { return c.Model + "-" + c.Manufacturer.Name }, carIndexer)
	var c Car
	if err := faker.FakeData(&c); err != nil {
		t.Fatal(err)
	}
	c.Manufacturer.Name = faker.Username()
	c.Manufacturer.Country = faker.Word()
	c.Engine.Type = "V6"
	c.Engine.Horsepower = 350
	c.Engine.Capacity = 3.5
	ds.Insert(c)
	c.Color = "Magenta"
	ds.Update(c)
	query := "model:" + c.Model + ":brand:" + c.Brand + ":year:" + strconv.Itoa(c.Year) + ":color:" + c.Color + ":manufacturer:" + c.Manufacturer.Name + ":country:" + c.Manufacturer.Country
	results := ds.Search(query)
	if len(results) == 0 {
		t.Error("Expected updated car but got 0 results")
	} else {
		for _, r := range results {
			t.Log("Updated car found:")
			logCarKeys(r, t)
		}
	}
}

func TestCarDelete(t *testing.T) {
	ds := matrixsearch.NewDataStore(func(c Car) string { return c.Model + "-" + c.Manufacturer.Name }, carIndexer)
	var c Car
	if err := faker.FakeData(&c); err != nil {
		t.Fatal(err)
	}
	c.Manufacturer.Name = faker.Username()
	c.Manufacturer.Country = faker.Word()
	c.Engine.Type = "V8"
	c.Engine.Horsepower = 400
	c.Engine.Capacity = 4.0
	ds.Insert(c)
	ds.Delete(c)
	query := "model:" + c.Model + ":brand:" + c.Brand + ":year:" + strconv.Itoa(c.Year) + ":color:" + c.Color + ":manufacturer:" + c.Manufacturer.Name + ":country:" + c.Manufacturer.Country
	results := ds.Search(query)
	if len(results) != 0 {
		t.Error("Expected no results after deletion")
	} else {
		t.Log("Deletion confirmed, no results found")
	}
}

func TestCarSearchNonExistent(t *testing.T) {
	ds := matrixsearch.NewDataStore(func(c Car) string { return c.Model + "-" + c.Manufacturer.Name }, carIndexer)
	results := ds.Search("model:NonExistent:brand:Fake:year:0:color:None:manufacturer:None:country:None")
	if len(results) != 0 {
		t.Error("Expected no results for non-existent query")
	} else {
		t.Log("No results for non-existent query as expected")
	}
}

func TestCarMultipleInsert(t *testing.T) {
	ds := matrixsearch.NewDataStore(func(c Car) string { return c.Model + "-" + c.Manufacturer.Name }, carIndexer)
	for i := 0; i < 5; i++ {
		var c Car
		if err := faker.FakeData(&c); err != nil {
			t.Fatal(err)
		}
		c.Manufacturer.Name = faker.Username()
		c.Manufacturer.Country = faker.Word()
		c.Engine.Type = "V8"
		c.Engine.Horsepower = 400 + i
		c.Engine.Capacity = 4.0 + float64(i)
		ds.Insert(c)
	}
	query := "brand:" + "Tesla"
	results := ds.Search(query)
	t.Log("Multiple insert test, results:")
	for _, r := range results {
		logCarKeys(r, t)
	}
}

func TestCarComplexQuery(t *testing.T) {
	ds := matrixsearch.NewDataStore(func(c Car) string { return c.Model + "-" + c.Manufacturer.Name }, carIndexer)
	var c Car
	if err := faker.FakeData(&c); err != nil {
		t.Fatal(err)
	}
	c.Manufacturer.Name = "ComplexManu"
	c.Manufacturer.Country = "ComplexCountry"
	c.Engine.Type = "V12"
	c.Engine.Horsepower = 500
	c.Engine.Capacity = 5.5
	ds.Insert(c)
	query := "model:" + c.Model + ":brand:" + c.Brand + ":year:" + strconv.Itoa(c.Year) + ":color:" + c.Color + ":manufacturer:" + c.Manufacturer.Name + ":country:" + c.Manufacturer.Country
	results := ds.Search(query)
	t.Log("Complex query test, results:")
	for _, r := range results {
		logCarKeys(r, t)
	}
}

func BenchmarkCarSearchRandom(b *testing.B) {
	sizes := []int{10000, 100000, 1000000}
	for _, size := range sizes {
		b.Run(fmt.Sprintf("Size_%d", size), func(b *testing.B) {
			ds := matrixsearch.NewDataStore(func(c Car) string { return c.Model + "-" + c.Manufacturer.Name }, carIndexer)
			var known Car
			for i := 0; i < size; i++ {
				var c Car
				if err := faker.FakeData(&c); err != nil {
					b.Fatal(err)
				}
				c.Manufacturer.Name = faker.Username()
				c.Manufacturer.Country = faker.Word()
				c.Engine.Type = "V8"
				c.Engine.Horsepower = 400
				c.Engine.Capacity = 4.0
				ds.Insert(c)
				if i == size/2 {
					known = c
				}
			}
			query := "model:" + known.Model + ":brand:" + known.Brand + ":year:" + strconv.Itoa(known.Year) + ":color:" + known.Color + ":manufacturer:" + known.Manufacturer.Name + ":country:" + known.Manufacturer.Country
			b.Log("Car SearchRandom Benchmark - Size:", size, "Query:", query)
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
