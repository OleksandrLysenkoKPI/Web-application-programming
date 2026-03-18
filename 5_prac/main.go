package main

import (
	"fmt"
	"html/template"
	"math"
	"net/http"
	"strconv"
)

type Data struct {
	N            string
	AcPrice      string
	PlPrice      string
	Results      *CalcResults
	ShowResults  bool
}

type CalcResults struct {
	Woc        string
	TvOc       string
	KaOc       string
	KpOc       string
	Wdk        string
	Wdc        string
	MathNedA   string
	MathNedP   string
	MathLosses string
}

func calculate(n, acPrice, plPrice float64) *CalcResults {
	wOc := 0.01 + 0.07 + 0.015 + 0.02 + 0.03*n
	tvOc := (0.01*30 + 0.07*10 + 0.015*100 + 0.02*15 + (0.03*n)*2) / wOc
	kaOc := (wOc * tvOc) / 8760
	kpOc := 1.2 * (43.0 / 8760.0)
	wDk := 2 * wOc * (kaOc + kpOc)
	wDc := wDk + 0.02

	mathNedA := 0.01 * 45 * math.Pow(10, -3) * 5.12 * math.Pow(10, 3) * 6451
	mathNedP := 4 * math.Pow(10, 3) * 5.12 * math.Pow(10, 3) * 6451
	mathLosses := acPrice*mathNedA + plPrice*mathNedP

	return &CalcResults{
		Woc:        fmt.Sprintf("%.4f", wOc),
		TvOc:       fmt.Sprintf("%.1f", tvOc),
		KaOc:       fmt.Sprintf("%.5f", kaOc),
		KpOc:       fmt.Sprintf("%.5f", kpOc),
		Wdk:        fmt.Sprintf("%.5f", wDk),
		Wdc:        fmt.Sprintf("%.4f", wDc),
		MathNedA:   fmt.Sprintf("%.2f", mathNedA),
		MathNedP:   fmt.Sprintf("%.2f", mathNedP),
		MathLosses: fmt.Sprintf("%.2f", mathLosses),
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	data := Data{}

	if r.Method == http.MethodPost {
		n, _ := strconv.ParseFloat(r.FormValue("connection"), 64)
		ac, _ := strconv.ParseFloat(r.FormValue("accident_price"), 64)
		pl, _ := strconv.ParseFloat(r.FormValue("planed_price"), 64)

		data.N = r.FormValue("connection")
		data.AcPrice = r.FormValue("accident_price")
		data.PlPrice = r.FormValue("planed_price")
		data.Results = calculate(n, ac, pl)
		data.ShowResults = true
	}

	tmpl.Execute(w, data)
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", indexHandler)
	fmt.Println("Сервер запущено на http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}