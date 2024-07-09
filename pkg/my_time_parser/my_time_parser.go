package my_time_parser

import (
	"time"
)

const (
	MyRfc      = "2006-01-02T15:04:05"
	MyRfc2     = "02.01.2006T15:04"
	MyRfcMili  = "2006-01-02T15:04:05.000"
	MyDateOnly = "02.01.2006"
)

var (
	Msk, _ = time.LoadLocation("Europe/Moscow")
)

func Parse(timeVal string) (time.Time, error) {
	// fmt.Println(time.RFC3339)
	return time.Parse(MyRfc, timeVal)
}

func ParseInLocation(timeVal string, loc *time.Location) (time.Time, error) {
	return time.ParseInLocation(MyRfc, timeVal, loc)
}

func ParseInLocation_V2(timeVal string, loc *time.Location) (time.Time, error) {
	return time.ParseInLocation(MyRfc2, timeVal, loc)
}

func ParseInLocation_V3(timeVal string, loc *time.Location) (time.Time, error) {
	return time.ParseInLocation(MyDateOnly, timeVal, loc)
}

func ParseInLocation_V4(timeVal string, loc *time.Location) (time.Time, error) {
	return time.ParseInLocation(MyRfcMili, timeVal, loc)
}
