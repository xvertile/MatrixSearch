package main

import (
	"fmt"
	"github.com/xvertile/matrixsearch"
	"math/rand"
	"strconv"
	"time"
)

type Geo struct {
	City    string `text:"city"`
	State   string `text:"state"`
	Country string `text:"country"`
	Lat     float64
	Lon     float64
	Postal  string
}

type Privacy struct {
	Vpn     bool
	Proxy   bool
	Tor     bool
	Relay   bool
	Hosting bool
	Service string
}

type ASN struct {
	Number int
	Asn    string
	Name   string
	Domain string
	Type   string
}

type Proxy struct {
	ID        string `text:"id"`
	IP        string `text:"ip"`
	Speed     int    `text:"speed"`
	SpeedType string `text:"speedtype"`
	Mobile    bool   `text:"mobile"`
	Timezone  string `text:"timezone"`
	Anonymous bool
	Satellite bool
	Hosting   bool
	Geo       Geo
	Privacy   Privacy
	ASN       ASN
}

func strictProxyIndexer(p Proxy) []string {
	return []string{
		"country:" + p.Geo.Country,
		"state:" + p.Geo.State,
		"speedtype:" + p.SpeedType,
		"mobile:" + strconv.FormatBool(p.Mobile),
	}
}

func randomProxy(i int) Proxy {
	countries := []string{"us", "ca", "uk", "de", "fr"}
	states := []string{"CA", "NY", "TX", "FL", "IL"}
	speedVal := rand.Intn(200)
	speedClass := "slow"
	if speedVal > 150 {
		speedClass = "fast"
	} else if speedVal > 75 {
		speedClass = "medium"
	}
	return Proxy{
		ID:        strconv.Itoa(i),
		IP:        "192.168.1." + strconv.Itoa(rand.Intn(255)),
		Speed:     speedVal,
		SpeedType: speedClass,
		Mobile:    rand.Intn(2) == 0,
		Timezone:  "tz" + strconv.Itoa(rand.Intn(10)),
		Anonymous: rand.Intn(2) == 0,
		Satellite: rand.Intn(2) == 0,
		Hosting:   rand.Intn(2) == 0,
		Geo: Geo{
			City:    "City" + strconv.Itoa(rand.Intn(100)),
			State:   states[rand.Intn(len(states))],
			Country: countries[rand.Intn(len(countries))],
			Lat:     rand.Float64() * 100,
			Lon:     rand.Float64() * 100,
			Postal:  "Postal" + strconv.Itoa(rand.Intn(1000)),
		},
		Privacy: Privacy{
			Vpn:     rand.Intn(2) == 0,
			Proxy:   rand.Intn(2) == 0,
			Tor:     rand.Intn(2) == 0,
			Relay:   rand.Intn(2) == 0,
			Hosting: rand.Intn(2) == 0,
			Service: "service" + strconv.Itoa(rand.Intn(5)),
		},
		ASN: ASN{
			Number: rand.Intn(10000),
			Asn:    "asn" + strconv.Itoa(rand.Intn(100)),
			Name:   "name" + strconv.Itoa(rand.Intn(100)),
			Domain: "domain" + strconv.Itoa(rand.Intn(100)) + ".com",
			Type:   "type" + strconv.Itoa(rand.Intn(3)),
		},
	}
}

func strictProxyExample() {
	fmt.Println("=== Proxy Strict Example ===")
	ds := matrixsearch.NewDataStore(func(p Proxy) string { return p.ID }, strictProxyIndexer)
	proxy := randomProxy(1012)
	ds.Insert(proxy)
	query := "country:" + proxy.Geo.Country + ":state:" + proxy.Geo.State + ":speedtype:" + proxy.SpeedType + ":mobile:" + strconv.FormatBool(proxy.Mobile)
	fmt.Println("Strict Query:", query)
	results := ds.Search(query)
	if len(results) > 0 {
		for _, r := range results {
			fmt.Printf("Found Proxy: %+v\n", r)
		}
	} else {
		fmt.Println("No proxy found")
	}
	filename := "proxies.svg"
	err := ds.Dump(filename)
	if err != nil {
		fmt.Println("Error dumping DataStore:", err)
	} else {
		fmt.Println("Dump successfully written to", filename)
	}
}

func autoProxyExample() {
	fmt.Println("=== Proxy Auto Example ===")
	ds := matrixsearch.NewDataStore(func(p Proxy) string { return p.ID }, matrixsearch.AutoIndexer[Proxy])
	proxy := randomProxy(2020)
	ds.Insert(proxy)
	query := "speedtype:" + proxy.SpeedType +
		":mobile:" + strconv.FormatBool(proxy.Mobile) +
		":country:" + proxy.Geo.Country
	fmt.Println("Auto Query:", query)
	results := ds.Search(query)
	if len(results) > 0 {
		for _, r := range results {
			fmt.Printf("Found Proxy: %+v\n", r)
		}
	} else {
		fmt.Println("No proxy found")
	}
}

func bulkProxyExample() {
	fmt.Println("=== Proxy Bulk Example ===")
	ds := matrixsearch.NewDataStore(func(p Proxy) string { return p.IP }, strictProxyIndexer)
	for i := 0; i < 100; i++ {
		ds.Insert(randomProxy(i))
	}
	fmt.Println("Total Proxies Inserted:", 100)
	query := "speedtype:fast"
	fmt.Println("Bulk Query:", query)
	results := ds.Search(query)
	fmt.Println("Results found:", len(results))
	if len(results) > 0 {
		fmt.Printf("Sample Proxy: %+v\n", results[0])
	}
	r, ok := ds.SearchRandom(query)
	if ok {
		fmt.Printf("Random Proxy: %+v\n", r)
	} else {
		fmt.Println("SearchRandom found no proxy")
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	strictProxyExample()
	autoProxyExample()
	bulkProxyExample()
}
