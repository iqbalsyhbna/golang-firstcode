// internal/background/job.go
package background

import (
	"bytes"
	"encoding/json"
	"fmt"
	"golang-firstcode/internal/helpers"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/robfig/cron/v3"
)

type PresenceData struct {
	CheckType string  `json:"check_type"`
	Code      string  `json:"code"`
	CheckTime string  `json:"check_time"`
	CheckDate string  `json:"check_date"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func StartBackgroundJob() *cron.Cron {
	c := cron.New(cron.WithSeconds())
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	days := []string{"1", "2", "3", "4", "5"}

	for _, day := range days {
		exactHourCheckIn := 7
		exactMinuteCheckIn := r.Intn(16) + 25
		randomSecondCheckIn := r.Intn(60)

		today := time.Now()
		if int(today.Weekday()) != 0 && int(today.Weekday()) <= 5 && helpers.IsHoliday(today) {
			log.Printf("Skipping check-in job for day %s because it's a holiday", day)
			continue
		}

		checkInSchedule := fmt.Sprintf("CRON_TZ=Asia/Jakarta %d %d %d * * %s", randomSecondCheckIn, exactMinuteCheckIn, exactHourCheckIn, day)
		_, err := c.AddFunc(checkInSchedule, func() {
			if err := postToAPI("Check-in", "851269"); err != nil {
				log.Printf("Failed to perform check-in: %v", err)
			}
		})

		if err != nil {
			log.Fatalf("Failed to schedule check-in job for day %s: %v", day, err)
		}
		log.Printf("Scheduled check-in job for day %s at %02d:%02d:%02d", day, exactHourCheckIn, exactMinuteCheckIn, randomSecondCheckIn)
	}

	for _, day := range days {
		randomSecondCheckOut := r.Intn(60)
		randomMinuteCheckOut := r.Intn(60)
		randomHourCheckOut := r.Intn(2) + 17

		today := time.Now()
		if int(today.Weekday()) != 0 && int(today.Weekday()) <= 5 && helpers.IsHoliday(today) {
			log.Printf("Skipping check-out job for day %s because it's a holiday", day)
			continue
		}

		checkOutSchedule := fmt.Sprintf("CRON_TZ=Asia/Jakarta %d %d %d * * %s", randomSecondCheckOut, randomMinuteCheckOut, randomHourCheckOut, day)
		_, err := c.AddFunc(checkOutSchedule, func() {
			if err := postToAPI("Check-out", "851269"); err != nil {
				log.Printf("Failed to perform check-out: %v", err)
			}
		})
		if err != nil {
			log.Fatalf("Failed to schedule check-out job for day %s: %v", day, err)
		}
		log.Printf("Scheduled check-out job for day %s at %02d:%02d:%02d", day, randomHourCheckOut, randomMinuteCheckOut, randomSecondCheckOut)
	}

	c.Start()
	log.Println("Cron jobs started with random times for check-in and check-out.")
	return c
}

func postToAPI(checkType, code string) error {
	now := time.Now()
	checkTime := now.Format("2006-01-02 15:04:05")
	checkDate := now.Format("2006-01-02")

	data := PresenceData{
		CheckType: checkType,
		Code:      code,
		CheckTime: checkTime,
		CheckDate: checkDate,
		Latitude:  -6.184814158122739,
		Longitude: 106.9310312061378,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error encoding JSON: %w", err)
	}

	log.Printf("Sending data: %s", string(jsonData))

	req, err := http.NewRequest("POST", getAPIURL(), bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+getToken())
	req.Header.Set("DeviceId", getDeviceID())
	req.Header.Set("PlatformID", getPlatformID())
	req.Header.Set("site", "mobile-site")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}

	log.Printf("Response body: %s", string(body))

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status: %s, response: %s", resp.Status, string(body))
	}

	log.Printf("Request successful: %s", checkType)
	return nil
}

func getToken() string {
	if token := os.Getenv("AUTH_TOKEN"); token != "" {
		return token
	}
	log.Println("Warning: Using default token. Set AUTH_TOKEN environment variable.")
	return "your-default-token-here"
}

func getDeviceID() string {
	if deviceID := os.Getenv("DEVICE_ID"); deviceID != "" {
		return deviceID
	}
	log.Println("Warning: Using default device ID. Set DEVICE_ID environment variable.")
	return "default-device-id"
}

func getPlatformID() string {
	if platformID := os.Getenv("PLATFORM_ID"); platformID != "" {
		return platformID
	}
	log.Println("Warning: Using default platform ID. Set PLATFORM_ID environment variable.")
	return "default-platform-id"
}

func getAPIURL() string {
	if url := os.Getenv("API_URL"); url != "" {
		return url
	}
	log.Println("Warning: Using default API URL. Set API_URL environment variable.")
	return "https://api-gateway.triatra.co.id/api2/time-cards/mobile-presence"
}
