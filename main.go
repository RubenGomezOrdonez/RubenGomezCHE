package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	deletecsvFile()
	CreateCSV()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logfile := os.Getenv("LOGFILE")
	if logfile == "" {
		logfile = "/var/log/golang/golang-server.log"
	}
	f, _ := os.Create(logfile)
	defer f.Close()
	log.SetOutput(f)

	log.Println("Starting app")

	http.HandleFunc("/view/pizzas", pizzaHandler)
	http.HandleFunc("/view/orders", orderHandler)
	http.HandleFunc("/view/orderdetails", orderdetailsHandler)
	http.HandleFunc("/view/customers", customerHandler)
	http.HandleFunc("/view/orderpizza", orderpizzaHandler)
	http.HandleFunc("/view/custompizza", custompizzaHandler)
	http.HandleFunc("/view/nieuwpizza", nieuwpizzaHandler)
	http.HandleFunc("/view/bewerkpagina", bewerkpaginaHandler)
	http.HandleFunc("/view/nieuwpizzaterugkoppeling", nieuwpizzaterugkoppelingHandler)
	http.HandleFunc("/view/bewerkpizza", bewerkpizzaHandler)
	http.HandleFunc("/view/thankyou", opslaantotaalorderHandler)
	http.HandleFunc("/view/toptien", toptienHandler)
	http.HandleFunc("/view/pizzadelete", Pizzadelete)
	http.HandleFunc("/view/login", Loginhandler)
	http.HandleFunc("/view/nieuwingredient", nieuwingredient)
	http.HandleFunc("/view/nieuwingredientterugkoppeling", nieuwingredientterugkoppeling)
	http.HandleFunc("/view/verwijderingredient", verwijderingredient)
	http.HandleFunc("/", redirecthtml)
	http.Handle("/html/", http.StripPrefix("/html/", http.FileServer(http.Dir("./html"))))
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./images"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./css"))))
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))
	http.Handle("/data/", http.StripPrefix("/data/", http.FileServer(http.Dir("./data"))))


	log.Fatal(http.ListenAndServe(":"+port, nil))
}
