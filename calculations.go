package main

import "log"

func calculateTotal(o Order) float64 {
	var total float64

	for _, l := range o.Lines {
		p := l.CustomPizza.Ingredienten
		if len(p) == 0 {
			p, err := GetPizza(l.Pizza)
			total += float64(p.Price * float64(l.Qty))
			if err != nil {
				log.Printf("Pizza %v not found, total price may be incorrect\n", l.Pizza)
			}
		} else if len(p) != 0 {
			for _, s := range p {
				total += float64(s.Price * float64(l.Qty))
			}
		}
	}
	return total
}
