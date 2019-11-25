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

func orderpizzaHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Preparing pizza order page")

	pp, _ := LoadPizzas()
	// prepare the page
	t, err := template.ParseFiles("./templates/orderpizza.html")
	if err != nil {
		log.Println(err)
	}
	t.Execute(w, pp)

}

func getOrderNo() int {
	oo, _ := LoadOrders()
	if len(oo) == 0 {
		return 1
	}
	sort.Slice(oo, func(i, j int) bool {
		return oo[i].No < oo[j].No
	})

	return oo[len(oo)-1].No + 1
}

func opslaantotaalorderHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./html/thankyou.html")
	log.Println("Ordering pizza")
	r.ParseForm()
	c := Customer{Email: r.FormValue("email"), Name: r.FormValue("name"), Address: r.FormValue("address"), Postalcode: r.FormValue("postalcode"), PhoneNumber: r.FormValue("phoneNumber")}

	// create the customer
	SaveCustomer(c)
	d := time.Now()
	o := Order{No: getOrderNo(), Date: d, DateString: d.Format(time.RFC1123Z), Email: c.Email}

	readjson3, err := ioutil.ReadFile("./data/pizzas.json")

	var pizzaslice2 []Pizza

	err = json.Unmarshal(readjson3, &pizzaslice2)
	if err != nil {
		panic(err)
	}
	for i := 0; i < len(pizzaslice2); i++ {
		pp := r.FormValue("pizza" + pizzaslice2[i].Name)
		qq := r.FormValue("qty" + pizzaslice2[i].Name)
		if pp != "" && qq != "" {
			qty, _ := strconv.Atoi(qq)
			o.Lines = append(o.Lines, OrderLine{Pizza: pp, Qty: qty})
		}
	}
	o.Total = calculateTotal(o)

	// create the order
	SaveOrder(o)

	log.Printf("Placing order %v for customer: %v\n ", o, c)

	file1 := "./data/Livepizza.json"
	file2 := "./data/Toptienvorigeweek.json"

	readjson, err := ioutil.ReadFile(file1)

	var pizzaslice []Pizza
	var toptienpizzas []Pizza
	var resultaatpizza []Pizza

	err = json.Unmarshal(readjson, &pizzaslice)
	if err != nil {
		panic(err)
	}

	var status int
	for i := 0; i < len(pizzaslice); i++ {

		naamvanpizza := r.FormValue("pizza" + pizzaslice[i].Name)

		for _, value := range pizzaslice {
			if value.Name == naamvanpizza {
				fmt.Println("test:" + value.Name)
				status = value.Aantal

				fmt.Println("De huidige aantal van ", naamvanpizza, " is:", status)
			}

		}

		aantalint, _ := strconv.Atoi(r.FormValue("qty" + pizzaslice[i].Name))
		fmt.Println("De hoeveelheid ", naamvanpizza, " Pizza's die erbij gekomen zijn is:", r.FormValue("qty"+pizzaslice[i].Name))
		k := 1
		var nieuweaantal int
		if time.Now().Weekday() == 1 {
			k = 2
		}
		switch k {
		case 1:
			nieuweaantal = aantalint + status
		case 2:
			nieuweaantal = 0
		}
		prijsfloat, _ := strconv.ParseFloat(r.FormValue("price"+pizzaslice[i].Name), 64)

		pizzaslice[i].Aantal = nieuweaantal
		pizzaslice[i].Image = r.FormValue("image" + pizzaslice[i].Name)
		pizzaslice[i].Price = prijsfloat

		fmt.Println("Toptienvorigeweek aantal pizza van", naamvanpizza, " is:", nieuweaantal)

		sort.SliceStable(pizzaslice, func(l, j int) bool { return pizzaslice[l].Aantal > pizzaslice[j].Aantal })

		var buf2 = new(bytes.Buffer)
		enc := json.NewEncoder(buf2)
		enc.Encode(pizzaslice)
		f2, err := os.Create("./data/Livepizza.json")
		if nil != err {
			log.Fatalln(err)
		}
		defer f2.Close()
		io.Copy(f2, buf2)

		toptienpizzas = pizzaslice[0:10]

		fmt.Println("Dit is het resultaat van top tien pizza's :", toptienpizzas)

	}

	if time.Now().Weekday() == 0 {
		var buf = new(bytes.Buffer)

		enc := json.NewEncoder(buf)
		enc.Encode(toptienpizzas)
		f, err := os.Create("./data/Toptienvorigeweek.json")
		if nil != err {
			log.Fatalln(err)
		}
		defer f.Close()
		io.Copy(f, buf)
		fmt.Println("Toptien pizza is geupdated!")
	}

	readjson2, err := ioutil.ReadFile(file2)

	err = json.Unmarshal(readjson2, &resultaatpizza)
	if err != nil {
		panic(err)
	}

	fmt.Println("dit is resultaatpizza:", resultaatpizza)
	t.Execute(w, resultaatpizza)
}
