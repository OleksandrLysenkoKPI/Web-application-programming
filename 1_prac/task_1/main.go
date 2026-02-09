package main

import (
	"fmt"
	"html/template"
	"math"
	"net/http"
	"strconv"
)

type FuelData struct {
	H, C, S, N, O, W, A float64
	Kpc, Kpg            float64
	Hc, Oc, Nc, Ac, Sc, Cc float64
	Hg, Og, Ng, Sg, Cg     float64
	Qph, Qch, Qgh          float64
	HasResult              bool
	Error                  string
}

func round(val float64) float64 {
	return math.Round(val*100) / 100
}

func fuelHandler(w http.ResponseWriter, r *http.Request) {
	data := FuelData{}

	if r.Method == http.MethodPost {
		data.H, _ = strconv.ParseFloat(r.FormValue("H"), 64)
		data.C, _ = strconv.ParseFloat(r.FormValue("C"), 64)
		data.S, _ = strconv.ParseFloat(r.FormValue("S"), 64)
		data.N, _ = strconv.ParseFloat(r.FormValue("N"), 64)
		data.O, _ = strconv.ParseFloat(r.FormValue("O"), 64)
		data.W, _ = strconv.ParseFloat(r.FormValue("W"), 64)
		data.A, _ = strconv.ParseFloat(r.FormValue("A"), 64)

		sum := data.H + data.C + data.S + data.N + data.O + data.W + data.A
		if math.Abs(sum-100) > 0.01 {
			data.Error = "Елементарний склад робочого палива не дорівнює 100%!"
		} else {
			data.Kpc = 100 / (100 - data.W)
			data.Kpg = 100 / (100 - data.W - data.A)

			data.Hc = data.H * data.Kpc
			data.Oc = data.O * data.Kpc
			data.Nc = data.N * data.Kpc
			data.Ac = data.A * data.Kpc
			data.Sc = data.S * data.Kpc
			data.Cc = data.C * data.Kpc

			data.Cg = data.C * data.Kpg
			data.Og = data.O * data.Kpg
			data.Ng = data.N * data.Kpg
			data.Sg = data.S * data.Kpg
			data.Hg = data.H * data.Kpg

			data.Qph = (339*data.C + 1030*data.H - 108.8*(data.O-data.S) - 25*data.W) / 1000
			data.Qch = (data.Qph + 0.025*data.W) * (100 / (100 - data.W))
			data.Qgh = (data.Qph + 0.025*data.W) * (100 / (100 - data.W - data.A))

			data.HasResult = true
		}
	}

    tmpl, err := template.ParseFiles("index.html")
    if err != nil {
        http.Error(w, "Файл шаблону не знайдено", http.StatusInternalServerError)
        return
    }

    tmpl.Execute(w, data)
}

func main() {
    fs := http.FileServer(http.Dir("./static"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))

    http.HandleFunc("/", fuelHandler)

    fmt.Println("Сервер запущено на http://localhost:8081")
    http.ListenAndServe("127.0.0.1:8081", nil)
}