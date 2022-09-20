package service

import (
	"github.com/uniplaces/carbon"
	"strconv"
	"strings"
	"time"
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

func GetCarbonDate(date string) (*carbon.Carbon, error) {
	splits := strings.Split(date, "-")
	year, err := strconv.Atoi(splits[0])
	if err != nil {
		myError := MyError{Message: "yearを数値に変換できませんでした。"}
		return &carbon.Carbon{}, myError
	}
	month, err := strconv.Atoi(splits[1])
	if err != nil {
		myError := MyError{Message: "monthを数値に変換できませんでした。"}
		return &carbon.Carbon{}, myError
	}
	day, err := strconv.Atoi(splits[2])
	if err != nil {
		myError := MyError{Message: "dateを数値に変換できませんでした。"}
		return &carbon.Carbon{}, myError
	}
	carbonDate, err := carbon.CreateFromDate(year, time.Month(month), day, "Asia/Tokyo")
	if err != nil {
		myError := MyError{Message: "carbonに変換できませんでした。"}
		return &carbon.Carbon{}, myError
	}
	return carbonDate, nil
}

func CarbonToString(date *carbon.Carbon) string {
	year := date.Year()
	month := strconv.Itoa(int(date.Month()))
	day := date.Day()
	return strconv.Itoa(year) + "-" + month + "-" + strconv.Itoa(day)
}
