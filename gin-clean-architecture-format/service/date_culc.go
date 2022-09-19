package service

import (
	"github.com/uniplaces/carbon"
	"strconv"
)

func AfterWeek(carbon carbon.Carbon, afterWeek int) *carbon.Carbon {
	after := carbon.AddWeeks(afterWeek)
	return after
}

func GetToday() (string, error) {
	todayCarbon, err := carbon.Today("Asia/Tokyo")
	if err != nil {
		return "", err
	}
	year := todayCarbon.Year()
	month := todayCarbon.Month()
	date := todayCarbon.Day()
	today := strconv.Itoa(year) + "-" + strconv.Itoa(int(month)) + "-" + strconv.Itoa(date)
	return today, nil
}
