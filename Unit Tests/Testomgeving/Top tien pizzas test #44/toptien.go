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
	"os"
	"sort"
	"strconv"
	"time"
)

type Struct1 struct {
	Naam   string
	Aantal int
}

func main() {

	http.HandleFunc("/index", pizza)
	http.HandleFunc("/index2", toptienHandler)

	http.Handle("/", logRequest(http.FileServer(http.Dir("./"))))
	log.Fatal(http.ListenAndServe(":5000", nil))

}
func pizza(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./index.html")
	pizzas, _ := LoadPizzas()
	t.Execute(w, pizzas)

}
func toptienHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./index2.html")
	//var ticker *time.Ticker = nil

	file1 := "./Livepizza.json"
	file2 := "./Toptienvorigeweek.json"

	readjson, err := ioutil.ReadFile(file1)

	var pizzaslice []Struct1
	var toptienpizzas []Struct1
	var resultaatpizza []Struct1

	err = json.Unmarshal(readjson, &pizzaslice)
	if err != nil {
		panic(err)
	}

	var status int
	for i := 0; i < len(pizzaslice); i++ {

		naamvanpizza := r.FormValue("naam" + pizzaslice[i].Naam)

		for _, value := range pizzaslice {
			if value.Naam == naamvanpizza {
				status = value.Aantal

				fmt.Println("De huidige aantal van ", naamvanpizza, " is:", status)
			}
		}

		aantalint, _ := strconv.Atoi(r.FormValue("aantal" + pizzaslice[i].Naam))
		fmt.Println("De hoeveelheid ", naamvanpizza, " Pizza's die erbij gekomen zijn is:", r.FormValue("aantal"+pizzaslice[i].Naam))

		Toptienvorigeweekpizza := aantalint + status

		pizzaslice[i].Aantal = Toptienvorigeweekpizza

		fmt.Println("Toptienvorigeweek aantal pizza van", naamvanpizza, " is:", Toptienvorigeweekpizza)

		sort.SliceStable(pizzaslice, func(l, j int) bool { return pizzaslice[l].Aantal > pizzaslice[j].Aantal })

		var buf2 = new(bytes.Buffer)
		enc := json.NewEncoder(buf2)
		enc.Encode(pizzaslice)
		f2, err := os.Create("./Livepizza.json")
		if nil != err {
			log.Fatalln(err)
		}
		defer f2.Close()
		io.Copy(f2, buf2)

		toptienpizzas = pizzaslice[0:10]

		fmt.Println("Dit is het resultaat van top tien pizza's :", toptienpizzas)

	}

	if time.Now().Weekday() == 7 {
		var buf = new(bytes.Buffer)

		enc := json.NewEncoder(buf)
		enc.Encode(toptienpizzas)
		f, err := os.Create("./Toptienvorigeweek.json")
		if nil != err {
			log.Fatalln(err)
		}
		defer f.Close()
		io.Copy(f, buf)
	}

	readjson2, err := ioutil.ReadFile(file2)

	err = json.Unmarshal(readjson2, &resultaatpizza)
	if err != nil {
		panic(err)
	}

	fmt.Println("dit is resultaatpizza:", resultaatpizza)
	t.Execute(w, resultaatpizza)

	fmt.Println("encoded json toptien: ", toptienpizzas)
}

func logRequest(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Handling request for ", r.URL)
		h.ServeHTTP(w, r)
	})
}

func LoadPizzas() ([]Struct1, error) {
	var pizzas []Struct1

	raw, err := ioutil.ReadFile("./test1.json")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(raw, &pizzas)
	if err != nil {
		return nil, err
	}
	return pizzas, nil

}
