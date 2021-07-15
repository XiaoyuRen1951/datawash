package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	
	"strconv"
	"os"
	"io"
	"bufio"
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

func cmain() {
	ReadMachineMemory()
	ReadGPUTemp()
	return
}

func main() {

	Podmp = make(map[string]TaskLog)
	Nodeio := make(map[string]NodeIO)

	var date []string
	/*var date []string = []string{"/raw-data/result-08-27-30","/raw-data/result-08-31-03","/raw-data/result-09-24-26","/raw-data/result-09-27-30", "/raw-data/result-13-16", "/raw-data/result-17-19",
		"/raw-data/result-20-22",
		"/raw-data/result-23-24",
		"/raw-data/result-4-4",
		"/raw-data/result-5-8",
		"/raw-data/result-9-12"}*/
	// var date []string = []string{
	// 	// "2020-11-19",
	// 	// "2020-11-20",
	// 	// "2020-11-21",
	// 	// "2020-11-22",
	// 	// "2020-11-23",
	// 	// "2020-11-24",
	// 	// "2020-11-25",
	// 	// "2020-11-26",
	// 	// "2020-11-27",
	// 	// "2020-11-28",
	// 	// "2020-11-29",
	// 	// "2020-11-30",
	// 	// "2020-12-01",
	// 	// "2020-12-02",
	// 	// "2020-12-03",
	// 	// "2020-12-04",
	// 	// "2020-12-05",
	// 	// "2020-12-06",
	// 	// "2020-12-07",
	// 	// "2020-12-08",
	// 	// "2020-12-09",
	// 	// "2020-12-10",
	// 	// "2020-12-11",
	// 	// "2020-12-12",
	// 	// "2020-12-13",
	// 	// "2020-12-14",
	// 	// "2020-12-15",
	// 	// "2020-12-16",
	// 	// "2020-12-17",
	// 	// "2020-12-18",
	// 	// "2020-12-19",
	// 	// "2020-12-20",
	// 	// "2020-12-21",
	// 	// "2020-12-22",
	// 	// "2020-12-23",
	// 	// "2020-12-24",

	// 	// "2020-12-25",
	// 	// "2020-12-26",
	// 	// "2020-12-27",
	// 	// "2020-12-28",
	// 	// "2020-12-29",
	// 	// "2020-12-30",
	// 	// "2020-12-31",
	// 	// "2021-01-01",
	// 	// "2021-01-02",
	// 	// "2021-01-03",
	// 	// "2021-01-04",
	// 	// "2021-01-05",
	// 	// "2021-01-06",
	// 	// "2021-01-07",
	// 	// "2021-01-08",
	// 	// "2021-01-09",
	// 	// "2021-01-10",
	// 	// "2021-01-11",
	// 	// "2021-01-12",
	// 	// "2021-01-13",
	// 	// "2021-01-14",
	// 	// "2021-01-15",
	// 	// "2021-01-16",
	// 	// "2021-01-17",
	// 	// "2021-01-18",
	// 	// "2021-01-19",
	// 	// "2021-01-20",
	// 	// "2021-01-21",
	// 	// "2021-01-22",
	// 	// "2021-01-23",

	// 	// "2021-01-24",
	// 	// "2021-01-25",
	// 	// "2021-01-26",
	// 	// "2021-01-27",
	// 	// "2021-01-28",
	// 	// "2021-01-29",
	// 	// "2021-01-30",
	// 	// "2021-01-31",
	// 	// "2021-02-01",
	// 	// "2021-02-02",
	// 	// "2021-02-03",
	// 	// "2021-02-04",
	// 	// "2021-02-05",
	// 	// "2021-02-06",
	// 	// "2021-02-07",
	// 	// "2021-02-08",
	// 	// "2021-02-09",
	// 	// "2021-02-10",
	// 	"2021-02-18",
	// 	"2021-02-19",
	// 	"2021-02-20",
	// 	"2021-02-21",
	// 	"2021-02-22",
	// 	"2021-02-23",
	// 	"2021-02-24",
	// 	"2021-02-25",
	// 	"2021-02-26",
	// 	"2021-02-27",

	// 	"2021-02-28",
	// 	"2021-03-01",
	// 	"2021-03-02",
	// 	"2021-03-03",
	// 	"2021-03-04",
	// 	"2021-03-05",
	// 	"2021-03-06",
	// 	"2021-03-07",
	// 	"2021-03-08",
	// 	"2021-03-09",
	// 	"2021-03-10",
	// 	"2021-03-11",
	// 	"2021-03-12",

	// 	"2021-03-13",
	// 	"2021-03-14",

	// 	"2021-03-15",
	// 	"2021-03-16",

	// 	"2021-03-17",
	// 	"2021-03-18",
	// 	"2021-03-19",
	// 	"2021-03-20",
	// 	"2021-03-21",
	// 	"2021-03-22",
	// 	"2021-03-23",
	// 	"2021-03-24",
	// 	"2021-03-25",

	// 	"2021-03-26",
	// 	"2021-03-27",
	// 	"2021-03-28",
	// 	"2021-03-29",
	// 	"2021-03-30",
	// 	"2021-03-31",

	// 	"2021-04-01",

	// 	"2021-04-02",
	// 	"2021-04-03",

	// 	"2021-04-04",
	// 	"2021-04-05",
	// 	"2021-04-06",
	// 	"2021-04-07",
	// 	"2021-04-08",
	// 	"2021-04-09",
	// 	"2021-04-10",
	// 	"2021-04-11",
	// 	"2021-04-12",
	// 	"2021-04-13",
	// 	"2021-04-14",
	// 	"2021-04-15",
	// 	"2021-04-16",
	// 	"2021-04-17",
	// 	"2021-04-18",
	// 	"2021-04-19",
	// }
	File, err := os.Open("./date.log")
	defer File.Close()

	if err != nil {
		fmt.Println(err)
		return
	}

	rd := bufio.NewReader(File)

	for {
		line, err := rd.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}
		
		date = append(date,line[:len(line)-1])
		
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
