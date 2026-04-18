package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/xuri/excelize/v2"
)

// Data represents a row of data from the excel file
type Data struct {
	ID      string
	Name    string
	TallyIn string
	Diff    string
}

// PageData is passed to the HTML template
type PageData struct {
	Results []Data
	Query   string
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/search", searchHandler)
	//http.HandleFunc("/update", updateHandler)

	fmt.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// Panggil searchExcel dengan query kosong untuk mengambil semua data
	results, err := searchExcel("")
	if err != nil {
		http.Error(w, "Error reading database: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, PageData{
		Results: results,
		Query:   "",
	})
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	results, err := searchExcel(query)
	if err != nil {
		http.Error(w, "Error reading database: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, PageData{
		Results: results,
		Query:   query,
	})
}

func searchExcel(query string) ([]Data, error) {
	// Open the excel file
	f, err := excelize.OpenFile("database.xlsx")
	if err != nil {
		return nil, fmt.Errorf("make sure database.xlsx exists with 'Sheet1': %v", err)
	}
	defer f.Close()

	// Get all rows in the first sheet
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return nil, err
	}

	var results []Data
	if len(rows) == 0 {
		return results, nil
	}

	// Skip header row (i=0)
	for i := 1; i < len(rows); i++ {
		row := rows[i]
		if len(row) < 2 {
			continue
		}

		// Simple case-insensitive search in any column
		match := false
		if query == "" {
			match = true
		} else {
			for _, col := range row {
				if strings.Contains(strings.ToLower(col), strings.ToLower(query)) {
					match = true
					break
				}
			}
		}

		if match {
			// Ensure we don't go out of bounds if some columns are missing
			tallyIn := ""
			diff := ""
			if len(row) >= 3 {
				tallyIn = row[2]
			}
			if len(row) >= 4 {
				diff = row[3]
			}

			results = append(results, Data{
				ID:      row[0],
				Name:    row[1],
				TallyIn: tallyIn,
				Diff:    diff,
			})
		}
	}

	return results, nil
}
