package helpers

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// Struktur untuk menyimpan respon dari API
type Holiday struct {
	Date string `json:"tanggal"`
	Name string `json:"keterangan"`
}

type CustomHoliday struct {
	Date        string `json:"date"`
	Description string `json:"description"`
}

var customHolidays = []CustomHoliday{
	{
		Date:        "2024-10-30",
		Description: "Cuti CPNS",
	},
}

// Fungsi untuk mendapatkan daftar tanggal merah dari API
func fetchHolidays() ([]Holiday, error) {
	// Mengabaikan verifikasi sertifikat SSL
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := httpClient.Get("https://dayoffapi.vercel.app/api")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch holidays: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response status: %v", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var holidays []Holiday
	if err := json.Unmarshal(body, &holidays); err != nil {
		return nil, fmt.Errorf("failed to unmarshal holidays: %v", err)
	}

	return holidays, nil
}

func IsHoliday(date time.Time) bool {
	holidays, err := fetchHolidays()

	if err != nil {
		log.Printf("Error fetching holidays: %v", err)
		return false
	}

	dateStr := date.Format("2006-01-02")

	for _, holiday := range customHolidays {
		if dateStr == holiday.Date {
			log.Printf("%s is a custom holiday: %s", dateStr, holiday.Description)
			return true
		}
	}

	for _, holiday := range holidays {
		if dateStr == holiday.Date {
			return true
		}
	}
	return false
}
