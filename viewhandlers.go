package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"sort"
	"strconv"
	"time"
)

func pizzaHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Viewing pizzas")
	pp, _ := LoadPizzas()
	t, _ := template.ParseFiles("./templates/pizzas.html")

	t.Execute(w, pp)
}

func custompizzaHandler(w http.ResponseWriter, r *http.Request) {
	ingredientsMap, _ := LoadIngredients()

	if r.Method == "POST" {
		log.Println("Ordering custom pizza")
		r.ParseForm()
		customer := Customer{Email: r.FormValue("email"), Name: r.FormValue("name"), Address: r.FormValue("address"), PhoneNumber: r.FormValue("phone")}
		// create the customer
		SaveCustomer(customer)

		date := time.Now()
		order := Order{No: getOrderNo(), Date: date, DateString: date.Format(time.RFC1123Z), Email: customer.Email}

		bodem := r.FormValue("bodem")
		saus := r.FormValue("saus")
		toppings := r.Form["toppings"]
		log.Printf("%#v\n", toppings)

		findingredientbodem, _ := FindIngredient(bodem, ingredientsMap["bodem"])
		findingredientsaus, _ := FindIngredient(saus, ingredientsMap["saus"])

		var ingredienten []Ingredient
		ingredienten = append(ingredienten, findingredientbodem)
		ingredienten = append(ingredienten, findingredientsaus)

		for i := 0; i < len(toppings); i++ {
			findingredienttoppings, _ := FindIngredient(toppings[i], ingredientsMap["topping"])
			if toppings[i] != "" {
				ingredienten = append(ingredienten, findingredienttoppings)
			}
		}

		customPizza := CustomPizza{Ingredienten: ingredienten}
		var orderlines []OrderLine
		orderlines = append(orderlines, OrderLine{Qty: 1, CustomPizza: customPizza})
		order.Lines = orderlines
		order.Total = calculateTotal(order)

		// create the order
		SaveOrder(order)
		http.Redirect(w, r, "../html/thankyou.html", http.StatusSeeOther)
	} else {
		log.Println("Viewing custom pizza")
		t, _ := template.ParseFiles("./templates/custompizza.html")
		t.Execute(w, ingredientsMap)
	}
	from := "pizzahouse.ede@gmail.com"
	pass := "CHE2019!"
	var to string
	var body string
	var email string
	var naam string
	var adres string
	var postcode string
	var telefoon string

	naam = r.FormValue("name")
	adres = r.FormValue("address")
	postcode = r.FormValue("postalcode")
	telefoon = r.FormValue("phoneNumber")
	email = r.FormValue("email")

	fmt.Println(email)
	body = "Bestelling is aangekomen. Deze wordt zo snel mogelijk bereid." + " " + " " +
		"Uw gegevens: " + " " +
		"Naam: " + " " + naam + ", " + "Adres: " + " " + adres + ", " + "Postcode " + " " + postcode + ", " + "Telefoon: " + " " + telefoon

	to = email

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: Je Pizza bestelling\n\n" +
		body

	sendmail := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if sendmail != nil {
		log.Printf("smtp error: %s", sendmail)
		return
	}
}

func customerHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Viewing customers")
	cc, _ := LoadCustomersWithOrders()
	t, _ := template.ParseFiles("./templates/customers.html")

	t.Execute(w, cc)
}

func orderHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Viewing orders")
	oo, err := LoadOrders()
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "../html/error.html", http.StatusInternalServerError)
	}
	t, _ := template.ParseFiles("./templates/orders.html")

	t.Execute(w, oo)
}

func orderdetailsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Viewing order details")
	tmp := r.URL.Query()["no"]
	orderno, err := strconv.Atoi(tmp[0])
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "../html/error.html", http.StatusBadRequest)
	}

	oo, err := LoadOrders()
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "../html/error.html", http.StatusInternalServerError)
	}

	for _, o := range oo {
		if o.No == orderno {
			t, _ := template.ParseFiles("./templates/orderdetails.html")

			t.Execute(w, o)
			return
		}
	}

	http.Redirect(w, r, "../html/error.html", http.StatusNotFound)
}

func redirecthtml(w http.ResponseWriter, r *http.Request) {

	http.Redirect(w, r, "/html/", 301)
}

func nieuwpizzaHandler(w http.ResponseWriter, r *http.Request) {
	pp, _ := LoadIngredients()
	t, _ := template.ParseFiles("./templates/nieuwpizza.html")
	t.Execute(w, pp)
}

func bewerkpaginaHandler(w http.ResponseWriter, r *http.Request) {
	pp, _ := LoadPizzas()
	t, _ := template.ParseFiles("./templates/bewerkpagina.html")
	t.Execute(w, pp)
}

func nieuwpizzaterugkoppelingHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./templates/nieuwpizzaterugkoppeling.html")
	ingredientsMap, _ := LoadIngredients()
	var pizzaslice []Pizza
	var resultaat Pizza
	var nieuwpizza []Pizza
	var toptienupdate Pizza
	var toptienslice []Pizza
	var toptienresultaat []Pizza
	readjson, err := ioutil.ReadFile("./data/pizzas.json")

	err = json.Unmarshal(readjson, &pizzaslice)
	if err != nil {
		panic(err)
	}
	for i := 0; i < len(pizzaslice); i++ {
		prijsfloat, _ := strconv.ParseFloat(r.FormValue("price"), 64)

		naam := r.FormValue("name")
		prijs := prijsfloat
		toppings := r.Form["toppings"]

		bodem := r.FormValue("bodem")
		saus := r.FormValue("saus")

		findingredientbodem, _ := FindIngredient(bodem, ingredientsMap["bodem"])
		findingredientsaus, _ := FindIngredient(saus, ingredientsMap["saus"])

		var ingredienten []Ingredient
		ingredienten = append(ingredienten, findingredientbodem)
		ingredienten = append(ingredienten, findingredientsaus)

		for i := 0; i < len(toppings); i++ {
			findingredienttoppings, _ := FindIngredient(toppings[i], ingredientsMap["topping"])
			if toppings[i] != "" {
				ingredienten = append(ingredienten, findingredienttoppings)
			}
		}
		fmt.Println("dit is ingridients:", ingredienten)
		k := 2
		if r.FormValue("image2") == "" {
			k = 1
		}
		switch k {
		case 1:
			plaatje := r.FormValue("image")
			plaatjenav := "../images/"
			plaatje = plaatjenav + plaatje
			resultaat = Pizza{Name: naam, Price: prijs, Image: plaatje, Ingredient: ingredienten}
			toptienupdate = Pizza{Name: naam, Price: prijs, Image: plaatje, Aantal: 0}
		case 2:
			plaatje := r.FormValue("image2")
			resultaat = Pizza{Name: naam, Price: prijs, Image: plaatje, Ingredient: ingredienten}
			toptienupdate = Pizza{Name: naam, Price: prijs, Image: plaatje, Aantal: 0}
		}

		nieuwpizza = append(pizzaslice, resultaat)
		raw, err := json.MarshalIndent(nieuwpizza, "", "\t")
		if err != nil {
			panic(err)
		}

		ioutil.WriteFile("./data/pizzas.json", raw, 0644)

	}

	Livepizzajson, err := ioutil.ReadFile("./data/Livepizza.json")
	err = json.Unmarshal(Livepizzajson, &toptienslice)
	if err != nil {
		panic(err)
	}
	for j := 0; j < len(toptienslice); j++ {

		toptienresultaat = append(toptienslice, toptienupdate)

		rawtoptien, err := json.MarshalIndent(toptienresultaat, "", "\t")
		if err != nil {
			panic(err)
		}
		ioutil.WriteFile("./data/Livepizza.json", rawtoptien, 0644)

	}
	deletecsvFile()
	CreateCSV()
	t.Execute(w, nieuwpizza)
}

func bewerkpizzaHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./templates/bewerkpizza.html")

	var pizzaslice []Pizza
	var liveslice []Pizza
	var toptienpizzas []Pizza

	readjson, err := ioutil.ReadFile("./data/pizzas.json")

	err = json.Unmarshal(readjson, &pizzaslice)
	if err != nil {
		panic(err)
	}
	for i := 0; i < len(pizzaslice); i++ {
		prijsfloat, _ := strconv.ParseFloat(r.FormValue("price"+pizzaslice[i].Name), 64)

		naam := r.FormValue("name" + pizzaslice[i].Name)

		prijs := prijsfloat
		aantal, _ := strconv.Atoi(r.FormValue("aantal" + pizzaslice[i].Name))
		k := 2
		if r.FormValue("image4"+pizzaslice[i].Name) == "" {
			k = 1
		}

		switch k {
		case 1:
			plaatje := r.FormValue("image3" + pizzaslice[i].Name)
			plaatjenav := "../images/"
			plaatje = plaatjenav + plaatje
			pizzaslice[i].Name = naam
			pizzaslice[i].Price = prijs
			pizzaslice[i].Image = plaatje
			pizzaslice[i].Aantal = aantal
		case 2:
			plaatje := r.FormValue("image4" + pizzaslice[i].Name)
			pizzaslice[i].Name = naam
			pizzaslice[i].Price = prijs
			pizzaslice[i].Image = plaatje
			pizzaslice[i].Aantal = aantal
		}
		var buf = new(bytes.Buffer)

		enc := json.NewEncoder(buf)
		enc.Encode(pizzaslice)
		f, err := os.Create("./data/pizzas.json")
		if nil != err {
			log.Fatalln(err)
		}
		defer f.Close()
		io.Copy(f, buf)

		readjson2, err := ioutil.ReadFile("./data/Livepizza.json")
		err = json.Unmarshal(readjson2, &liveslice)
		if err != nil {
			panic(err)
		}

		fmt.Println("dit is aantal:", aantal)
		sort.SliceStable(pizzaslice, func(l, j int) bool { return pizzaslice[l].Aantal > pizzaslice[j].Aantal })

		var buf2 = new(bytes.Buffer)

		enc2 := json.NewEncoder(buf2)
		enc2.Encode(pizzaslice)
		f2, err := os.Create("./data/Livepizza.json")
		if nil != err {
			log.Fatalln(err)
		}
		defer f2.Close()
		io.Copy(f2, buf2)

		toptienpizzas = pizzaslice[0:10]

		var buf3 = new(bytes.Buffer)

		enc3 := json.NewEncoder(buf3)
		enc3.Encode(toptienpizzas)
		f3, err := os.Create("./data/Toptienvorigeweek.json")
		if nil != err {
			log.Fatalln(err)
		}
		defer f3.Close()
		io.Copy(f3, buf3)

	}
	deletecsvFile()
	CreateCSV()
	t.Execute(w, pizzaslice)
}

func toptienHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Viewing pizzas")
	pp, _ := LoadToptienPizzas()
	t, _ := template.ParseFiles("./templates/toptien.html")

	t.Execute(w, pp)
}

func Pizzadelete(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./templates/pizzadelete.html")
	delete := r.FormValue("delete")
	DeletePizza(delete)
	DeletePizza2(delete)
	DeletePizza3(delete)
	t.Execute(w, delete)
}
func Loginhandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Loggin in")
	t, _ := template.ParseFiles("./templates/login.html")

	t.Execute(w, t)
}
