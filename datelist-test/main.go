package main

import (
	"fmt"
	"strings"
	"time"

	. "github.com/ahmetb/go-linq/v3"
	mapset "github.com/deckarep/golang-set/v2"
)

const ShortDayFormat = "2006-01-02"

type SearchStruct struct {
	FullMonths []string
	Days       []string
}

func main() {
	start_date := "2021-01-31"
	end_date := "2021-03-01"

	var result SearchStruct

	// Slice of months, data can be repeated, helps to know if month is complete or not
	var dates []string

	//dateSet := set.New(set.ThreadSafe)

	// This set contains all months to be analyzed for adding to search struct, even if no full months : YYYY-MM
	dateSet := mapset.NewSet[string]()

	// This set contains all days to be analyzed for adding to search struct : YYYY-MM
	daySet := mapset.NewSet[string]()

	// This set contains only full months which will be send to monthly search
	fullMonthSet := mapset.NewSet[string]()

	// This sets contain months with either 31 or 30 days in string-int format
	s31 := mapset.NewSet[string]("01", "03", "05", "07", "08", "10", "12")
	s30 := mapset.NewSet[string]("04", "06", "09", "11")

	// Parsed dates for iteration
	from_Param, _ := time.Parse(ShortDayFormat, start_date)
	until_Param, _ := time.Parse(ShortDayFormat, end_date)

	// This cycle creates all days between the start param and end param
	for d := from_Param; !d.After(until_Param); d = d.AddDate(0, 0, 1) {
		day := d.Format(ShortDayFormat)
		slicedDay := strings.Split(day, "-")
		//month := fmt.Sprintf("%s-%s", slicedDay[0], slicedDay[1])
		month := monthFormat(slicedDay)
		dates = append(dates, month)
		dateSet.Add(month)
		daySet.Add(day)
	}

	// This cycle adds a month to the fullMonth slice if day range is full
	for _, item := range dateSet.ToSlice() {
		m := strings.Split(item, "-")
		month := m[1]
		if s31.Contains(month) {
			count := From(dates).
				CountWith(
					func(i interface{}) bool { return i == item },
				)
			if count == 31 {
				result.AddMonth(item)
				fullMonthSet.Add(item)
			}
		}
		if s30.Contains(month) {
			count := From(dates).
				CountWith(
					func(i interface{}) bool { return i == item },
				)
			if count == 30 {
				result.AddMonth(item)
				fullMonthSet.Add(item)
			}
		}
		if month == "02" {
			count := From(dates).
				CountWith(
					func(i interface{}) bool { return i == item },
				)
			if count == 28 || count == 29 {
				result.AddMonth(item)
				fullMonthSet.Add(item)
			}
		}
	}

	// This cycle adds to single day struct if YYYY-MM side is not on the full month struct
	for _, item := range daySet.ToSlice() {
		d := strings.Split(item, "-")
		//month := fmt.Sprintf("%s-%s", d[0], d[1])
		month := monthFormat(d)
		if !fullMonthSet.Contains(month) {
			result.AddDay(item)
		}
	}

	fmt.Println(len(result.FullMonths))
	fmt.Println(result.FullMonths)

	fmt.Println(len(result.Days))
	fmt.Println(result.Days)
}

func EvaluateMonth(dates []string, month string, result *SearchStruct, months mapset.Set[string]) mapset.Set[string] {
	fullMonthSet := mapset.NewSet[string]()
	return fullMonthSet
}

func (result *SearchStruct) AddMonth(item string) {
	result.FullMonths = append(result.FullMonths, item)
}

func (result *SearchStruct) AddDay(item string) {
	result.Days = append(result.Days, item)
}

// Returns sliced day into YYYY-MM format
func monthFormat(sliced []string) string {
	result := fmt.Sprintf("%s-%s", sliced[0], sliced[1])
	return result
}
