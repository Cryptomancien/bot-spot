package server

import (
	"fmt"
	"github.com/ostafen/clover"
	"html/template"
	"log"
	"net/http"
)

// todo dynamic path

func Serve() {
	// Open database once at the start
	db, err := clover.Open("cycles.db")
	if err != nil {
		log.Fatal("Error opening database:", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Println("Error closing database:", err)
		}
	}()

	fmt.Println("http://localhost:8080")
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		docs, err := db.Query("cycles").FindAll()
		if err != nil {
			http.Error(w, "Error fetching data", http.StatusInternalServerError)
			return
		}

		var cycles []map[string]interface{}
		for _, doc := range docs {
			cycles = append(cycles, map[string]interface{}{
				"_id":       doc.Get("_id"),
				"buyPrice":  doc.Get("buyPrice"),
				"sellPrice": doc.Get("sellPrice"),
				"exchange":  doc.Get("exchange"),
				"quantity":  doc.Get("quantity"),
				"status":    doc.Get("status"),
			})
		}

		tmpl, err := template.ParseFiles("server/index.html")
		if err != nil {
			http.Error(w, "Error loading template", http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, map[string]interface{}{"Cycles": cycles})
		if err != nil {
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
		}
	})

	err = http.ListenAndServe("localhost:8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
