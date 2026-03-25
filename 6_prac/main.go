package main

import (
	"fmt"
	"html/template"
	"math"
	"net/http"
	"strconv"
)

type PageData struct {
	Pn          string
	Kv          string
	Tg          string
	Results     *CalcResults
	ShowResults bool
}

type CalcResults struct {
	GroupKv             string
	EffectiveEP         string
	ActivePowerCoef     string
	EstActLoad          string
	EstReactLoad        string
	FullPower           string
	GroupElectricity    string
	TotalKv             string
	TotalEffectiveEP    string
	TotalActivePower    string
	EstTireActLoad      string
	EstTireReactLoad    string
	FullTirePower       string
	TireGroupElectricity string
}

func calculateLoads(pn, kv, tg float64) *CalcResults {
	productNominalSand := 4 * pn
	
	sumProdNomMultKv := (productNominalSand * 0.15) + 3.36 + 25.2 + 10.8 + 10 + (40 * kv) + 12.8 + 13
	sumProdNom := productNominalSand + 28 + 168 + 36 + 20 + 40 + 64 + 20
	squarePSumProdNom := (4 * math.Pow(pn, 2)) + 392 + 7056 + 1296 + 400 + 1600 + 2048 + 400
	sumProdNomMultKvTg := (productNominalSand * 0.15 * 1.33) + 3.36 + 33.5 + (36 * 0.3 * tg) + 7.5 + (40 * kv * 1) + 12.8 + 9.5

	groupKv := sumProdNomMultKv / sumProdNom
	effectiveEP := (math.Pow(sumProdNom, 2) / squarePSumProdNom) + 1
	estActLoad := 1.25 * sumProdNomMultKv
	estReactLoad := 1.0 * sumProdNomMultKvTg
	fullPower := math.Sqrt(math.Pow(estActLoad, 2) + math.Pow(estReactLoad, 2))
	groupElectricity := estActLoad / 0.38

	totalKv := 752.0 / 2330.0
	totalEffectiveEP := math.Pow(2330, 2) / 96399.0
	estTireActLoad := 0.7 * 752.0
	estTireReactLoad := 0.7 * 657.0
	fullTirePower := math.Sqrt(math.Pow(estTireActLoad, 2) + math.Pow(estTireReactLoad, 2))
	tireGroupElectricity := estTireActLoad / 0.38

	return &CalcResults{
		GroupKv:              fmt.Sprintf("%.4f", groupKv),
		EffectiveEP:          fmt.Sprintf("%.0f", effectiveEP),
		ActivePowerCoef:      "1.25",
		EstActLoad:           fmt.Sprintf("%.2f", estActLoad),
		EstReactLoad:         fmt.Sprintf("%.2f", estReactLoad),
		FullPower:            fmt.Sprintf("%.3f", fullPower),
		GroupElectricity:     fmt.Sprintf("%.2f", groupElectricity),
		TotalKv:              fmt.Sprintf("%.2f", totalKv),
		TotalEffectiveEP:     fmt.Sprintf("%.0f", totalEffectiveEP),
		TotalActivePower:     "0.7",
		EstTireActLoad:       fmt.Sprintf("%.1f", estTireActLoad),
		EstTireReactLoad:     fmt.Sprintf("%.1f", estTireReactLoad),
		FullTirePower:        fmt.Sprintf("%.0f", fullTirePower),
		TireGroupElectricity: fmt.Sprintf("%.2f", tireGroupElectricity),
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	data := PageData{}

	if r.Method == http.MethodPost {
		pn, _ := strconv.ParseFloat(r.FormValue("power"), 64)
		kv, _ := strconv.ParseFloat(r.FormValue("coef"), 64)
		tg, _ := strconv.ParseFloat(r.FormValue("tan"), 64)

		data.Pn = r.FormValue("power")
		data.Kv = r.FormValue("coef")
		data.Tg = r.FormValue("tan")
		data.Results = calculateLoads(pn, kv, tg)
		data.ShowResults = true
	}

	tmpl.Execute(w, data)
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", handler)
	fmt.Println("Сервер запущено на http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}