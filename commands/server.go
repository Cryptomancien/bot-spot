package commands

import (
	"fmt"
	"html/template"
	"log"
	"main/database"
	"net/http"
	"strconv"
)

var perPage int = 200

func getPage(r *http.Request) int {
	query := r.URL.Query()
	pageStr := query.Get("page")

	if pageStr == "" {
		return 1
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		return 1
	}

	return page
}

func Server() {

	fmt.Println("http://localhost:8080")
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		page := getPage(r)
		fmt.Println(page)

		//docs := database.ListPerPage(page, perPage)
		docs := database.List()
		cyclesCount := 0
		cyclesCompleted := 0
		totalBuy := 0.0
		totalSell := 0.0

		var cycles []map[string]interface{}
		for _, doc := range docs {
			quantity := doc.Get("quantity")
			buyPrice := doc.Get("buyPrice")
			sellPrice := doc.Get("sellPrice")
			status := doc.Get("status")

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
				"idInt":     doc.Get("idInt"),
				"exchange":  doc.Get("exchange"),
				"status":    doc.Get("status"),
				"quantity":  quantityStr,
				"buyPrice":  buyPriceFloat,
				"sellPrice": sellPriceFloat,
				"change":    percentageChange,
				"buyId":     doc.Get("buyId"),
				"sellId":    doc.Get("sellId"),
			})

			// Update stats
			cyclesCount++

			if status == "completed" {
				cyclesCompleted++

				b := (buyPrice).(float64) * (quantity).(float64)
				totalBuy += b

				s := (sellPrice).(float64) * (quantity).(float64)
				totalSell += s
			}
		}

		gainAbs := totalSell - totalBuy

		// Pagination

		tmpl, err := template.ParseFiles("commands/misc/template.html")
		if err != nil {
			http.Error(w, "Error loading template", http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, map[string]interface{}{
			"Cycles":          cycles,
			"cyclesCount":     cyclesCount,
			"cyclesCompleted": cyclesCompleted,
			"totalBuy":        totalBuy,
			"totalSell":       totalSell,
			"gainAbs":         gainAbs,
			"page":            page,
		})
		if err != nil {
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
		}
	})

	err := http.ListenAndServe("localhost:8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
