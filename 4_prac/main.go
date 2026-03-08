package main

import (
	"fmt"
	"html/template"
	"math"
	"net/http"
	"strconv"
)

type Results struct {
	Task1 struct {
		Im, Impa, Sek, Smin string
	}
	Task2 struct {
		Xc, Xt, Xsum, Ip0 string
	}
	Task3 struct {
		Ish3, Ish2, Ish3min, Ish2min     string
		Ishn3, Ishn2, Ishn3min, Ishn2min string
		Iln3, Iln2, Iln3min, Iln2min     string
	}
	Calculated bool
}

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", handler)

	fmt.Println("Сервер запущено: http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	data := Results{}

	if r.Method == http.MethodPost {
		r.ParseForm()
		data.Calculated = true

		sm, _ := strconv.ParseFloat(r.FormValue("s_m"), 64)
		ik, _ := strconv.ParseFloat(r.FormValue("i_k"), 64)
		tf, _ := strconv.ParseFloat(r.FormValue("t_f"), 64)
		im := (sm / 2) / (math.Sqrt(3) * 10)
		data.Task1.Im = fmt.Sprintf("%.1f", im)
		data.Task1.Impa = fmt.Sprintf("%.0f", 2*im)
		data.Task1.Sek = fmt.Sprintf("%.1f", im/1.4)
		data.Task1.Smin = fmt.Sprintf("%.0f", (ik*math.Sqrt(tf))/92)

		sk, _ := strconv.ParseFloat(r.FormValue("s_k"), 64)
		xc := math.Pow(10.5, 2) / sk
		xt2 := (10.5 / 100) * (math.Pow(10.5, 2) / 6.3)
		xsum := xc + xt2
		data.Task2.Xc, data.Task2.Xt = fmt.Sprintf("%.2f", xc), fmt.Sprintf("%.2f", xt2)
		data.Task2.Xsum, data.Task2.Ip0 = fmt.Sprintf("%.2f", xsum), fmt.Sprintf("%.1f", 10.5/(math.Sqrt(3)*xsum))

		rsn, _ := strconv.ParseFloat(r.FormValue("r_sn"), 64)
		xsn, _ := strconv.ParseFloat(r.FormValue("x_sn"), 64)
		rsnMin, _ := strconv.ParseFloat(r.FormValue("r_sn_min"), 64)
		xsnMin, _ := strconv.ParseFloat(r.FormValue("x_sn_min"), 64)
		xt3 := (11.1 * math.Pow(115, 2)) / (100 * 6.3)
		
		ish3 := (115 * 1000) / (math.Sqrt(3) * math.Sqrt(math.Pow(rsn, 2)+math.Pow(xsn+xt3, 2)))
		data.Task3.Ish3, data.Task3.Ish2 = fmt.Sprintf("%.1f", ish3), fmt.Sprintf("%.0f", ish3*0.866)
		ish3min := (115 * 1000) / (math.Sqrt(3) * math.Sqrt(math.Pow(rsnMin, 2)+math.Pow(xsnMin+xt3, 2)))
		data.Task3.Ish3min, data.Task3.Ish2min = fmt.Sprintf("%.0f", ish3min), fmt.Sprintf("%.0f", ish3min*0.866)

		kpr := math.Pow(11, 2) / math.Pow(115, 2)
		zshn := math.Sqrt(math.Pow(rsn*kpr, 2) + math.Pow((xsn+xt3)*kpr, 2))
		ishn3 := (11 * 1000) / (math.Sqrt(3) * zshn)
		data.Task3.Ishn3, data.Task3.Ishn2 = fmt.Sprintf("%.0f", ishn3), fmt.Sprintf("%.0f", ishn3*0.866)
		
		zshnMin := math.Sqrt(math.Pow(rsnMin*kpr, 2) + math.Pow((xsnMin+xt3)*kpr, 2))
		ishn3min := (11 * 1000) / (math.Sqrt(3) * zshnMin)
		data.Task3.Ishn3min, data.Task3.Ishn2min = fmt.Sprintf("%.0f", ishn3min), fmt.Sprintf("%.0f", ishn3min*0.866)

		il := 12.42
		rl, xl := il*0.64, il*0.363
		zSumN := math.Sqrt(math.Pow(rl+rsn*kpr, 2) + math.Pow(xl+(xsn+xt3)*kpr, 2))
		iln3 := (11 * 1000) / (math.Sqrt(3) * zSumN)
		data.Task3.Iln3, data.Task3.Iln2 = fmt.Sprintf("%.0f", iln3), fmt.Sprintf("%.0f", iln3*0.866)
		
		zSumNMin := math.Sqrt(math.Pow(rl+rsnMin*kpr, 2) + math.Pow(xl+(xsnMin+xt3)*kpr, 2))
		iln3min := (11 * 1000) / (math.Sqrt(3) * zSumNMin)
		data.Task3.Iln3min, data.Task3.Iln2min = fmt.Sprintf("%.0f", iln3min), fmt.Sprintf("%.0f", iln3min*0.866)
	}

	tmpl, _ := template.ParseFiles("templates/index.html")
	tmpl.Execute(w, data)
}