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

type WorkDay struct {
	Name string
	Date time.Time
}

func getDayName(t time.Time) string {
	switch t.Weekday() {
	case time.Monday:
		return "MON"
	case time.Tuesday:
		return "TUE"
	case time.Wednesday:
		return "WED"
	case time.Thursday:
		return "THU"
	case time.Friday:
		return "FRI"
	default:
		return ""
	}
}

func getWorkDaysThisWeek() []WorkDay {
	now := time.Now()
	currentWeekDay := now.Weekday()

	monday := now.AddDate(0, 0, -int(currentWeekDay)+1)

	workDays := []WorkDay{}
	for i := 0; i < 5; i++ {
		date := monday.AddDate(0, 0, i)
		day := WorkDay{
			Name: getDayName(date),
			Date: date,
		}
		workDays = append(workDays, day)
	}

	return workDays
}

func StartBackgroundJob() *cron.Cron {
	c := cron.New(cron.WithSeconds())
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	workDays := getWorkDaysThisWeek()

	for _, workDay := range workDays {
		exactHourCheckIn := 7
		exactMinuteCheckIn := r.Intn(11) + 25
		randomSecondCheckIn := r.Intn(60)

		if helpers.IsHoliday(workDay.Date) {
			log.Printf("Skipping check-in job for %s (%s) because it's a holiday",
				workDay.Name,
				workDay.Date.Format("2006-01-02"))
			continue
		}

		day := workDay.Name
		date := workDay.Date

		checkInSchedule := fmt.Sprintf("CRON_TZ=Asia/Jakarta %d %d %d * * %s", randomSecondCheckIn, exactMinuteCheckIn, exactHourCheckIn, day)
		_, err := c.AddFunc(checkInSchedule, func() {
			if err := postToAPI("Check-in", "851269"); err != nil {
				log.Printf("Failed to perform check-in for %s (%s): %v",
					day, date.Format("2006-01-02"), err)
			}
		})

		if err != nil {
			log.Fatalf("Failed to schedule check-in job for %s (%s): %v",
				day, date.Format("2006-01-02"), err)
		}

		log.Printf("Scheduled check-in job for %s (%s) at %02d:%02d:%02d",
			day, date.Format("2006-01-02"),
			exactHourCheckIn, exactMinuteCheckIn, randomSecondCheckIn)
	}

	for _, workDay := range workDays {
		randomSecondCheckOut := r.Intn(60)
		randomMinuteCheckOut := r.Intn(30) + 2
		randomHourCheckOut := 17

		if helpers.IsHoliday(workDay.Date) {
			log.Printf("Skipping check-out job for %s (%s) because it's a holiday",
				workDay.Name,
				workDay.Date.Format("2006-01-02"))
			continue
		}

		day := workDay.Name
		date := workDay.Date

		checkOutSchedule := fmt.Sprintf("CRON_TZ=Asia/Jakarta %d %d %d * * %s",
			randomSecondCheckOut, randomMinuteCheckOut, randomHourCheckOut, day)

		_, err := c.AddFunc(checkOutSchedule, func() {
			if err := postToAPI("Check-out", "851269"); err != nil {
				log.Printf("Failed to perform check-out for %s (%s): %v",
					day, date.Format("2006-01-02"), err)
			}
		})

		if err != nil {
			log.Fatalf("Failed to schedule check-out job for %s (%s): %v",
				day, date.Format("2006-01-02"), err)
		}

		log.Printf("Scheduled check-out job for %s (%s) at %02d:%02d:%02d",
			day, date.Format("2006-01-02"),
			randomHourCheckOut, randomMinuteCheckOut, randomSecondCheckOut)
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
