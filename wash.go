package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
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

func DeepCopy(value interface{}) interface{} {
	if valueMap, ok := value.(map[string]interface{}); ok {
		newMap := make(map[string]interface{})
		for k, v := range valueMap {
			newMap[k] = DeepCopy(v)
		}

		return newMap
	} else if valueSlice, ok := value.([]interface{}); ok {
		newSlice := make([]interface{}, len(valueSlice))
		for k, v := range valueSlice {
			newSlice[k] = DeepCopy(v)
		}

		return newSlice
	}

	return value
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
//CPU Utilization
var mpcpu = make(map[string]int64)

var timestampstart int64 = INT_MAX
var timestampend int64 = 0

//不同ResourceType的利用率
func DiffResourceTGPUUti(filepath string) {
	tmp, err := os.Create(filepath)
	
	if err != nil {
		fmt.Println("tmp File creating error", err)
		return
	}
	resourcetymp := make(map[string]([]int))
	for _,v := range Podmp {
		if v.ResourceT == "debug" {
			continue
		}
		tmpr, ok := resourcetymp[v.ResourceT]
		if !ok {
			tmpr = make([]int, 0)
		}
		if v.GPU.NumGPU < 1 {
			continue
		}
		if len(v.GPU.GPUUtil) == 0 {
			tmpr = append(tmpr, 0)
			resourcetymp[v.ResourceT] = tmpr
			continue
		}
		var re float64 = 0
		for _, vv := range v.GPU.GPUUtil {
			re = re+Cal_Average(vv.History)
		}
		re = math.Floor(re / float64(len(v.GPU.GPUUtil)))
		tmpr = append(tmpr, int(re))
		resourcetymp[v.ResourceT] = tmpr
	}
	
	for k,v := range resourcetymp {
		tmp.WriteString(k)
		for _,vv := range v {
			tmp.WriteString(fmt.Sprintf(" %d", vv))
		}
		tmp.WriteString("\n")
	}
	tmp.Close()
}

//单卡任务利用率
func SCGPUUtui(filepath string) {
	tmp, err := os.Create(filepath)

	if err != nil {
		fmt.Println("tmp File creating error", err)
		return
	}
	scutil := make([]int,102)
	scutig := make([]int,102)
	for _,v := range Podmp {
		if v.ResourceT == "debug" {
			continue
		}
		if v.GPU.NumGPU != 1 {
			continue
		}
		//if v.Endtime - v.Starttime >=7200 {
		//	continue
		//}
		if len(v.GPU.GPUUtil) == 0 {
			if v.Endtime - v.Starttime >=10800 {
				scutig[0]++
			} else {
				scutil[0]++
			}
			continue
		}
		var res float64 = 0
		for _,vv := range v.GPU.GPUUtil {
	
			res += Cal_Average(vv.History)
	
		}
		res = math.Floor(res / float64(len(v.GPU.GPUUtil)))
		if v.Endtime - v.Starttime >=10800 {
			scutig[int(res)]++
		} else {
			scutil[int(res)]++
		}
	}
	sum := 0
	tmp.WriteString("scutil:")
	for i:=0;i<=100;i++ {
		tmp.WriteString(fmt.Sprintf("%d ", scutil[i]))
		sum += scutil[i]
	}
	tmp.WriteString("\n")
	fmt.Println(sum)
	
	sum = 0
	tmp.WriteString("scutig:")
	for i:=0;i<=100;i++ {
		tmp.WriteString(fmt.Sprintf("%d ", scutig[i]))
		sum += scutig[i]
	}
	tmp.WriteString("\n")
	fmt.Println(sum)

	tmp.Close()
}

//计算多卡最大最小利用率
func MCGPUUti(filepath string) {
	tmp, err := os.Create(filepath)

	if err != nil {
		fmt.Println("tmp File creating error", err)
		return
	}

	mcutil := make(map[int]int)
	mcutil[-1]=0
	mcutig := make(map[int]int)
	mcutig[-1]=0
	for _,v := range Podmp {
		if v.ResourceT == "debug" {
			continue
		}
		if v.GPU.NumGPU <= 1 {
			continue
		}
		//if v.Endtime - v.Starttime >=10800 {
		//	continue
		//}
		if len(v.GPU.GPUUtil) == 0 {
			if v.Endtime - v.Starttime >=10800 {
				mcutig[-1]++
			} else {
				mcutil[-1]++
			}
			continue
		}
		var mn int64 = 102
		var mx int64 = 0
		for _,vv := range v.GPU.GPUUtil {
			mn = MIN(int64(Cal_Average(vv.History)),mn).(int64)
			mx = MAX(int64(Cal_Average(vv.History)),mx).(int64)
		}
		val :=-1
		if mx != 0 {
			val = int(math.Floor(float64(mn)*100/float64(mx)))
		}
		if v.Endtime - v.Starttime >=10800 {
			if _,ok:=mcutig[val];ok {
				mcutig[val]++
			} else {
				mcutig[val]=1
			}
		} else {
			if _,ok:=mcutil[val];ok {
				mcutil[val]++
			} else {
				mcutil[val]=1
			}
		}
	
	}
	sum := 0
	tmp.WriteString("mcutil:")
	for i:=-1;i<=100;i++ {
		if _,ok:=mcutil[i];ok {
			tmp.WriteString(fmt.Sprintf("%d ", mcutil[i]))
			sum += mcutil[i]
		} else {
			tmp.WriteString(fmt.Sprintf("0 "))
		}
	}
	fmt.Println(sum)
	tmp.WriteString("\n")
	
	sum = 0
	tmp.WriteString("mcutig:")
	for i:=-1;i<=100;i++ {
		if _,ok:=mcutil[i];ok {
			tmp.WriteString(fmt.Sprintf("%d ", mcutig[i]))
			sum += mcutig[i]
		} else {
			tmp.WriteString(fmt.Sprintf("0 "))
		}
	}
	fmt.Println(sum)
	tmp.WriteString("\n")

	tmp.Close()
}

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
		"2021-03-01",
		"2021-03-02",
		"2021-03-03",
		"2021-03-04",
		"2021-03-05",
		"2021-03-06",
		"2021-03-07",
		"2021-03-08",
		"2021-03-09",
		"2021-03-10",
		"2021-03-11",
		"2021-03-12",
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

	AverageVal("./result/renxy-task.log")

	// tmp, err := os.Create("./result/tmp.log")
	// defer tmp.Close()
	// if err != nil {
	// 	fmt.Println("tmp File creating error", err)
	// 	return
	// }

	
	// sum := 0
	// for i:=-100;i<=100;i++ {
	// 	if _,ok:=mcutil[i];ok {
	// 		tmp.WriteString(fmt.Sprintf("%d ", mcutil[i]))
	// 		sum += mcutil[i]
	// 	} else {
	// 		tmp.WriteString(fmt.Sprintf("0 "))
	// 	}
	// }
	// fmt.Println(sum)
	// tmp.WriteString("\n")

	// sum = 0
	// for i:=-100;i<=100;i++ {
	// 	if _,ok:=mcutig[i];ok {
	// 		tmp.WriteString(fmt.Sprintf("%d ", mcutig[i]))
	// 		sum += mcutig[i]
	// 	} else {
	// 		tmp.WriteString(fmt.Sprintf("0 "))
	// 	}
	// }
	// fmt.Println(sum)
	// tmp.WriteString("\n")

	return
}
