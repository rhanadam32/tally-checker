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
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/search", searchHandler)
	http.HandleFunc("/update", updateHandler)

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
			tallyIn := ""
			diff := ""
			stock := ""

			// Logic to extract stock value from Column B (Name/Target Qty)
			// Example: "kt 1/2 = 2 mc" -> result "2"
			if len(row) >= 2 {
				parts := strings.Split(row[1], "=")
				if len(parts) > 1 {
					val := strings.TrimSpace(parts[1])
					// Remove non-numeric characters except decimal point
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

// updateHandler menangani pengiriman data dari form web untuk disimpan ke file Excel
func updateHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Hanya izinkan metode POST
	if r.Method != http.MethodPost {
		http.Error(w, "Metode tidak diizinkan", http.StatusMethodNotAllowed)
		return
	}

	// 2. Ambil data dari form
	id := r.FormValue("id")
	tallyInVal := r.FormValue("tallyin")

	if id == "" {
		http.Error(w, "ID harus diisi", http.StatusBadRequest)
		return
	}

	// 3. Buka file database.xlsx
	f, err := excelize.OpenFile("database.xlsx")
	if err != nil {
		http.Error(w, "Gagal membuka file Excel: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()

	// 4. Baca semua baris untuk mencari ID yang sesuai
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		http.Error(w, "Gagal membaca baris Excel", http.StatusInternalServerError)
		return
	}

	found := false
	for i, row := range rows {
		if i == 0 {
			continue // Lewati baris header (baris 1)
		}

		// Jika ID di kolom A cocok dengan ID yang dikirim dari form
		if len(row) > 0 && row[0] == id {
			// --- SIMPAN DATA TALLY-IN ---
			// Simpan ke Kolom C (Indeks 3)
			cellTally := fmt.Sprintf("C%d", i+1)
			f.SetCellValue("Sheet1", cellTally, tallyInVal)

			// --- HITUNG SELISIH OTOMATIS ---
			var targetQty, actualQty float64
			// Ambil Target Qty dari Kolom B (Indeks 2)
			if len(row) >= 2 {
				fmt.Sscanf(row[1], "%f", &targetQty)
			}
			// Ambil nilai input Tally-In
			fmt.Sscanf(tallyInVal, "%f", &actualQty)

			// Rumus: Tally-In - Target Qty
			diff := actualQty - targetQty

			// Simpan hasil selisih ke Kolom D (Indeks 4)
			cellDiff := fmt.Sprintf("D%d", i+1)
			f.SetCellValue("Sheet1", cellDiff, diff)

			found = true
			break
		}
	}

	if !found {
		http.Error(w, "ID tidak ditemukan di dalam Excel", http.StatusNotFound)
		return
	}

	// 5. Simpan perubahan ke file secara permanen
	if err := f.SaveAs("database.xlsx"); err != nil {
		http.Error(w, "Gagal menyimpan perubahan ke file", http.StatusInternalServerError)
		return
	}

	// 6. Redirect kembali ke halaman search untuk menampilkan data terbaru
	http.Redirect(w, r, "/search?q="+id, http.StatusSeeOther)
}
