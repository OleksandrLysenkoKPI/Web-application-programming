package main

import (
	"fmt"
	"html/template"
	"math"
	"net/http"
	"strconv"
)

type CalcData struct {
	Coal, Mazut, Gas     float64
	CoalEm, CoalTotal    string
	MazutEm, MazutTotal  string
	GasEm, GasTotal      string
	HasResults           bool
}

func calculate(w http.ResponseWriter, r *http.Request) {
	data := CalcData{}

	if r.Method == http.MethodPost {
		data.Coal, _ = strconv.ParseFloat(r.FormValue("coal"), 64)
		data.Mazut, _ = strconv.ParseFloat(r.FormValue("mazut"), 64)
		data.Gas, _ = strconv.ParseFloat(r.FormValue("gas"), 64)

		coalEmVal := (math.Pow(10, 6) / 20.47) * 0.8 * (25.2 / (100 - 1.5)) * (1 - 0.985)
		coalTotalVal := 1e-6 * coalEmVal * 20.47 * data.Coal

		mazutEmVal := (math.Pow(10, 6) / 39.48) * 1 * (0.15 / 100) * (1 - 0.985)
		mazutTotalVal := 1e-6 * mazutEmVal * 39.48 * data.Mazut

		gasEmVal := 0.0
		gasTotalVal := 0.0

		data.CoalEm = fmt.Sprintf("%.2f", coalEmVal)
		data.CoalTotal = fmt.Sprintf("%.2f", coalTotalVal)
		data.MazutEm = fmt.Sprintf("%.2f", mazutEmVal)
		data.MazutTotal = fmt.Sprintf("%.2f", mazutTotalVal)
		data.GasEm = fmt.Sprintf("%.2f", gasEmVal)
		data.GasTotal = fmt.Sprintf("%.2f", gasTotalVal)
		data.HasResults = true
	}

	tmpl, _ := template.ParseFiles("templates/index.html")
	tmpl.Execute(w, data)
}

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", calculate)

	fmt.Println("Сервер запущено на http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}