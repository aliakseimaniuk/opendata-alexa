package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/opendata/airports", getOpenDataAirports).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", router))
}

func getOpenDataAirports(w http.ResponseWriter, r *http.Request) {
	var airportUrl = "https://data.delaware.gov/api/views/mh8v-eba6/rows.json?accessType=DOWNLOAD"
	resp, err := http.Get(airportUrl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	w.Write(body)
}
