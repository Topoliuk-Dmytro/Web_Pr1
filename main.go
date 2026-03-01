package main

import (
	"fuel-calculator/fuel"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

var tmpl *template.Template

func main() {
	var err error
	tmpl, err = template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/task1", task1Handler)
	http.HandleFunc("/task2", task2Handler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Println("Сервер запущено на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	data := map[string]interface{}{
		"Variant2": fuel.WorkingMassComposition{
			HP: 4.2,
			CP: 62.1,
			SP: 3.2,
			NP: 20.6,
			OP: 6.4,
			WP: 0.7,
			AP: 2.8,
		},
	}
	tmpl.ExecuteTemplate(w, "index.html", data)
}

func task1Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не дозволено", http.StatusMethodNotAllowed)
		return
	}

	wc := fuel.WorkingMassComposition{}
	parseFloat := func(name string) float64 {
		s := r.FormValue(name)
		v, _ := strconv.ParseFloat(s, 64)
		return v
	}
	wc.HP = parseFloat("hp")
	wc.CP = parseFloat("cp")
	wc.SP = parseFloat("sp")
	wc.NP = parseFloat("np")
	wc.OP = parseFloat("op")
	wc.WP = parseFloat("wp")
	wc.AP = parseFloat("ap")

	result, err := fuel.Task1Calculate(wc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data := map[string]interface{}{
		"Input":  wc,
		"Result": result,
	}
	tmpl.ExecuteTemplate(w, "task1_result.html", data)
}

func task2Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не дозволено", http.StatusMethodNotAllowed)
		return
	}

	in := fuel.CombustibleMassInput{}
	parseFloat := func(name string) float64 {
		s := r.FormValue(name)
		v, _ := strconv.ParseFloat(s, 64)
		return v
	}
	in.Cg = parseFloat("cg")
	in.Hg = parseFloat("hg")
	in.Og = parseFloat("og")
	in.Sg = parseFloat("sg")
	in.QgH = parseFloat("qgh")
	in.Wr = parseFloat("wr")
	in.Ad = parseFloat("ad")
	in.Vg = parseFloat("vg")

	result, err := fuel.Task2Calculate(in)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data := map[string]interface{}{
		"Input":  in,
		"Result": result,
	}
	tmpl.ExecuteTemplate(w, "task2_result.html", data)
}
