package main

import (
	"fmt"
	"html/template"
	"math"
	"net/http"
	"strconv"
)

var tmpls = template.Must(template.ParseFiles("index.html"))

type Result struct {
	NormLaw      string
	DeltaW1      string
	W1           string
	Income1      string
	W2           string
	Fine1        string
	DeltaW2      string
	W3           string
	Income2      string
	W4           string
	Fine2        string
	IncomeFinale string
	HasResult    bool
}

func calculatePd(p, pc, sigma float64) float64 {
	return (1 / (sigma * math.Sqrt(2 * math.Pi))) *
		math.Exp(-math.Pow(p-pc, 2) / (2 * math.Pow(sigma, 2)))
}

func integrate(pc, sigma, start, end float64, steps int) float64 {
	step := (end - start) / float64(steps)
	sum := 0.5 * (calculatePd(start, pc, sigma) + calculatePd(end, pc, sigma))

	for i := 1; i < steps; i++ {
		sum += calculatePd(start+float64(i)*step, pc, sigma)
	}
	return sum * step
}

func calculateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpls.Execute(w, Result{HasResult: false})
		return
	}

	pc, _ := strconv.ParseFloat(r.FormValue("avr_power"), 64)
	sigma1, _ := strconv.ParseFloat(r.FormValue("sigma_val_1"), 64)
	price, _ := strconv.ParseFloat(r.FormValue("electricity_price"), 64)
	
	const sigma2 = 0.25
	const pTarget = 5.0
	const steps = 1000

	deltaW1 := integrate(pc, sigma1, 4.75, 5.25, steps)
	w1 := pc * 24 * deltaW1
	income1 := w1 * price
	w2 := pc * 24 * (1 - deltaW1)
	fine1 := w2 * price

	deltaW2 := integrate(pc, sigma2, 4.75, 5.25, steps)
	w3 := pc * 24 * deltaW2
	income2 := w3 * price
	w4 := pc * 24 * (1 - deltaW2)
	fine2 := w4 * price

	finalIncome := income2 - fine2

	res := Result{
		NormLaw:      fmt.Sprintf("%.2f", calculatePd(pTarget, pc, sigma1)),
		DeltaW1:      fmt.Sprintf("%.2f", deltaW1*100),
		W1:           fmt.Sprintf("%.0f", w1),
		Income1:      fmt.Sprintf("%.0f", income1),
		W2:           fmt.Sprintf("%.0f", w2),
		Fine1:        fmt.Sprintf("%.0f", fine1),
		DeltaW2:      fmt.Sprintf("%.2f", deltaW2*100),
		W3:           fmt.Sprintf("%.1f", w3),
		Income2:      fmt.Sprintf("%.1f", income2),
		W4:           fmt.Sprintf("%.1f", w4),
		Fine2:        fmt.Sprintf("%.1f", fine2),
		IncomeFinale: fmt.Sprintf("%.1f", finalIncome),
		HasResult:    true,
	}

	tmpls.Execute(w, res)
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	
	http.HandleFunc("/", calculateHandler)

	fmt.Println("Сервер запущено на http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Помилка старту сервера: %v\n", err)
	}
}