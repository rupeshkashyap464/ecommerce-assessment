package main

import (
	"encoding/json"
	"log"
	"os"
)

type DataStore struct {
	Users  []User  `json:"users"`
	Items  []Item  `json:"items"`
	Carts  []Cart  `json:"carts"`
	Orders []Order `json:"orders"`
}

var DB DataStore

// simple ID counters
var nextUserID uint = 1
var nextItemID uint = 1
var nextCartID uint = 1
var nextCartItemID uint = 1
var nextOrderID uint = 1
var nextOrderItemID uint = 1

// InitDB loads db.json if present and recomputes next IDs.
func InitDB(_ string) (*DataStore, error) {
	data, err := os.ReadFile("db.json")
	if err == nil {
		if err := json.Unmarshal(data, &DB); err != nil {
			log.Println("could not parse db.json:", err)
		}
	}

	// recompute next IDs
	for _, u := range DB.Users {
		if u.ID >= nextUserID {
			nextUserID = u.ID + 1
		}
	}
	for _, it := range DB.Items {
		if it.ID >= nextItemID {
			nextItemID = it.ID + 1
		}
	}
	for _, c := range DB.Carts {
		if c.ID >= nextCartID {
			nextCartID = c.ID + 1
		}
		for _, ci := range c.Items {
			if ci.ID >= nextCartItemID {
				nextCartItemID = ci.ID + 1
			}
		}
	}
	for _, o := range DB.Orders {
		if o.ID >= nextOrderID {
			nextOrderID = o.ID + 1
		}
		for _, oi := range o.Items {
			if oi.ID >= nextOrderItemID {
				nextOrderItemID = oi.ID + 1
			}
		}
	}

	return &DB, nil
}

// SaveDB writes DB to db.json
func SaveDB() {
	data, err := json.MarshalIndent(DB, "", "  ")
	if err != nil {
		log.Println("failed to marshal db:", err)
		return
	}
	if err := os.WriteFile("db.json", data, 0644); err != nil {
		log.Println("failed to write db.json:", err)
	}
}
