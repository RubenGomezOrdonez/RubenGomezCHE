package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadPizzas(t *testing.T) {
	if _, err := os.Stat(PIZZAFILE); os.IsNotExist(err) {
		log.Println("Pizza file not found")
		return
	}
	pp, err := LoadPizzas()
	assert.Nil(t, err)

	assert.Equal(t, 3, len(pp))
}

func TestLoadCustomers(t *testing.T) {
	if _, err := os.Stat(CUSTOMERFILE); os.IsNotExist(err) {
		log.Println("Customer file not found")
		return
	}
	cc, err := LoadCustomers()
	assert.Nil(t, err)

	assert.True(t, len(cc) >= 2)
}

func TestLoadOrders(t *testing.T) {
	if _, err := os.Stat(ORDERFILE); os.IsNotExist(err) {
		log.Println("Order file not found")
		return
	}
	oo, err := LoadOrders()
	assert.Nil(t, err)

	assert.True(t, len(oo) >= 3)
}

func TestLoadOrdersForCustomer(t *testing.T) {
	oo, err := LoadOrdersForCustomer("rlcomte@che.nl")
	assert.Nil(t, err)

	assert.True(t, len(oo) >= 1)
}

func TestSaveCustomer(t *testing.T) {
	cc, err := LoadCustomers()
	assert.Nil(t, err)
	n := len(cc)

	c := Customer{Email: "petedoe@acme.com", Name: "Pete Doe"}
	err = SaveCustomer(c)
	assert.Nil(t, err)

	cc, err = LoadCustomers()
	assert.Nil(t, err)
	assert.True(t, len(cc) >= n)

	err = DeleteCustomer("petedoe@acme.com")
	assert.Nil(t, err)

	cc, err = LoadCustomers()
	assert.Nil(t, err)
	assert.True(t, len(cc) == n)
}

func TestSaveOrder(t *testing.T) {
	cc, err := LoadOrders()
	assert.Nil(t, err)
	n := len(cc)

	o := Order{
		No:    (n + 1),
		Email: "petedoe@acme.com",
		Lines: []OrderLine{
			OrderLine{
				Pizza: "margarita",
				Qty:   1,
			},
		},
	}
	err = SaveOrder(o)
	assert.Nil(t, err)

	cc, err = LoadOrders()
	assert.Nil(t, err)
	assert.True(t, len(cc) >= n)

	err = DeleteOrder(n + 1)
	assert.Nil(t, err)

	cc, err = LoadOrders()
	assert.Nil(t, err)
	assert.True(t, len(cc) == n)
}

func TestIngredients(t *testing.T) {
	imap, err := LoadIngredients()
	assert.Nil(t, err)

	assert.Equal(t, 3, len(imap["bodem"]))
	assert.Equal(t, "Italiaans", imap["bodem"][0].Name)
	assert.Equal(t, float32(0.25), imap["topping"][0].Price)
}

func TestToptien(t *testing.T) {
	pp, err := LoadToptienPizzas()
	assert.Nil(t, err)
	if len(pp) != 10 {
		fmt.Println("Totaal aantal pizza is niet 10")
	}

	assert.Equal(t, 10, len(pp))
}

func TestToptiensorteren(t *testing.T) {
	var toptienpizzas []Pizza
	pizzalijst := []Pizza{{"marinara", 10, 7.25, "", nil}, {"quattro stagioni", 7, 11.26, "../images/img61l.jpg", nil}, {"calzone", 4, 10.25, "../images/Calzone-Pizza-Slices-e1459845640761.jpg", nil}, {"margherita", 2, 8.99, "../images/Margarita.png", nil}, {"siciliano", 1, 9.29, "../images/siciliano.png", nil}, {"hawaii", 2, 9.99, "../images/hawaii.png", nil}, {"dönner", 8, 10.25, "../images/Doner_Pizza-924.png", nil}, {"hawaii zonder ananas", 5, 8.25, "../images/depositphotos_16813655-stockafbeelding-pizza-ham-en-mozzarella.jpg", nil}, {"salami", 6, 8.55, " ../images/Salami_pizza-2887.png", nil}, {"carciofi", 3, 12.25, "../images/DSC_6645.jpg", nil}, {"testpizza", 12, 12.25, "../images/DSC_6645.jpg", nil}} //testpizza is de hoogste

	sort.SliceStable(pizzalijst, func(l, j int) bool { return pizzalijst[l].Aantal > pizzalijst[j].Aantal })

	toptienpizzas = pizzalijst[0:10]
	expectedaantal := 11 // zet hier de expected aantal
	var werkelijkeaantal int
	werkelijkeaantal = toptienpizzas[0].Aantal
	if werkelijkeaantal != expectedaantal {
		t.Errorf("Pizza is niet gesorteert van groot naar klein heeft: %v moet zijn: %v.", werkelijkeaantal, expectedaantal)
	} else {
		t.Log("Pizza is goed gesorteerd van groot naar laag!")
	}
}
