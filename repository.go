package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

const (
	CUSTOMERFILE   = "./data/customers.json"
	PIZZAFILE      = "./data/pizzas.json"
	ORDERFILE      = "./data/orders.json"
	INGREDIENTFILE = "./data/ingredients.json"
	TOPTIENPIZZAS  = "./data/Toptienvorigeweek.json"
	LIVEPIZZAS     = "./data/Livepizza.json"
)

func loadData(f string, vv interface{}) error {
	raw, err := ioutil.ReadFile(f)
	if err != nil {
		return err
	}
	err = json.Unmarshal(raw, &vv)
	if err != nil {
		return err
	}
	return nil
}
func GetPizza(name string) (Pizza, error) {
	pp, _ := LoadPizzas()
	for _, p := range pp {
		if p.Name == name {
			return p, nil
		}
	}
	return Pizza{}, errors.New("pizza not found")
}

func LoadPizzas() ([]Pizza, error) {
	var pp []Pizza
	err := loadData(PIZZAFILE, &pp)

	return pp, err
}
func LoadIngredients() (map[string][]Ingredient, error) {
	var ii []Ingredient
	err := loadData(INGREDIENTFILE, &ii)

	var ingredientsMap = make(map[string][]Ingredient)
	for _, i := range ii {
		ingredientsMap[i.Type] = append(ingredientsMap[i.Type], i)
	}
	return ingredientsMap, err
}
func LoadCustomers() ([]Customer, error) {
	var cc []Customer
	err := loadData(CUSTOMERFILE, &cc)

	return cc, err
}
func LoadCustomersWithOrders() ([]CustomerWithOrders, error) {
	cco := []CustomerWithOrders{}
	cc, err := LoadCustomers()
	if err != nil {
		return nil, err
	}
	for _, c := range cc {
		co, _ := LoadOrdersForCustomer(c.Email)
		cco = append(cco, CustomerWithOrders{Customer: c, Orders: co})
	}
	return cco, nil
}
func LoadOrders() ([]Order, error) {
	var oo []Order
	err := loadData(ORDERFILE, &oo)

	return oo, err
}
func LoadOrdersForCustomer(email string) ([]Order, error) {
	oo, err := LoadOrders()
	if err != nil {
		return nil, err
	}
	var co []Order
	for _, o := range oo {
		if o.Email == email {
			co = append(co, o)
		}
	}
	return co, nil
}

func saveData(f string, vv interface{}) error {
	raw, err := json.Marshal(vv)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(f, raw, 0644)
}

func SaveCustomer(c Customer) error {
	log.Println("Saving customer: ", c)
	if c.Email == "" {
		return errors.New("customer not defined")
	}
	cc, err := LoadCustomers()
	if err != nil {
		return err
	}
	for _, ec := range cc {
		if ec.Email == c.Email {
			return errors.New("customer already exists")
		}
	}
	cc = append(cc, c)
	return saveData(CUSTOMERFILE, cc)
}
func SaveOrder(o Order) error {
	log.Println("Saving order: ", o)
	if o.No == 0 || o.Email == "" {
		return errors.New("order not defined")
	}
	oo, err := LoadOrders()
	if err != nil {
		return err
	}
	for _, eo := range oo {
		if eo.No == o.No {
			return errors.New("order already exists")
		}
	}
	oo = append(oo, o)
	return saveData(ORDERFILE, oo)
}
func DeleteCustomer(email string) error {
	cc, err := LoadCustomers()
	if err != nil {
		return err
	}
	for i, ec := range cc {
		if ec.Email == email {
			// found ==> delete, save, return
			newcc := append(cc[:i], cc[i+1:]...)
			return saveData(CUSTOMERFILE, newcc)
		}
	}
	return nil
}
func DeleteOrder(no int) error {
	oo, err := LoadOrders()
	if err != nil {
		return err
	}
	for i, eo := range oo {
		if eo.No == no {
			// found ==> delete, save, return
			newoo := append(oo[:i], oo[i+1:]...)
			return saveData(ORDERFILE, newoo)
		}
	}
	return nil
}

func DeletePizza(delete string) error {
	cc, err := LoadPizzas()
	if err != nil {
		return err
	}
	for i, ec := range cc {
		if ec.Name == delete {
			fmt.Println(ec.Name)
			// found ==> delete, save, return
			newcc := append(cc[:i], cc[i+1:]...)
			return saveData(PIZZAFILE, newcc)
		}
		deletecsvFile()
		CreateCSV()
	}
	return nil

}

func DeletePizza2(delete string) error {
	cc, err := LoadLivePizzas()
	if err != nil {
		return err
	}

	for i, ec := range cc {
		if ec.Name == delete {
			fmt.Println(ec.Name)
			// found ==> delete, save, return
			newcc := append(cc[:i], cc[i+1:]...)
			return saveData(LIVEPIZZAS, newcc)
		}
	}

	return nil
}

func DeletePizza3(delete string) error {
	cc, err := LoadToptienPizzas()
	if err != nil {
		return err
	}

	for i, ec := range cc {
		if ec.Name == delete {
			fmt.Println(ec.Name)
			// found ==> delete, save, return
			newcc := append(cc[:i], cc[i+1:]...)
			return saveData(TOPTIENPIZZAS, newcc)
		}
	}

	return nil
}

func FindIngredient(name string, ingredients []Ingredient) (Ingredient, error) {
	for _, ingredient := range ingredients {
		if ingredient.Name == name {
			return ingredient, nil
		}
	}
	return Ingredient{}, errors.New("ingredient doesnt exist")

}

func LoadToptienPizzas() ([]Pizza, error) {
	var pp []Pizza
	err := loadData(TOPTIENPIZZAS, &pp)

	return pp, err
}

func LoadLivePizzas() ([]Pizza, error) {
	var pp []Pizza
	err := loadData(LIVEPIZZAS, &pp)

	return pp, err
}

func CreateCSV() {
	pizza, _ := LoadPizzas()
	var naam string
	var prijs float64
	f, err := os.OpenFile("./data/pizzas.csv", os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	w := csv.NewWriter(f)
	for i := 0; i < len(pizza); i++ {

		naam = pizza[i].Name
		prijs = pizza[i].Price

		s := fmt.Sprintf("%.2f", prijs)
		w.Write([]string{naam, s})
	}
	w.Flush()
	fmt.Println("==> done making file")
}
func nieuwingredient(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./templates/nieuwingredient.html")

	ingredient, _ := LoadIngredients()

	t.Execute(w, ingredient)
}
func nieuwingredientterugkoppeling(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./templates/nieuwingredientterugkoppeling.html")
	var ingredientslice []Ingredient

	readjson, err := ioutil.ReadFile(INGREDIENTFILE)
	err = json.Unmarshal(readjson, &ingredientslice)
	if err != nil {
		panic(err)
	}
	nieuwingredienttype := r.FormValue("type")

	nieuwingredientname := r.FormValue("name")
	nieuwingredientprice, _ := strconv.ParseFloat(r.FormValue("price"), 64)

	resultaat := Ingredient{Type: nieuwingredienttype, Name: nieuwingredientname, Price: nieuwingredientprice}

	ingredientslice = append(ingredientslice, resultaat)
	fmt.Println(ingredientslice)
	raw, err := json.MarshalIndent(ingredientslice, "", "\t")
	if err != nil {
		panic(err)
	}
	ioutil.WriteFile(INGREDIENTFILE, raw, 0644)

	laatsteingredient := ingredientslice[len(ingredientslice)-1]
	t.Execute(w, laatsteingredient)
}

func verwijderingredient(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./templates/verwijderingredient.html")
	var ingredientslice []Ingredient

	readjson, err := ioutil.ReadFile(INGREDIENTFILE)
	err = json.Unmarshal(readjson, &ingredientslice)
	if err != nil {
		panic(err)
	}
	naam := r.FormValue("naam")
	for i, ec := range ingredientslice {

		if ec.Name == naam {
			fmt.Println(ec.Name)
			// found ==> delete, save, return
			newcc := append(ingredientslice[:i], ingredientslice[i+1:]...)
			saveData(INGREDIENTFILE, newcc)
		}
	}
	t.Execute(w, naam)
}
func deletecsvFile() {
	// delete file
	var err = os.Remove("./data/pizzas.csv")
	if err != nil {
		return
	}

	fmt.Println("==> done deleting file")
}
