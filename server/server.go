package server

import (
	"fmt"
	"html/template"
	"log"
	"main/database"
	"net/http"
)

func Serve() {

	fmt.Println("http://localhost:8080")
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		docs := database.List()

		var cycles []map[string]interface{}
		for _, doc := range docs {
			quantity := doc.Get("quantity")
			buyPrice := doc.Get("buyPrice")
			sellPrice := doc.Get("sellPrice")

			var quantityFloat, buyPriceFloat, sellPriceFloat float64
			var quantityStr string

			if q, ok := quantity.(float64); ok {
				quantityFloat = q
				quantityStr = fmt.Sprintf("%.8f", q)
			} else {
				quantityStr = fmt.Sprintf("%v", quantity)
			}

			if bp, ok := buyPrice.(float64); ok {
				buyPriceFloat = bp
			}

			if sp, ok := sellPrice.(float64); ok {
				sellPriceFloat = sp
			}

			var percentageChange string
			if buyPriceFloat > 0 {
				change := ((quantityFloat * sellPriceFloat) - (quantityFloat * buyPriceFloat)) / (quantityFloat * buyPriceFloat) * 100
				percentageChange = fmt.Sprintf("%.2f%%", change)
			} else {
				percentageChange = "N/A"
			}

			cycles = append(cycles, map[string]interface{}{
				"_id":       doc.Get("_id"),
				"exchange":  doc.Get("exchange"),
				"status":    doc.Get("status"),
				"quantity":  quantityStr,
				"buyPrice":  buyPriceFloat,
				"sellPrice": sellPriceFloat,
				"change":    percentageChange,
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

	err := http.ListenAndServe("localhost:8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
