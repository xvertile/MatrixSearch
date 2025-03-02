package tests

import (
	"fmt"
	"github.com/xvertile/matrixsearch"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

type Geo struct {
	City    string
	State   string
	Country string
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
	ID        string
	IP        string
	Speed     int
	SpeedType string
	Mobile    bool
	Timezone  string
	Anonymous bool
	Satellite bool
	Hosting   bool
	Geo       Geo
	Privacy   Privacy
	ASN       ASN
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func getProxyID(p Proxy) string {
	return p.ID
}

func randomProxy(i int) Proxy {
	countries := []string{"us", "ca", "uk", "de", "fr"}
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
			City:    "city" + strconv.Itoa(rand.Intn(100)),
			State:   "state" + strconv.Itoa(rand.Intn(50)),
			Country: countries[rand.Intn(len(countries))],
			Lat:     rand.Float64() * 100,
			Lon:     rand.Float64() * 100,
			Postal:  "postal" + strconv.Itoa(rand.Intn(1000)),
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

func indexProxy(p Proxy) []string {
	return []string{
		"country:" + p.Geo.Country,
		"state:" + p.Geo.State,
		"speedtype:" + p.SpeedType,
		"mobile:" + fmt.Sprintf("%t", p.Mobile),
	}
}

func logProxyKeys(p Proxy, t *testing.T) {
	for _, k := range indexProxy(p) {
		t.Log(k)
	}
}

func TestProxyDataStore(t *testing.T) {
	ds := matrixsearch.NewDataStore(getProxyID, indexProxy)
	p := randomProxy(1)
	ds.Insert(p)
	query := "country:" + p.Geo.Country + ":state:" + p.Geo.State + ":speedtype:" + p.SpeedType + ":mobile:" + fmt.Sprintf("%t", p.Mobile)
	results := ds.Search(query)
	t.Log("Query:", query)
	if len(results) == 0 {
		t.Error("Expected to find inserted proxy but got 0 results")
	} else {
		for _, r := range results {
			t.Log("Proxy found:")
			logProxyKeys(r, t)
		}
	}
}

func TestProxyUpdate(t *testing.T) {
	ds := matrixsearch.NewDataStore(getProxyID, indexProxy)
	p := randomProxy(2)
	ds.Insert(p)
	p.SpeedType = "ultrafast"
	ds.Update(p)
	query := "country:" + p.Geo.Country + ":state:" + p.Geo.State + ":speedtype:" + p.SpeedType + ":mobile:" + fmt.Sprintf("%t", p.Mobile)
	results := ds.Search(query)
	t.Log("Query after update:", query)
	if len(results) == 0 {
		t.Error("Expected updated proxy but got 0 results")
	} else {
		for _, r := range results {
			t.Log("Updated proxy found:")
			logProxyKeys(r, t)
		}
	}
}

func TestProxyDelete(t *testing.T) {
	ds := matrixsearch.NewDataStore(getProxyID, indexProxy)
	p := randomProxy(3)
	ds.Insert(p)
	ds.Delete(p)
	query := "country:" + p.Geo.Country + ":state:" + p.Geo.State + ":speedtype:" + p.SpeedType + ":mobile:" + fmt.Sprintf("%t", p.Mobile)
	results := ds.Search(query)
	t.Log("Query after deletion:", query)
	if len(results) != 0 {
		t.Error("Expected no results after deletion")
	} else {
		t.Log("Deletion confirmed, no results found")
	}
}

func TestProxySearchNonExistent(t *testing.T) {
	ds := matrixsearch.NewDataStore(getProxyID, indexProxy)
	query := "country:nonexistent:state:nonexistent:speedtype:nonexistent:mobile:false"
	results := ds.Search(query)
	t.Log("Non-existent query:", query)
	if len(results) != 0 {
		t.Error("Expected no results for non-existent query")
	} else {
		t.Log("No results for non-existent query as expected")
	}
}

func BenchmarkProxySearchRandom(b *testing.B) {
	sizes := []int{10000, 100000, 1000000}
	for _, size := range sizes {
		b.Run(fmt.Sprintf("Size_%d", size), func(b *testing.B) {
			ds := matrixsearch.NewDataStore(getProxyID, indexProxy)
			var known Proxy
			for i := 0; i < size; i++ {
				p := randomProxy(i)
				ds.Insert(p)
				if i == size/2 {
					known = p
				}
			}
			query := "country:" + known.Geo.Country + ":state:" + known.Geo.State + ":speedtype:" + known.SpeedType + ":mobile:" + fmt.Sprintf("%t", known.Mobile)
			b.Log("Proxy SearchRandom Benchmark - Size:", size, "Query:", query)
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
