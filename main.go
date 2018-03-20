package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"

	. "github.com/ahmetb/go-linq"
	"github.com/gorilla/mux"
	"github.com/jinzhu/now"
)

var (
	events []EventModel
)

func main() {
	jsonFile, err := os.Open("./data/events.json")
	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &events)

	fmt.Println("Events count: ", len(events))

	for i, e := range events {
		events[i].DateString = strings.Join([]string{e.Date.Weekday().String(), e.Date.Month().String(), strconv.Itoa(e.Date.Day())}, " ")
	}

	router := mux.NewRouter()
	router.HandleFunc("/opendata/airports", getOpenDataAirports).Methods("GET")
	router.HandleFunc("/events/weekend/random", getRandomEventForWeekend).Methods("GET")
	router.HandleFunc("/events/today/random", getRandomEventForToday).Methods("GET")
	router.HandleFunc("/events/tomorrow/random", getRandomEventForTomorrow).Methods("GET")

	log.Fatal(http.ListenAndServe(getPort(), router))
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
	err = json.Unmarshal(body, &a)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var buffer bytes.Buffer
	for i := range a {
		buffer.WriteString(a[i].Name)
		if i != len(a)-1 {
			buffer.WriteString(", ")
		}
	}

	airports := Airports{AirportsNames: strings.Replace(buffer.String(), "&", "and", -1)}
	err = json.NewEncoder(w).Encode(airports)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getRandomEventForWeekend(w http.ResponseWriter, r *http.Request) {
	sunday := now.Sunday()
	fmt.Println("Sunday : ", sunday)

	today := now.BeginningOfDay()
	diff := sunday.Sub(today)
	diffDays := int(diff.Hours() / 24)
	if diffDays > 2 {
		diffDays = 2
	}

	weekendStart := sunday.AddDate(0, 0, -diffDays)
	fmt.Println("Weekend start : ", weekendStart)

	var weekendEvents []EventModel

	From(events).
		Where(
			func(f interface{}) bool {
				return (f.(EventModel).Date.After(weekendStart) && f.(EventModel).Date.Before(sunday)) || (f.(EventModel).Date.Equal(sunday)) || (f.(EventModel).Date.Equal(weekendStart))
			},
		).
		ToSlice(&weekendEvents)

	fmt.Println("Weekend events count: ", len(weekendEvents))
	err := json.NewEncoder(w).Encode(weekendEvents[rand.Intn(len(weekendEvents))])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getRandomEventForToday(w http.ResponseWriter, r *http.Request) {
	today := now.BeginningOfDay()
	fmt.Println("Today : ", today)
	var todayEvents []EventModel

	From(events).
		Where(
			func(f interface{}) bool {
				return f.(EventModel).Date.Equal(today)
			},
		).
		ToSlice(&todayEvents)

	fmt.Println("Today events count: ", len(todayEvents))
	err := json.NewEncoder(w).Encode(todayEvents[rand.Intn(len(todayEvents))])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getRandomEventForTomorrow(w http.ResponseWriter, r *http.Request) {
	today := now.BeginningOfDay()
	tomorrow := today.AddDate(0, 0, 1)

	fmt.Println("Tomorrow: ", tomorrow)
	var tomorrowEvents []EventModel

	From(events).
		Where(
			func(f interface{}) bool {
				return f.(EventModel).Date.Equal(tomorrow)
			},
		).
		ToSlice(&tomorrowEvents)

	fmt.Println("Tomorrow events count: ", len(tomorrowEvents))
	err := json.NewEncoder(w).Encode(tomorrowEvents[rand.Intn(len(tomorrowEvents))])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getPort() string {
	var port = os.Getenv("PORT")
	if port == "" {
		port = "8000"
		fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
	}

	return ":" + port
}
