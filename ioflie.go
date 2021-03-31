package main

import (
	"fmt"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"encoding/json"
	"bufio"
	"io"
)

func ReadFile(file string) *PrometheusInfo {
	File, err := os.Open(file)

	res := new(PrometheusInfo)
	if err != nil {
		fmt.Println("File reading error", err)
		return res
	}
	prdec := json.NewDecoder(File)

	err = prdec.Decode(&res)
	if err != nil {
		fmt.Println(err)
	}
	File.Close()
	return res
}

func ReadPodLifecycle(dir string) error {
	File, err := os.Open("./"+dir+"/PodLifecycle_log.log")
	defer File.Close()

	if err != nil {
		fmt.Println(err)
		return err
	}

	rd := bufio.NewReader(File)
	for {
		line, err := rd.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}
		if line[0] != '{' {
			continue
		}
		res := new(PodMetricsList)
		prdec := json.NewDecoder(strings.NewReader(line))
		err = prdec.Decode(&res)
		if err != nil {
			fmt.Println(err)
			return err
		}
		if res.Metadata.Namespace == "ingress-nginx" || res.Metadata.Namespace == "kube-system" || res.Metadata.Namespace == "default" || res.Metadata.Namespace == "lens-metrics" || len(res.Containers) < 1 {
			continue
		}

		tmptask, ok := Podmp[res.Metadata.Name]
		if !ok {
			continue
		}
		tmpcpu, _ := strconv.ParseFloat(res.Containers[0].Usage.CPU,64)
		tmptask.CPU.History = append(tmptask.CPU.History, tmpcpu)
		//if len(res.Containers[0].Usage.CPU) > 1 {
			//fmt.Println(res.Containers[0].Usage.CPU,"a")
			//tmptask.CPU.Max = Deal_CPU_Str_Compare(res.Containers[0].Usage.CPU,tmptask.CPU.Max,true)
			//tmptask.CPU.Min = Deal_CPU_Str_Compare(res.Containers[0].Usage.CPU,tmptask.CPU.Min,false)
		//}
		Podmp[res.Metadata.Name] = tmptask
	}
	return nil
}

func ReadGPUMemUsed(dir string) {
	dcgm_fb_used := ReadFile("./"+dir+"/dcgm_fb_used.log")


	for _,v := range dcgm_fb_used.Data.Result {
		tmptask, ok := Podmp[v.Metric.Pod_name]
		if !ok {
			continue
		}

		value := reflect.ValueOf(v.RValue)

		flag := false
		for k,_ := range tmptask.GPU.GPUMem {
			if tmptask.GPU.GPUMem[k].Uuid == v.Metric.Uuid {
				flag = true
				for i:=0;i<value.Len();i++ {
					gpumem,_ := strconv.ParseInt(value.Index(i).Elem().Index(1).Elem().String(),10,64)
					tmptask.GPU.GPUMem[k].MaxR = MAX(tmptask.GPU.GPUMem[k].MaxR,gpumem).(int64)
					//tmpMin = MIN(tmpMin,gpuuti)
					tmptask.GPU.GPUMem[k].History = append(tmptask.GPU.GPUMem[k].History,gpumem)
				}

				Podmp[v.Metric.Pod_name]=tmptask
			}
		}
		if flag {
			continue
		}

		var tmpgmemhis GPUMemHistory
		tmpgmemhis.Uuid = v.Metric.Uuid
		tmpgmemhis.Pod = v.Metric.Pod_name
		tmpgmemhis.Total = 0
		tmpgmemhis.MaxR = 0
		
		var tmpMax int64 = 0
		//var tmpMin int64 = INT_MAX

		for i:=0;i<value.Len();i++ {
			gpumem,_ := strconv.ParseInt(value.Index(i).Elem().Index(1).Elem().String(),10,64)
			tmpMax = MAX(tmpMax,gpumem).(int64)
			//tmpMin = MIN(tmpMin,gpumem)
			tmpgmemhis.History = append(tmpgmemhis.History,gpumem)
		}
		tmpgmemhis.MaxR = tmpMax
		//tmpgmemhis.Min=tmpMin
		tmptask.GPU.GPUMem = append(tmptask.GPU.GPUMem, tmpgmemhis)
		Podmp[v.Metric.Pod_name]=tmptask

	}
}

func ReadGPUMemFree(dir string) {
	dcgm_fb_free := ReadFile("./"+dir+"/dcgm_fb_free.log")


	for _,v := range dcgm_fb_free.Data.Result {
		tmptask, ok := Podmp[v.Metric.Pod_name]
		if !ok {
			continue
		}

		value := reflect.ValueOf(v.RValue)

		flag := false
		for k,_ := range tmptask.GPU.GPUMem {
			if tmptask.GPU.GPUMem[k].Uuid == v.Metric.Uuid {
				flag = true
				for i:=0;i<value.Len();i++ {
					gpumem,_ := strconv.ParseInt(value.Index(i).Elem().Index(1).Elem().String(),10,64)
					//tmpMax = MAX(tmpMax,gpuuti)
					//tmpMin = MIN(tmpMin,gpuuti)
					tmptask.GPU.GPUMem[k].Total = MAX(tmptask.GPU.GPUMem[k].History[i]+gpumem, tmptask.GPU.GPUMem[k].Total).(int64)
					break
				}
				
				Podmp[v.Metric.Pod_name]=tmptask
			}
		}
		if flag {
			continue
		}

		var tmpgmemhis GPUMemHistory
		tmpgmemhis.Uuid = v.Metric.Uuid
		tmpgmemhis.Pod = v.Metric.Pod_name
		tmpgmemhis.Total = 0
		//var tmpMax int64 = 0
		//var tmpMin int64 = INT_MAX

		for i:=0;i<value.Len();i++ {
			gpumem,_ := strconv.ParseInt(value.Index(i).Elem().Index(1).Elem().String(),10,64)
			//tmpMax = MAX(tmpMax,gpumem)
			//tmpMin = MIN(tmpMin,gpumem)
			tmpgmemhis.Total = gpumem
			break
		}
		//tmpgmemhis.Max=tmpMax
		//tmpgmemhis.Min=tmpMin
		tmptask.GPU.GPUMem = append(tmptask.GPU.GPUMem, tmpgmemhis)
		Podmp[v.Metric.Pod_name]=tmptask

	}
}

func ReadGPUMemmemCopyUtil(dir string) {
	dcgm_mem_copy_utilization := ReadFile("./"+dir+"/dcgm_mem_copy_utilization.log")

	for _,v := range dcgm_mem_copy_utilization.Data.Result {
		tmptask, ok := Podmp[v.Metric.Pod_name]
		if !ok || v.Metric.Pod_name == "" {
			continue
		}

		value := reflect.ValueOf(v.RValue)
		
		flag := false
		for k,_ := range tmptask.GPU.GPUMemCopy {
			if v.Metric.Uuid == "" {
				continue
			}
			if tmptask.GPU.GPUMemCopy[k].Uuid == v.Metric.Uuid {
				flag = true
				var tmpMax int64 = tmptask.GPU.GPUMemCopy[k].MaxR
				for i:=0;i<value.Len();i++ {
					
					gpumemutil,_ := strconv.ParseInt(value.Index(i).Elem().Index(1).Elem().String(),10,64)
					tmpMax = MAX(tmpMax,gpumemutil).(int64)
					
					//tmpMin = MIN(tmpMin,gpuuti)
					
					tmptask.GPU.GPUMemCopy[k].History = append(tmptask.GPU.GPUMemCopy[k].History,gpumemutil)
					
				}
				tmptask.GPU.GPUMemCopy[k].MaxR = tmpMax
				Podmp[v.Metric.Pod_name]=tmptask
			}
		}

		if flag {
			continue
		}

		var tmpghis GPUHistory
		tmpghis.Uuid = v.Metric.Uuid
		tmpghis.Pod = v.Metric.Pod_name
		tmpghis.MaxR = 0
		tmpghis.History = make([]int64, 0)
		var tmpMax int64 = 0
		//var tmpMin int64 = INT_MAX

		for i:=0;i<value.Len();i++ {
			gpumemutil,_ := strconv.ParseInt(value.Index(i).Elem().Index(1).Elem().String(),10,64)
			tmpMax = MAX(tmpMax , gpumemutil).(int64)
			//tmpMin = MIN(tmpMin,gpuuti)
			tmpghis.History = append(tmpghis.History,gpumemutil)
			
		}
		tmpghis.MaxR = tmpMax
		//tmpgmemhis.Min=tmpMin
		tmptask.GPU.GPUMemCopy = append(tmptask.GPU.GPUMemCopy, tmpghis)
		Podmp[v.Metric.Pod_name]=tmptask

	}
}

func ReadContainerTasksState(dir string) {
	container_tasks_state := ReadFile("./"+dir+"/container_tasks_state.log")

	for _,v := range container_tasks_state.Data.Result {
		if v.Metric.Container == "POD" || v.Metric.Container == "" || v.Metric.Namespace == "" || v.Metric.Namespace == "ingress-nginx" || v.Metric.Namespace == "kube-system" || v.Metric.Namespace == "default" || v.Metric.Namespace == "lens-metrics"{
			continue
		}

		value := reflect.ValueOf(v.RValue)
		if tmptask,ok := Podmp[v.Metric.Pod]; ok {
			tmptask.Starttime = MIN(tmptask.Starttime, int64(value.Index(0).Elem().Index(0).Elem().Float())).(int64)
			tmptask.Endtime = MAX(tmptask.Endtime, int64(value.Index(value.Len()-1).Elem().Index(0).Elem().Float())).(int64)
			Podmp[v.Metric.Pod] = tmptask
		} else {
			var tmp TaskLog
			tmp.Pod = v.Metric.Pod
			tmp.Namespace = v.Metric.Namespace
			tmp.ResourceT = v.Metric.ResourceT

			var tmpNode NodeInfo
			tmpNode.Name = v.Metric.Kbiohost
			tmp.Node = tmpNode


			tmp.Starttime = int64(value.Index(0).Elem().Index(0).Elem().Float())
			tmp.Endtime = int64(value.Index(value.Len()-1).Elem().Index(0).Elem().Float())

			Podmp[v.Metric.Pod] = tmp
		}

	}
}

func ReadDcgmGpuUtilization(dir string) {
	dcgm_gpu_utilization := ReadFile("./"+dir+"/dcgm_gpu_utilization.log")

	for _,v := range dcgm_gpu_utilization.Data.Result {
		tmptask, ok := Podmp[v.Metric.Pod_name]

		if _,ok := NodetoGPUtot[v.Metric.Name];!ok {
			NodetoGPUtot[v.Metric.Name] = make(map[int64]int64)
		}
		tmpv := NodetoGPUtot[v.Metric.Name]
		value := reflect.ValueOf(v.RValue)
		for i:=0;i<value.Len();i++ {
			timestamp := int64(value.Index(i).Elem().Index(0).Elem().Float())
			if _,ok := tmpv[timestamp];ok {
				tmpv[timestamp] = tmpv[timestamp] + 1
			} else {
				tmpv[timestamp] = 1
			}
		}
		NodetoGPUtot[v.Metric.Name] = tmpv

		if !ok {
			continue
		}

		if _,ok := NodetoGPUuse[v.Metric.Name];!ok {
			NodetoGPUuse[v.Metric.Name] = make(map[int64]int64)
		}
		tmpv = NodetoGPUuse[v.Metric.Name]

		flag := false
		for k,_ := range tmptask.GPU.GPUUtil {
			if tmptask.GPU.GPUUtil[k].Uuid == v.Metric.Uuid {
				flag = true
				for i:=0;i<value.Len();i++ {
					gpuuti,_ := strconv.ParseInt(value.Index(i).Elem().Index(1).Elem().String(),10,64)
					tmptask.GPU.GPUUtil[k].MaxR = MAX(tmptask.GPU.GPUUtil[k].MaxR,gpuuti).(int64)
					//tmpMin = MIN(tmpMin,gpuuti)
					tmptask.GPU.GPUUtil[k].History = append(tmptask.GPU.GPUUtil[k].History,gpuuti)

					timestamp := int64(value.Index(i).Elem().Index(0).Elem().Float())
					if _,ok := tmpv[timestamp];ok {
						tmpv[timestamp] = tmpv[timestamp] + 1
					} else {
						tmpv[timestamp] = 1
					}
				}
				NodetoGPUuse[v.Metric.Name] = tmpv
				Podmp[v.Metric.Pod_name]=tmptask
			}
		}
		if flag {
			continue
		}
		var tmpgutihis GPUHistory
		tmpgutihis.Uuid = v.Metric.Uuid
		tmpgutihis.Pod = v.Metric.Pod_name
		tmpgutihis.MaxR = 0
		
		var tmpMax int64 = 0
		//var tmpMin int64 = INT_MAX
		for i:=0;i<value.Len();i++ {
			gpuuti,_ := strconv.ParseInt(value.Index(i).Elem().Index(1).Elem().String(),10,64)
			tmpMax = MAX(tmpMax,gpuuti).(int64)
			//tmpMin = MIN(tmpMin,gpuuti)
			tmpgutihis.History = append(tmpgutihis.History,gpuuti)

			timestamp := int64(value.Index(i).Elem().Index(0).Elem().Float())
			if _,ok := tmpv[timestamp];ok {
				tmpv[timestamp] = tmpv[timestamp] + 1
			} else {
				tmpv[timestamp] = 1
			}
		}
		NodetoGPUuse[v.Metric.Name] = tmpv
		tmpgutihis.MaxR = tmpMax
		//tmpgutihis.Min=tmpMin
		tmptask.GPU.GPUUtil = append(tmptask.GPU.GPUUtil, tmpgutihis)

		Podmp[v.Metric.Pod_name]=tmptask
	}
}

func ReadPodContainerResourceLimits(dir string) (error){
	var err error

	kube_pod_container_resource_limits := ReadFile("./"+dir+"/kube_pod_container_resource_limits.log")

	for _,v := range kube_pod_container_resource_limits.Data.Result {
		tmptask, ok := Podmp[v.Metric.Pod]
		if !ok {
			continue
		}
		value := reflect.ValueOf(v.RValue)
		if v.Metric.Resource == "memory" {
			tmpmem := tmptask.Memory
			tmpmem.Pod = tmptask.Pod
			tmpmem.Node = tmptask.Node

			tmpmem.Limit, err = strconv.ParseInt(value.Index(0).Elem().Index(1).Elem().String(),10,64)
			if Err_Handle(err) {
				return err
			}

			tmptask.Memory = tmpmem

		} else if v.Metric.Resource == "nvidia_com_gpu" {
			tmpgpu := tmptask.GPU
			tmpgpu.Pod = tmptask.Pod
			tmpgpu.Node = tmptask.Node

			tmpgpu.NumGPU,err = strconv.ParseInt(value.Index(0).Elem().Index(1).Elem().String(),10,64)
			if Err_Handle(err) {
				return err
			}

			tmptask.GPU = tmpgpu

		} else if v.Metric.Resource == "cpu" {
			tmpcpu := tmptask.CPU
			tmpcpu.Node = tmptask.Node
			tmpcpu.Pod = tmptask.Pod

			tmpcpu.Limit, err = strconv.ParseInt(value.Index(0).Elem().Index(1).Elem().String(),10,64)
			if Err_Handle(err) {
				return err
			}

			tmptask.CPU = tmpcpu
		}
		namepos := strings.Index(v.Metric.Container,"-")
		tmptask.Container = v.Metric.Container
		tmptask.User = v.Metric.Container[:namepos]
		Podmp[v.Metric.Pod]=tmptask
	}

	return nil
}

func ReadContainerMemoryUsageBytes(dir string) {
	container_memory_usage_bytes := ReadFile("./"+dir+"/container_memory_usage_bytes.log")

	for _,v := range container_memory_usage_bytes.Data.Result {
		tmptask, ok := Podmp[v.Metric.Pod]
		if !ok {
			continue
		}
		if v.Metric.Container == "" || v.Metric.Container == "POD" {
			continue
		}
		value := reflect.ValueOf(v.RValue)

		//var tmpMax int64 = 0
		//var tmpMin int64 = INT_MAX

		for i:=0; i<value.Len(); i++ {
			memuse, _ := strconv.ParseInt(value.Index(i).Elem().Index(1).Elem().String(),10,64)
			//tmpMax = MAX(memuse,tmpMax)
			//tmpMin = MIN(memuse,tmpMin)
			tmptask.Memory.History = append(tmptask.Memory.History,memuse)
		}
		//tmptask.Memory.Max = tmpMax
		//tmptask.Memory.Min = tmpMin

		Podmp[v.Metric.Pod]=tmptask
	}
}

func Deal_Oneday_data( dir string, Nodeio map[string]NodeIO) (error) {
	
	ReadContainerTasksState(dir)

	ReadPodContainerResourceLimits(dir)

	ReadDcgmGpuUtilization(dir)

	// if err := ReadPodLifecycle(dir);err != nil {
	// 	fmt.Println(err)
	// }

	ReadGPUMemUsed(dir)

	ReadGPUMemFree(dir)

	ReadGPUMemmemCopyUtil(dir)

	ReadContainerMemoryUsageBytes(dir)

	rate_container_cpu_usage_seconds_total := ReadFile("./"+dir+"/rate(container_cpu_usage_seconds_total%5B1m%5D).log")

	for _,v := range rate_container_cpu_usage_seconds_total.Data.Result {
		tmptask, ok := Podmp[v.Metric.Pod]
		if !ok {
			continue
		}
		if v.Metric.Container == "" || v.Metric.Container == "POD" {
			continue
		}
		value := reflect.ValueOf(v.RValue)

		for i:=0; i<value.Len(); i++ {
			memuse, _ := strconv.ParseFloat(value.Index(i).Elem().Index(1).Elem().String(),10)
			//tmpMax = MAX(memuse,tmpMax)
			//tmpMin = MIN(memuse,tmpMin)
			tmptask.CPU.History = append(tmptask.CPU.History,memuse)
		}
		Podmp[v.Metric.Pod]=tmptask
	}
	

	rate_node_network_receive_bytes_total := ReadFile("./"+dir+"/rate%20(container_network_receive_bytes_total%7Bid%3D%22%2F%22%7D%5B1m%5D).log")
	//rate_container_network_receive_bytes_total := ReadFile("./rate%20(container_network_receive_bytes_total%7Bid%3D%22%2F%22%7D%5B1m%5D).log")

	for _,v := range rate_node_network_receive_bytes_total.Data.Result {
		tmpnode,ok := Nodeio[v.Metric.Kbiohost]
		if !ok {
			var tmp NodeIO
			tmp.Node = v.Metric.Kbiohost
			tmp.IRate = make(map[int64]float64)
			tmp.ORate = make(map[int64]float64)
			tmp.IbIRate = make(map[int64]float64)
			tmp.IbORate = make(map[int64]float64)
			tmpnode = tmp
		}
		value := reflect.ValueOf(v.RValue)

		for i:=0; i<value.Len(); i++ {
			timestamp := int64(value.Index(i).Elem().Index(0).Elem().Float())
			iorate, _ := strconv.ParseFloat(value.Index(i).Elem().Index(1).Elem().String(),10)

			timestampstart = MIN(timestamp,timestampstart).(int64)
			timestampend = MAX(timestamp,timestampend).(int64)
			if _,ok := tmpnode.IRate[timestamp];ok {
				tmpnode.IRate[timestamp] = tmpnode.IRate[timestamp] + iorate
			} else {
				tmpnode.IRate[timestamp] = iorate
			}
		}
		Nodeio[v.Metric.Kbiohost] = tmpnode
	}

	rate_node_network_transmit_bytes_total := ReadFile("./"+dir+"/rate%20(container_network_transmit_bytes_total%7Bid%3D%22%2F%22%7D%5B1m%5D).log")
	//rate_container_network_transmit_bytes_total := ReadFile("./rate%20(container_network_transmit_bytes_total%7Bid%3D%22%2F%22%7D%5B1m%5D).log")

	for _,v := range rate_node_network_transmit_bytes_total.Data.Result {
		tmpnode,ok := Nodeio[v.Metric.Kbiohost]
		if !ok {
			var tmp NodeIO
			tmp.Node = v.Metric.Kbiohost
			tmp.IRate = make(map[int64]float64)
			tmp.ORate = make(map[int64]float64)
			tmpnode = tmp
		}
		value := reflect.ValueOf(v.RValue)

		for i:=0; i<value.Len(); i++ {
			timestamp := int64(value.Index(i).Elem().Index(0).Elem().Float())
			iorate, _ := strconv.ParseFloat(value.Index(i).Elem().Index(1).Elem().String(),10)

			timestampstart = MIN(timestamp,timestampstart).(int64)
			timestampend = MAX(timestamp,timestampend).(int64)
			if _,ok := tmpnode.ORate[timestamp];ok {
				tmpnode.ORate[timestamp] = tmpnode.ORate[timestamp] + iorate
			} else {
				tmpnode.ORate[timestamp] = iorate
			}
		}
		Nodeio[v.Metric.Kbiohost] = tmpnode
	}

	rate_infiniband_node_network_receive_bytes_total := ReadFile("./"+dir+"/rate(node_infiniband_port_data_received_bytes_total%5B1m%5D).log")
	//rate_container_network_receive_bytes_total := ReadFile("./rate%20(container_network_receive_bytes_total%7Bid%3D%22%2F%22%7D%5B1m%5D).log")

	for _,v := range rate_infiniband_node_network_receive_bytes_total.Data.Result {
		tmpnode,ok := Nodeio[v.Metric.Name]
		if !ok {
			var tmp NodeIO
			tmp.Node = v.Metric.Name
			tmp.IRate = make(map[int64]float64)
			tmp.ORate = make(map[int64]float64)
			tmp.IbIRate = make(map[int64]float64)
			tmp.IbORate = make(map[int64]float64)
			tmpnode = tmp
		}
		value := reflect.ValueOf(v.RValue)

		for i:=0; i<value.Len(); i++ {
			timestamp := int64(value.Index(i).Elem().Index(0).Elem().Float())
			iorate, _ := strconv.ParseFloat(value.Index(i).Elem().Index(1).Elem().String(),10)

			timestampstart = MIN(timestamp,timestampstart).(int64)
			timestampend = MAX(timestamp,timestampend).(int64)
			if _,ok := tmpnode.IbIRate[timestamp];ok {
				tmpnode.IbIRate[timestamp] = tmpnode.IbIRate[timestamp] + iorate
			} else {
				tmpnode.IbIRate[timestamp] = iorate
			}
		}
		Nodeio[v.Metric.Name] = tmpnode
	}

	rate_infiniband_node_network_transmit_bytes_total := ReadFile("./"+dir+"/rate(node_infiniband_port_data_transmitted_bytes_total%5B1m%5D).log")
	//rate_container_network_transmit_bytes_total := ReadFile("./rate%20(container_network_transmit_bytes_total%7Bid%3D%22%2F%22%7D%5B1m%5D).log")

	for _,v := range rate_infiniband_node_network_transmit_bytes_total.Data.Result {
		tmpnode,ok := Nodeio[v.Metric.Name]
		if !ok {
			var tmp NodeIO
			tmp.Node = v.Metric.Name
			tmp.IRate = make(map[int64]float64)
			tmp.ORate = make(map[int64]float64)
			tmp.IbIRate = make(map[int64]float64)
			tmp.IbORate = make(map[int64]float64)
			tmpnode = tmp
		}
		value := reflect.ValueOf(v.RValue)

		for i:=0; i<value.Len(); i++ {
			timestamp := int64(value.Index(i).Elem().Index(0).Elem().Float())
			iorate, _ := strconv.ParseFloat(value.Index(i).Elem().Index(1).Elem().String(),10)

			timestampstart = MIN(timestamp,timestampstart).(int64)
			timestampend = MAX(timestamp,timestampend).(int64)
			if _,ok := tmpnode.IbORate[timestamp];ok {
				tmpnode.IbORate[timestamp] = tmpnode.IbORate[timestamp] + iorate
			} else {
				tmpnode.IbORate[timestamp] = iorate
			}
		}
		Nodeio[v.Metric.Name] = tmpnode
	}
	return nil
}

func OuttoFile(Nodeio map[string]NodeIO) {
	nodeinfo, err := os.Create("./result/nodeinfo.csv")
	if err != nil {
		fmt.Println("node File creating error", err)
		return
	}
	nodeinfo.WriteString("Nodename,Podname\n")
	for _,v := range Podmp {
		nodeinfo.WriteString(v.Node.Name+","+v.Pod+"\n")
	}
	nodeinfo.Close()

	taskinfo, err := os.Create("./result/taskinfo.csv")
	if err != nil {
		fmt.Println("task File creating error", err)
		return
	}
	taskinfo.WriteString("Podname,Containername,Namespace,User,ResourceType,Nodename,")
	taskinfo.WriteString("Starttime,Endtime,CPULimit,MemoryLimit,GPULimit\n")
	for _,v := range Podmp {
		taskinfo.WriteString(v.Pod+","+v.Container+","+v.Namespace+","+v.User+","+v.ResourceT+","+v.Node.Name+",")
		taskinfo.WriteString(fmt.Sprintf("%d,%d,%d,%d,%d\n",v.Starttime,v.Endtime,v.CPU.Limit, v.Memory.Limit, v.GPU.NumGPU))
	}
	taskinfo.Close()

	cpuinfo, err := os.Create("./result/cpuinfo.json")
	if err != nil {
		fmt.Println("CPU File creating error", err)
		return
	}
	//cpuinfo.WriteString("Podname,Max,Min,History\n")
	enc = json.NewEncoder(cpuinfo)
	for _,v := range Podmp {
		//cpuinfo.WriteString(v.Pod)
		//cpuinfo.WriteString(","+v.CPU.Max+","+v.CPU.Min)
		//for _,vv := range v.CPU.History {
		//	cpuinfo.WriteString(","+vv)
		//}
		//cpuinfo.WriteString("\n")
		if v.CPU.Pod == "" {
			continue
		}
		err = enc.Encode(v.CPU)
		Err_Handle(err)
	}
	cpuinfo.Close()

	gpuinfomem, err := os.Create("./result/gpuinfomem.json")

	if err != nil {
		fmt.Println("GPU mem File creating error", err)
		return
	}
	//gpuinfomem.WriteString("Podname,Uuid,Max,Min,History\n")
	enc = json.NewEncoder(gpuinfomem)
	for _,v := range Podmp {
		for _,vv := range v.GPU.GPUMem {
			//gpuinfomem.WriteString(v.Pod + "," + vv.Uuid)
			//gpuinfomem.WriteString(","+fmt.Sprint(vv.Max,",",vv.Min))
			//for _, vvv := range vv.History {
			//	gpuinfomem.WriteString(","+fmt.Sprint(vvv))
			//}
			//gpuinfomem.WriteString("\n")
			err = enc.Encode(vv)
			Err_Handle(err)
		}

	}
	gpuinfomem.Close()

	gpuinfoutil, err := os.Create("./result/gpuinfoutil.json")
	if err != nil {
		fmt.Println("GPU util File creating error", err)
		return
	}
	//gpuinfoutil.WriteString("Podname,Uuid,Max,Min,History\n")
	enc = json.NewEncoder(gpuinfoutil)

	for _,v := range Podmp {
		for _,vv := range v.GPU.GPUUtil {
			//gpuinfoutil.WriteString(v.Pod + "," + vv.Uuid)
			//gpuinfoutil.WriteString(","+fmt.Sprint(vv.Max,",",vv.Min))
			//for _, vvv := range vv.History {
			//	gpuinfoutil.WriteString(fmt.Sprint(",",vvv))
			//}
			//gpuinfoutil.WriteString("\n")
			err = enc.Encode(vv)
			Err_Handle(err)
		}


	}
	gpuinfoutil.Close()

	gpumemcpyutil, err := os.Create("./result/gpumemcpyutil.json")
	if err != nil {
		fmt.Println("GPU Memcpy util File creating error", err)
		return
	}
	//gpumemcpyutil.WriteString("Podname,Uuid,Max,Min,History\n")
	enc = json.NewEncoder(gpumemcpyutil)

	for _,v := range Podmp {
		for _,vv := range v.GPU.GPUMemCopy {
			//gpuinfoutil.WriteString(v.Pod + "," + vv.Uuid)
			//gpuinfoutil.WriteString(","+fmt.Sprint(vv.Max,",",vv.Min))
			//for _, vvv := range vv.History {
			//	gpuinfoutil.WriteString(fmt.Sprint(",",vvv))
			//}
			//gpuinfoutil.WriteString("\n")
			err = enc.Encode(vv)
			Err_Handle(err)
		}


	}
	gpumemcpyutil.Close()

	meminfo, err := os.Create("./result/meminfo.json")
	if err != nil {
		fmt.Println("Mem File creating error", err)
		return
	}
	//meminfo.WriteString("Podname,Max,Min,History\n")

	enc = json.NewEncoder(meminfo)
	for _,v := range Podmp {
		//meminfo.WriteString(v.Pod)
		//meminfo.WriteString(","+fmt.Sprint(v.Memory.Max,",",v.Memory.Min))
		//for _, vv := range v.Memory.History {
		//	meminfo.WriteString(fmt.Sprint(",", vv))
		//}
		//meminfo.WriteString("\n")


		err = enc.Encode(v.Memory)
		Err_Handle(err)

	}
	meminfo.Close()

	nodegpuinfo, err := os.Create("./result/nodegpuinfo.json")
	if err != nil {
		fmt.Println("nodegpu File creating error", err)
		return
	}

	enc = json.NewEncoder(nodegpuinfo)
	for k,v := range NodetoGPUtot {
		var vv NodeGPUstate
		vv.Node = k
		timestamp := make([]int,0)
		for t,_ := range v {
			timestamp = append(timestamp,int(t))
		}
		sort.Ints(timestamp)

		for _,idx := range timestamp {
			var vvv NodeGPUstateElem
			vvv.Total = v[int64(idx)]
			//vvv.Time = idx
			if gpuuse,ok := NodetoGPUuse[k]; ok {
				vvv.Use=gpuuse[int64(idx)]
			} else {
				vvv.Use=0
			}
			vv.State = append(vv.State,vvv)
		}
		err = enc.Encode(vv)
		Err_Handle(err)
	}
	nodegpuinfo.Close()

	nodegpuuse, err := os.Create("./result/nodegpuuse.csv")
	if err != nil {
		fmt.Println("nodegpuuse File creating error", err)
		return
	}

	for _,v := range Podmp {
		for _,vv := range v.GPU.GPUUtil {

			var cnt float64 = 0
			for i:=0; i< len(vv.History); i++ {
				if vv.History[i] == 0 {
					cnt++
				}
			}

			ratio := cnt *100 / float64(len(vv.History))
			nodegpuuse.WriteString(fmt.Sprint(vv.Pod,",",vv.Uuid,",",ratio)+"\n")
		}
	}
	nodegpuuse.Close()

	nodecpuuti, err := os.Create("./result/nodecpuuti.json")
	if err != nil {
		fmt.Println("nodecpuuti File creating error", err)
		return
	}

	enc = json.NewEncoder(nodecpuuti)
	for _,v := range Podmp {
		if v.CPU.Limit == 0 {
			continue
		}
		var cpuco CPUCore
		cpuco.Node = v.CPU.Node
		cpuco.Pod = v.CPU.Pod
		
		for i:=0; i<len(v.CPU.History); i++ {
			cpuco.Utilization = append(cpuco.Utilization, v.CPU.History[i]*100/float64(v.CPU.Limit))
			
		}
		
		
		err = enc.Encode(cpuco)
		Err_Handle(err)
	}
	nodecpuuti.Close()

	nodereiverate, err := os.Create("./result/nodereiverate.json")

	if err != nil {
		fmt.Println("tmp File creating error", err)
		return
	}

	enc = json.NewEncoder(nodereiverate)
	for _,v := range Nodeio {
		var tmprate struct{
			Node string `json:"node"`
			Rate []float64 `json:"Irate"`
		}
		tmprate.Node = v.Node
		tmprate.Rate = make([]float64,0)
		
		for k := timestampstart; k<= timestampend; k+=30 {
			if _,ok := v.IRate[k]; ok {
				tmprate.Rate = append(tmprate.Rate,v.IRate[k])
			} else {
				tmprate.Rate = append(tmprate.Rate,0)
			}
			
		}
		
		//fmt.Println(tmpnoderate)
		err = enc.Encode(tmprate)
		Err_Handle(err)
	}
	nodereiverate.Close()

	nodetransrate, err := os.Create("./result/nodetransrate.json")

	if err != nil {
		fmt.Println("nodetransrate File creating error", err)
		return
	}
	
	enc = json.NewEncoder(nodetransrate)
	for _,v := range Nodeio {
		var tmprate struct{
			Node string `json:"node"`
			Rate []float64 `json:"Orate"`
		}
		tmprate.Node = v.Node
		tmprate.Rate = make([]float64,0)

		for k := timestampstart; k<= timestampend; k+=30 {
			if _,ok := v.ORate[k]; ok {
				tmprate.Rate = append(tmprate.Rate,v.ORate[k])
			} else {
				tmprate.Rate = append(tmprate.Rate,0)
			}
			
		}
		
		err = enc.Encode(tmprate)
		Err_Handle(err)
	}
	nodetransrate.Close()

	nodeibreiverate, err := os.Create("./result/nodeibreiverate.json")

	if err != nil {
		fmt.Println("tmp File creating error", err)
		return
	}

	enc = json.NewEncoder(nodeibreiverate)
	
	for _,v := range Nodeio {
		var tmprate struct{
			Node string `json:"node"`
			Rate []float64 `json:"Iibrate"`
		}
		tmprate.Node = v.Node
		tmprate.Rate = make([]float64,0)
		for k := timestampstart; k<= timestampend; k+=30 {
			if _,ok := v.IbIRate[k]; ok {
				tmprate.Rate = append(tmprate.Rate,v.IbIRate[k])
			} else {
				tmprate.Rate = append(tmprate.Rate,0)
			}
			
		}
		
		//fmt.Println(tmpnoderate)
		err = enc.Encode(tmprate)
		Err_Handle(err)
	}
	nodereiverate.Close()

	nodeibtransrate, err := os.Create("./result/nodeibtransrate.json")

	if err != nil {
		fmt.Println("nodetransrate File creating error", err)
		return
	}

	enc = json.NewEncoder(nodeibtransrate)
	
	for _,v := range Nodeio {
		var tmprate struct{
			Node string `json:"node"`
			Rate []float64 `json:"Oibrate"`
		}
		tmprate.Node = v.Node
		tmprate.Rate = make([]float64,0)

		for k := timestampstart; k<= timestampend; k+=30 {
			if _,ok := v.IbORate[k]; ok {
				tmprate.Rate = append(tmprate.Rate,v.IbORate[k])
			} else {
				tmprate.Rate = append(tmprate.Rate,0)
			}
			
		}
		
		err = enc.Encode(tmprate)
		Err_Handle(err)
	}
	nodetransrate.Close()
	
}
