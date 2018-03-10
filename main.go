package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/opendata/airports", getOpenDataAirports).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", router))
}

// Airport model.
type Airport struct {
	Name string
}

func getOpenDataAirports(w http.ResponseWriter, r *http.Request) {
	var airportURL = "https://data.delaware.gov/resource/mh8v-eba6.json?$select=name"
	resp, err := http.Get(airportURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	a := make([]Airport, 0)
	json.Unmarshal(body, &a)
	var buffer bytes.Buffer
	for i := range a {
		buffer.WriteString(a[i].Name)
		if i != len(a)-1 {
			buffer.WriteString(", ")
		}
	}

	w.Write(buffer.Bytes())
}

func getPort() string {
	var port = os.Getenv("ALEXAPORT")
	if port == "" {
		port = "8000"
		fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
	}

	return ":" + port
}
