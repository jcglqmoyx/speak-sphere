package util

import (
	"strconv"
	"time"
)

func DateToString(time time.Time) string {
	return time.Format("20060102")
}

func DatetimeToString(time time.Time) string {
	return time.Format("20060102150405")
}

func DateToInt(date string) int {
	today := DateToString(time.Now())
	res, _ := strconv.Atoi(today)
	return res
}

func AddDaysToIntDate(currentDate int, n int) int {
	currentDateObj, _ := time.Parse("20060102", strconv.Itoa(currentDate))
	newDateObj := currentDateObj.AddDate(0, 0, n).Format("20060102")
	newDate, _ := strconv.Atoi(newDateObj)
	return newDate
}
