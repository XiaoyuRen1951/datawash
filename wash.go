package main

import (
	"encoding/json"
	"fmt"
	"reflect"

	"strconv"
)

const INT_MAX = int64(^uint64(0) >> 1)

var Podmp map[string]TaskLog

func Err_Handle(err error) bool{
	if err != nil {
		fmt.Println(err)
		return true
	}
	return false
}

func MIN(x, y interface{}) interface{} {
	if reflect.TypeOf(x).Name() == "int64" {
		if reflect.ValueOf(x).Int() < reflect.ValueOf(y).Int() {
			return x
		}
		return y
	} else if reflect.TypeOf(x).Name() == "float64" {
		if reflect.ValueOf(x).Float() < reflect.ValueOf(y).Float() {
			return x
		}
		return y
	}

	return x
}

func MAX(x, y interface{}) interface{} {
	if reflect.TypeOf(x).Name() == "int64" {
		if reflect.ValueOf(x).Int() > reflect.ValueOf(y).Int() {
			return x
		}
		return y
	} else if reflect.TypeOf(x).Name() == "float64" {
		if reflect.ValueOf(x).Float() > reflect.ValueOf(y).Float() {
			return x
		}
		return y
	}
	return x
}

func Deal_CPU_Str_Compare(x,y string, flag bool) string {
	if y == "" {
		return x
	}
	if x == "" {
		return y
	}
	x_v,err := strconv.ParseInt(x[:len(x)-1],10,64)
	if Err_Handle(err) {
		return ""
	}
	y_v,err := strconv.ParseInt(y[:len(y)-1],10,64)
	if Err_Handle(err) {
		return ""
	}
	if flag {
		res := strconv.FormatInt(MAX(x_v,y_v).(int64),10)
		return res+"n"
	} else {
		res := strconv.FormatInt(MIN(x_v,y_v).(int64),10)
		return res+"n"
	}
}

var NodetoGPUtot = make(map[string](map[int64]int64))
var NodetoGPUuse = make(map[string](map[int64]int64))

var enc *json.Encoder

var timestampstart int64 = INT_MAX
var timestampend int64 = 0

func Decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}

func Cal_Average(src []int64) float64 {
	var sum int64 = 0
	for _,v := range src {
		sum = sum + v
	}

	return Decimal(float64(sum)/float64(len(src)))
}

func main() {

	Podmp = make(map[string]TaskLog)
	Nodeio := make(map[string]NodeIO)

	/*var date []string = []string{"/raw-data/result-08-27-30","/raw-data/result-08-31-03","/raw-data/result-09-24-26","/raw-data/result-09-27-30", "/raw-data/result-13-16", "/raw-data/result-17-19",
		"/raw-data/result-20-22",
		"/raw-data/result-23-24",
		"/raw-data/result-4-4",
		"/raw-data/result-5-8",
		"/raw-data/result-9-12"}*/
	var date []string = []string{
	//	"2020-11-18",
	// 	"2020-11-19",
	// 	"2020-11-20",
	// 	"2020-11-21",
	// 	"2020-11-22",
	// 	"2020-11-23",
	// 	"2020-11-24",
	// 	"2020-11-25",
	// 	"2020-11-26",
	// 	"2020-11-27",
	// 	"2020-11-28",
	// 	"2020-11-29",
	// 	"2020-11-30",
	// 	"2020-12-01",
	// 	"2020-12-02",
	// 	"2020-12-03",
	// 	"2020-12-04",
	// 	"2020-12-05",
	// 	"2020-12-06",
	// 	"2020-12-07",
	// 	"2020-12-08",
	// 	"2020-12-09",
	// 	"2020-12-10",
	// 	"2020-12-11",
	// 	"2020-12-12",
	// 	"2020-12-13",
	// 	"2020-12-14",
	// 	"2020-12-15",
	// 	"2020-12-16",
	// 	"2020-12-17",
	// 	"2020-12-18",
	// 	"2020-12-19",
	// 	"2020-12-20",
	// 	"2020-12-21",
	// 	"2020-12-22",
	// 	"2020-12-23",
	// 	"2020-12-24",
	// 	"2021-01-24",
	// 	"2021-01-25",
	//	"2021-02-25",
	//	"2021-02-26",
	//	"2021-02-27",
		// "2021-03-01",
		// "2021-03-02",
		// "2021-03-03",
		// "2021-03-04",
		// "2021-03-05",
		// "2021-03-06",
		// "2021-03-07",
		// "2021-03-08",
		// "2021-03-09",
		// "2021-03-10",
		// "2021-03-11",
		// "2021-03-12",

		// "2021-03-15",
		// "2021-03-16",
		// "2021-03-26",
		// "2021-03-27",
		// "2021-03-28",
		"2021-03-29",
	}

	var bar Bar
    bar.NewOption(0, int64(len(date)))

	for i,d := range date {
		Deal_Oneday_data(d,Nodeio)
		bar.Play(int64(i+1))
	}
	fmt.Println()
	// Deal_Oneday_data("2021-01-24",Nodeio)
	// Deal_Oneday_data("2021-01-25",Nodeio)
	// Deal_Oneday_data("2021-02-25",Nodeio)

	//fmt.Println("Read Finish")
	OuttoFile(Nodeio)

	//PrintGPUCPUUtiRange()



	return
}
