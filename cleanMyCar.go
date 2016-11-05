package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	forecast "github.com/mlbright/forecast/v2"
	"gopkg.in/gomail.v2"
)

const USER string = "USER"
const PASS string = "PASSWORD"
const FROM string = "FROM EMAIL"
const TO string = "TO EMAIL"
const PORT int = 587
const USER string = "USER"
const PASS string = "PASSWORD"
const FROM string = "FROM EMAIL"
const TO string = "TO EMAIL"
const SUBJECT string = "Car Clean!! The Summary"

type day_info struct {
	result bool
	moon   string
}

func checkWeekDay(cDate time.Time, cache_date map[string]day_info) (result bool, days int, moon string) {
	result = true
	days = 0
	moon = ""
	for i := 0; i < 7; i++ {
		moon_tmp := ""
		result, moon_tmp = checkHoursDay(cDate, cache_date)
		if moon == "" {
			moon = moon_tmp
		}
		if !result {
			result = false
			break
		} else {
			days++
		}
		cDate = cDate.Add(time.Hour * 24)
	}
	return
}

func checkHoursDay(cDate time.Time, cache_date map[string]day_info) (result bool, moon string) {
	result = true
	moon = ""
	if val, ok := cache_date[cDate.Format("2006-01-02")]; ok {
		result = val.result
		moon = val.moon
		return
	}
	hoursList := [...]string{
		"T07:00:00", "T10:00:00", "T12:00:00", "T13:00:00",
		"T15:00:00", "T16:00:00", "T18:00:00", "T19:00:00"}
	moonList := map[float64]string{
		0: "New Moon", 0.25: "First quarter Moon", 0.5: "Full Moon",
		0.75: "Last quarter Moon"}
	moonRangeList := [...]float64{0, 0.25, 0.5, 0.75}
	dangerList := [...]string{"rain", "snow", "sleet", "hail", "thunderstorm"}
	for _, value := range hoursList {
		checkDate := cDate.Format("2006-01-02") + value
		f, err := forecast.Get(KEY, LAT, LONG, checkDate, forecast.AUTO)
		if err != nil {
			log.Fatal(err)
		}
		for _, v := range dangerList {
			if f.Currently.Icon == v {
				result = false
				break
			}
		}
		if moon == "" {
			for _, mRL := range moonRangeList {
				MoonPhase := f.Daily.Data[0].MoonPhase
				if MoonPhase <= mRL {
					if val, ok := moonList[mRL]; ok {
						moon = strconv.FormatFloat(MoonPhase, 'f', 6, 64)
						moon = moon + " to " + val + " "
						moon = moon + strconv.FormatFloat(mRL, 'f', 6, 64)
						break
					}
				}
			}
		}

		if !result {
			break
		}
	}
	cache_date[cDate.Format("2006-01-02")] = day_info{result: result, moon: moon}
	return
}

func sendMail(text string) {
	m := gomail.NewMessage()
	m.SetHeader("From", FROM)
	m.SetHeader("To", TO)
	m.SetHeader("Subject", SUBJECT)
	m.SetBody("text/html", text)

	d := gomail.NewDialer(SMTP, PORT, USER, PASS)

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}

func main() {
	cache_date := make(map[string]day_info)
	cDate := time.Now()
	// cDate = cDate.Add(-time.Hour * 24 * 7)
	summary := ""
	for i := 0; i < 7; i++ {
		result := "NO"
		r, d, m := checkWeekDay(cDate, cache_date)
		if r {
			result = "OK"
		}
		fmt.Printf(
			"%s => %s Sunny days: %d Moon: %s\n",
			cDate.Format("2006-01-02"), result, d, m)
		summary = summary + " "
		summary = summary + cDate.Format("2006-01-02") + " => "
		summary = summary + result + " Sunny days: " + strconv.Itoa(d)
		summary = summary + " Moon: " + m + "<br>"
		cDate = cDate.Add(time.Hour * 24)
	}
	sendMail(summary)
}
