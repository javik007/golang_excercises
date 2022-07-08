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

type MonthTypes struct {
	Even mapset.Set[string]
	Odd  mapset.Set[string]
}

func main() {
	start_date := "2020-12-25"
	end_date := "2021-02-17"

	var result SearchStruct

	// Slice of months, data can be repeated, helps to know if month is complete or not
	var dates []string

	// This set contains all months to be analyzed for adding to search struct, even if no full months : YYYY-MM
	dateSet := mapset.NewSet[string]()

	// This set contains all days to be analyzed for adding to search struct : YYYY-MM
	daySet := mapset.NewSet[string]()

	// This set contains only full months which will be send to monthly search
	fullMonthSet := mapset.NewSet[string]()

	// Parsed dates for iteration
	from_Param, _ := time.Parse(ShortDayFormat, start_date)
	until_Param, _ := time.Parse(ShortDayFormat, end_date)

	// This cycle creates all days between the start param and end param
	for d := from_Param; !d.After(until_Param); d = d.AddDate(0, 0, 1) {
		day := d.Format(ShortDayFormat)
		slicedDay := strings.Split(day, "-")
		month := monthFormat(slicedDay)
		dates = append(dates, month)
		dateSet.Add(month)
		daySet.Add(day)
	}

	// This cycle adds a month to the fullMonth slice if day range is full
	for _, item := range dateSet.ToSlice() {
		// Evaluates month for adding to struct
		EvaluateMonths(dates, item, &result, fullMonthSet)
	}

	// This cycle adds to single day struct if YYYY-MM side is not on the full month struct
	for _, item := range daySet.ToSlice() {
		d := strings.Split(item, "-")
		month := monthFormat(d)
		if !fullMonthSet.Contains(month) {
			result.AddDay(item)
		}
	}

	fmt.Println(len(result.FullMonths))
	fmt.Println(result.FullMonths)

	fmt.Println(len(result.Days))
	From(result.Days).OrderBy(
		func(i interface{}) interface{} { return i },
	).Distinct().ToSlice(&result.Days)
	fmt.Println(result.Days)
}

// Evaluates month type for adding to struct according to if it is even, odd or february
func EvaluateMonths(dates []string, item string, result *SearchStruct, fullMonthSet mapset.Set[string]) {
	monthTypes := FillMonthStruct()
	m := strings.Split(item, "-")
	month := m[1]
	if month == "02" {
		count := CountDays(dates, item)
		if count == 28 || count == 29 {
			addDate(item, result, fullMonthSet)
		}
	} else {
		if monthTypes.Even.Contains(month) {
			EvaluateMonth(dates, 30, item, result, fullMonthSet)
		}
		if monthTypes.Odd.Contains(month) {
			EvaluateMonth(dates, 31, item, result, fullMonthSet)
		}
	}
}

// Evaluates single month for adding to struct whether it is even or odd
func EvaluateMonth(dates []string, monthDays int, item string, result *SearchStruct, fullMonthSet mapset.Set[string]) {
	count := CountDays(dates, item)
	if count == monthDays {
		addDate(item, result, fullMonthSet)
	}
}

// Counts repeated dates with linq style
func CountDays(dates []string, item string) int {
	return From(dates).
		CountWith(
			func(i interface{}) bool { return i == item },
		)
}

// Adds date in format YYYY-MM to both struct and set
func addDate(item string, result *SearchStruct, fullMonthSet mapset.Set[string]) {
	result.AddMonth(item)
	fullMonthSet.Add(item)
}

func (result *SearchStruct) AddMonth(item string) {
	result.FullMonths = append(result.FullMonths, item)
}

func (result *SearchStruct) AddDay(item string) {
	result.Days = append(result.Days, item)
}

// Fills month types struct according to quantity of days they have
func FillMonthStruct() *MonthTypes {
	var monthTypes MonthTypes
	monthTypes.Even = mapset.NewSet[string]("04", "06", "09", "11")
	monthTypes.Odd = mapset.NewSet[string]("01", "03", "05", "07", "08", "10", "12")
	return &monthTypes
}

// Returns sliced day into YYYY-MM format
func monthFormat(sliced []string) string {
	result := fmt.Sprintf("%s-%s", sliced[0], sliced[1])
	return result
}
