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
	Stock   string
	TallyIn string
	Diff    string
}

// PageData is passed to the HTML template
type PageData struct {
	Results []Data
	Query   string
}

func main() {
	// Route handler
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/search", searchHandler)
	http.HandleFunc("/update", updateHandler) // Handler untuk input data tally

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

// updateHandler menangani pengiriman data dari form web untuk disimpan ke file Excel
func updateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Metode tidak diizinkan", http.StatusMethodNotAllowed)
		return
	}

	id := r.FormValue("id")
	tallyInVal := r.FormValue("tallyin")

	if id == "" {
		http.Error(w, "ID harus diisi", http.StatusBadRequest)
		return
	}

	f, err := excelize.OpenFile("database.xlsx")
	if err != nil {
		http.Error(w, "Gagal membuka file Excel: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		http.Error(w, "Gagal membaca baris Excel", http.StatusInternalServerError)
		return
	}

	found := false
	var finalDiff float64
	for i, row := range rows {
		if i == 0 {
			continue // Skip header
		}
		if len(row) > 0 && row[0] == id {
			// 1. Update Tally-In di Kolom C
			cellTally := fmt.Sprintf("C%d", i+1)
			f.SetCellValue("Sheet1", cellTally, tallyInVal)

			// 2. HITUNG SELISIH (Tally-In - Stok Master)
			var targetQty, actualQty float64

			// Ekstrak angka stok master dari Kolom B dengan lebih teliti
			if len(row) >= 2 {
				rawB := row[1]
				if strings.Contains(rawB, "=") {
					parts := strings.Split(rawB, "=")
					val := strings.TrimSpace(parts[1])
					val = strings.TrimRight(val, " mc")
					val = strings.TrimSpace(val)
					fmt.Sscanf(val, "%f", &targetQty)
				} else {
					fmt.Sscanf(rawB, "%f", &targetQty)
				}
			}

			// Ambil nilai Tally-In dari input (pastikan bersih)
			fmt.Sscanf(tallyInVal, "%f", &actualQty)

			// Rumus: Tally In - Stok Master
			finalDiff = actualQty - targetQty

			// 3. Update Selisih di Kolom D
			cellDiff := fmt.Sprintf("D%d", i+1)
			f.SetCellValue("Sheet1", cellDiff, finalDiff)

			found = true
			break
		}
	}

	if !found {
		http.Error(w, "ID tidak ditemukan di dalam Excel", http.StatusNotFound)
		return
	}

	if err := f.SaveAs("database.xlsx"); err != nil {
		http.Error(w, "Gagal menyimpan perubahan ke file", http.StatusInternalServerError)
		return
	}

	// Kirim respon JSON
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status":"success", "diff": %.2f}`, finalDiff)
}

func searchExcel(query string) ([]Data, error) {
	f, err := excelize.OpenFile("database.xlsx")
	if err != nil {
		return nil, fmt.Errorf("make sure database.xlsx exists with 'Sheet1': %v", err)
	}
	defer f.Close()

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return nil, err
	}

	var results []Data
	if len(rows) == 0 {
		return results, nil
	}

	for i := 1; i < len(rows); i++ {
		row := rows[i]
		if len(row) < 2 {
			continue
		}

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
			tallyIn := ""
			diff := ""
			stock := ""

			if len(row) >= 2 {
				parts := strings.Split(row[1], "=")
				if len(parts) > 1 {
					val := strings.TrimSpace(parts[1])
					stock = strings.TrimRight(val, " mc")
					stock = strings.TrimSpace(stock)
				}
			}

			if len(row) >= 3 {
				tallyIn = row[2]
			}
			if len(row) >= 4 {
				diff = row[3]
			}

			results = append(results, Data{
				ID:      row[0],
				Name:    row[1],
				Stock:   stock,
				TallyIn: tallyIn,
				Diff:    diff,
			})
		}
	}

	return results, nil
}
