package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	cache "github.com/patrickmn/go-cache"
)

var c = cache.New(30*time.Minute, 60*time.Minute)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/opendata/airports", getOpenDataAirports).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", router))
}

func getOpenDataAirports(w http.ResponseWriter, r *http.Request) {
	var airportURL = "https://data.delaware.gov/api/views/mh8v-eba6/rows.json?accessType=DOWNLOAD"
	var airportJsonKey = "airportJson"

	airportJsonData, found := c.Get(airportJsonKey)
	if found {
		json := airportJsonData.(string)
		w.Write([]byte(json))
	} else {
		resp, err := http.Get(airportURL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		body, err := ioutil.ReadAll(resp.Body)
		c.Set(airportJsonKey, string(body[:]), cache.DefaultExpiration)
		w.Write(body)
	}
}
