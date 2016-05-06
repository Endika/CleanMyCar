package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	forecast "github.com/mlbright/forecast/v2"
	"gopkg.in/gomail.v2"
)

const KEY string = "forecast API KEY"
const LAT string = "LATITUDE"
const LONG string = "LONGITUDE"
const SMTP string = "smtp.server"
const PORT int = 587
const USER string = "USER"
const PASS string = "PASSWORD"
const FROM string = "FROM EMAIL"
const TO string = "TO EMAIL"
const SUBJECT string = "Car Clean!! The Summary"

func checkWeekDay(cDate time.Time) (result bool, days int) {
	result = true
	days = 0
	for i := 0; i < 7; i++ {
		if !checkHoursDay(cDate) {
			result = false
			break
		} else {
			days++
		}
		cDate = cDate.Add(time.Hour * 24)
	}
	return
}

func checkHoursDay(cDate time.Time) (result bool) {
	result = true
	hoursList := [...]string{
		"T07:00:00", "T10:00:00", "T12:00:00", "T13:00:00",
		"T15:00:00", "T16:00:00", "T18:00:00", "T19:00:00"}
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
		if !result {
			break
		}
	}
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
	cDate := time.Now()
	// cDate = cDate.Add(-time.Hour * 24 * 7)
	summary := ""
	for i := 0; i < 7; i++ {
		result := "NO"
		r, d := checkWeekDay(cDate)
		if r {
			result = "OK"
		}
		fmt.Printf("%s => %s (%d)\n", cDate.Format("2006-01-02"), result, d)
		summary = summary + " "
		summary = summary + cDate.Format("2006-01-02") + " => "
		summary = summary + result + " (" + strconv.Itoa(d) + ")<br>"
		cDate = cDate.Add(time.Hour * 24)
	}
	sendMail(summary)
}
