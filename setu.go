package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/corpix/uarand"
)

const (
	cDateFormat      = "02-01-2006"
	cEmailDateFormat = "Jan 02, 2006"
	cSetuURL         = "https://cdn-api.co-vin.in/api/v2/appointment/sessions/public/calendarByDistrict"
)

type apiResp struct {
	Centers []struct {
		CenterID int `json:"center_id"`
		Sessions []struct {
			AvailableCap int `json:"available_capacity"`
			MinAgeLimit  int `json:"min_age_limit"`
		} `json:"sessions"`
	} `json:"centers"`
}

func getSlotsForDays(districtID int, centers []int, days []time.Time) (int, error) {
	totalSlots := 0
	for _, d := range days {
		slots, err := getSlots(districtID, centers, d)
		if err != nil {
			date := d.Format(cEmailDateFormat)
			return 0, fmt.Errorf("error getting slots for %v: %w", date, err)
		}

		totalSlots += slots
	}

	return totalSlots, nil
}

func getSlots(districtID int, centers []int, t time.Time) (int, error) {
	date := t.Format(cDateFormat)
	url := fmt.Sprintf("%v?district_id=%v&date=%v", cSetuURL, districtID, date)

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", uarand.GetRandom())
	if err != nil {
		return 0, fmt.Errorf("error creating request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("error hitting API: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("error reading body: %w", err)
	}

	if resp.StatusCode != 200 {
		return 0, fmt.Errorf("code: %v, response: %v", resp.StatusCode, string(body))
	}

	var data apiResp
	if err := json.Unmarshal(body, &data); err != nil {
		return 0, fmt.Errorf("error unmarshalling: %w", err)
	}

	if len(data.Centers) == 0 {
		// This happens when dates are out of bound
		return 0, nil
	}

	for _, c := range data.Centers {
		if inSlice(c.CenterID, centers) {
			total := 0
			for _, s := range c.Sessions {
				if s.MinAgeLimit != 18 {
					continue
				}

				total += s.AvailableCap
			}

			return total, nil
		}
	}

	return 0, errors.New("centers not found")
}

func inSlice(n int, s []int) bool {
	// Handle special case
	if s[0] == 0 {
		return true
	}

	for _, e := range s {
		if n == e {
			return true
		}
	}

	return false
}
