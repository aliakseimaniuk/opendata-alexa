package main

import "time"

type EventModel struct {
	Name         string
	Description  string
	Link         string
	Date         time.Time
	Location     string
	LocationLink string
	Tags         []string
}
