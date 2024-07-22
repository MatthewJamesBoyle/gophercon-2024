package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Option struct {
	Name             string `json:"name"`
	PricePerNightGBP int    `json:"price_per_night_gbp"`
	StarRating       int    `json:"rating"`
}

type Response struct {
	HotelOptions []Option `json:"hotel_options"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	options := []Option{
		{Name: "Hotel A", PricePerNightGBP: 100, StarRating: 5},
		{Name: "Hotel B", PricePerNightGBP: 80, StarRating: 4},
		{Name: "Hotel C", PricePerNightGBP: 60, StarRating: 3},
	}

	response := Response{HotelOptions: options}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/hotels", handler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
