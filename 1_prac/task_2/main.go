package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

type OilData struct {
	C, H, O, S, Q, V, A, W float64
	Cr, Hr, Or, Sr, Ar, Wr float64
	Qr                     float64
	HasResult              bool
	Error                  string
}

func oilHandler(w http.ResponseWriter, r *http.Request) {
	data := OilData{}

	if r.Method == http.MethodPost {
		data.C, _ = strconv.ParseFloat(r.FormValue("C"), 64)
		data.H, _ = strconv.ParseFloat(r.FormValue("H"), 64)
		data.O, _ = strconv.ParseFloat(r.FormValue("O"), 64)
		data.S, _ = strconv.ParseFloat(r.FormValue("S"), 64)
		data.Q, _ = strconv.ParseFloat(r.FormValue("Q"), 64)
		data.V, _ = strconv.ParseFloat(r.FormValue("V"), 64)
		data.A, _ = strconv.ParseFloat(r.FormValue("A"), 64)
		data.W, _ = strconv.ParseFloat(r.FormValue("W"), 64)

		eq1 := (100 - data.V - data.A) / 100
		eq2 := (100 - data.V) / 100

		data.Cr = data.C * eq1
		data.Hr = data.H * eq1
		data.Or = data.O * eq1
		data.Sr = data.S * eq1
		data.Ar = data.A * eq2
		data.Wr = data.W * eq2

		data.Qr = data.Q*((100-data.V-data.Ar)/100) - 0.025*data.V

		data.HasResult = true
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

	http.HandleFunc("/", oilHandler)

	fmt.Println("Сервер запущено на http://localhost:8082")
	http.ListenAndServe("127.0.0.1:8082", nil)
}