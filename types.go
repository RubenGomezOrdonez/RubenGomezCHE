package main

import "time"

type Pizza struct {
	Name       string
	Aantal     int
	Price      float64
	Image      string
	Ingredient []Ingredient
}

type Customer struct {
	Email       string `json:"email"`
	Name        string `json:"name"`
	Address     string `json:"address"`
	Postalcode  string `json:"postalcode"`
	PhoneNumber string `json:"phoneNumber"`
}

type Order struct {
	No         int         `json:"no"`
	Date       time.Time   `json:"date"`
	DateString string      `json:"datestring"`
	Email      string      `json:"email"`
	Lines      []OrderLine `json:"lines"`
	Total      float64     `json:"total"`
}

type OrderLine struct {
	Pizza       string `json:"pizza"`
	Qty         int    `json:"qty"`
	CustomPizza CustomPizza
}

type CustomerWithOrders struct {
	Customer
	Orders []Order
}

type CustomPizza struct {
	Ingredienten []Ingredient
}

type Ingredient struct {
	Type  string  `json:"type"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}
