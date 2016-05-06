# CleanMyCar
This aplication detect, when is the perfect day to clean your car. And send email with summary.

## Install
For to build you need:

- golang

``sudo aptitude install golang``


- forecast

``go get github.com/mlbright/forecast/v2``

- gomail

``go get gopkg.in/gomail.v2``

## Configure
Open ``cleanMyCar.go`` file and edit this part to your information

```
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
```


## Finally build

``go build cleanMyCar.go``

and run:

``./cleanMyCar.go``

or without build execute

``go run clenMyCar.go``
