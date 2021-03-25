package main

import (
	"time"
	"fmt"
)

type MetricInfo struct {
	/*infiniband
	Device string `json:"device"`
	Instance string `json:"instance"`
	Name string `json:"name"`
	*/

	//container_tasks_state
	Container string `json:"container"`
	Id string `json:"id"`
	Image string `json:"image"`
	Instance string `json:"instance"`
	Kbiohost string `json:"kubernetes_io_hostname"`
	Name string `json:"name"`
	Namespace string `json:"namespace"`
	Octopus_ftp_daemon string `json:"octopus_openi_pcl_cn_ftp_daemon"`
	Octopus_node string `json:"octopus_openi_pcl_cn_node"`
	Pod string `json:"pod"`
	ResourceT string `json:"resourceType"`
	//State string `json:"state"`
	Node string `json:"node"`
	Resource string `json:"resource"`
	Container_name string `json:"container_name"`
	Pod_name string `json:"pod_name"`
	Pod_namespace string `json:"pod_namespace"`
	Uuid string `json:"uuid"`

	/*kube_pod_init_container_status_waiting_reason
	Container string `json:"container"`
	Instance string `json:"instance"`
	Name string `json:"name"`
	Namespace string `json:"namespace"`
	Pod string `json:"pod"`
	Reason string `json:"reason"`
	*/
	/*container_start_time_seconds
	Container string `json:"container"`
	Id string `json:"id"`
	Image string `json:"image"`
	Instance string `json:"instance"`
	Kbiohost string `json:"kubernetes_io_hostname"`
	Name string `json:"name"`
	Namespace string `json:"namespace"`
	Octopus_node string `json:"octopus_openi_pcl_cn_node"`
	Pod string `json:"pod"`
	ResourceT string `json:"resourceType"`

	Instance string `json:"instance"`
	Name string `json:"name"`
	Namespace string `json:"namespace"`
	Phase string `json:"phase"`
	Pod string `json:"pod"`
	*/
}

type ResultInfo struct {
	Metric MetricInfo `json:"metric"`
	RValue []interface{} `json:"values"`
}
type DataInfo struct {
	ResultType string `json:"resultType"`
	Result []ResultInfo `json:"result"`
}
type PrometheusInfo struct {
	Status string `json:"status"`
	Data DataInfo `json:"data"`
}

type NodeInfo struct {
	Name string	`json:"name"`

}

type MemoryInfo struct {
	Node NodeInfo	`json:"node"`
	Pod string	`json:"pod"`
	Limit int64	`json:"limit"`
	//Max int64	`json:"max"`
	//Min int64	`json:"min"`
	History []int64	`json:"history"`
}

type CPUInfo struct {
	Node NodeInfo `json:"node"`
	Pod string	`json:"pod"`
	Limit int64	`json:"limit"`
	//Max string	`json:"max"`
	//Min string	`json:"min"`
	History []float64	`json:"history"`

}

type CPUCore struct {
	Node NodeInfo `json:"node"`
	Pod string	`json:"pod"`
	Utilization []float64 `json:"utili"`
}

type GPUMemHistory struct {
	Pod string `json:"pod"`
	Uuid string	`json:"uuid"`
	Total int64 `json:"total"`
	MaxR int64	`json:"maxratio"`
	//Min int64	`json:"min"`
	History []int64	`json:"history"`
}

type GPUHistory struct {
	Pod string `json:"pod"`
	Uuid string	`json:"uuid"`
	MaxR int64	`json:"maxratio"`
	//Min int64	`json:"min"`
	History []int64	`json:"history"`
}

/*type GPUutilratio struct {
	Pod string `json:"pod"`
	Uuid string	`json:"uuid"`
	Ratio float64 `json:"ratio"`
}*/

type GPUInfo struct {
	Node NodeInfo	`json:"nodename"`
	Pod string	`json:"podname"`
	GPUUtil []GPUHistory	`json:"gpuutil"`
	NumGPU int64	`json:"numgpu"`
	GPUMem []GPUMemHistory	`json:"gpumem"`
	//Ratio []GPUutilratio
}


type TaskLog struct {
	JobName string	`json:"jobname"`
	Namespace string `json:"namespace"`
	Starttime int64	`json:"starttime"`
	Endtime int64	`json:"endtime"`
	SubmitTime int64	`json:"submittime"`
	Container string	`json:"container"`
	User string	`json:"user"`
	Pod string	`json:"pod"`
	Node NodeInfo	`json:"node"`
	GPU GPUInfo	`json:"gpu"`
	CPU CPUInfo	`json:"cpu"`
	Memory MemoryInfo	`json:"memory"`
	ResourceT string	`json:"resourcetype"`
}

type PodMetricsList struct {
	Metadata struct {
		Name              string    `json:"name"`
		Namespace         string    `json:"namespace"`
		SelfLink          string    `json:"selfLink"`
		CreationTimestamp time.Time `json:"creationTimestamp"`
	} `json:"metadata"`
	Timestamp  time.Time `json:"timestamp"`
	Window     string    `json:"window"`
	Containers []struct {
		Name  string `json:"name"`
		Usage struct {
			CPU    string `json:"cpu"`
			Memory string `json:"memory"`
			// add non exist GPU cnt field, fill by get pod
			GPUCnt int64  `json:"gpu_cnt, omitempty"`
		} `json:"usage"`
	} `json:"containers"`
}

type NodeIO struct {
	Node string `json:"node"`
	IRate map[int64]float64 `json:"recieve_rate"`
	ORate map[int64]float64 `json:"transmit_rate"`
	IbIRate map[int64]float64 `json:"ib_recieve_rate"`
	IbORate map[int64]float64 `json:"ib_transmit_rate"`
}

type NodeGPUstateElem struct {
	Total int64 `json:"total"`
	Use int64 `json:"use"`
	//Time int `json:"time"`
}

type NodeGPUstate struct {
	Node string `json:"name"`
	State []NodeGPUstateElem `json:"state"`
}

type Bar struct {
    percent int64  //百分比
    cur     int64  //当前进度位置
    total   int64  //总进度
    rate    string //进度条
    graph   string //显示符号
}

func (bar *Bar) NewOption(start, total int64) {
    bar.cur = start
    bar.total = total
    if bar.graph == "" {
        bar.graph = "█"
    }
    bar.percent = bar.getPercent()
    for i := 0; i < int(bar.percent); i += 2 {
        bar.rate += bar.graph //初始化进度条位置
    }
	fmt.Printf("\r[%-50s]%3d%%  %8d/%d", bar.rate, bar.percent, bar.cur, bar.total)
}

func (bar *Bar) getPercent() int64 {
    return int64(float32(bar.cur) / float32(bar.total) * 100)
}


func (bar *Bar) Play(cur int64) {
    bar.cur = cur
    last := bar.percent
    bar.percent = bar.getPercent()
    
	for i := last; i < int64(bar.percent); i += 2 {
        bar.rate += bar.graph 
    }
    fmt.Printf("\r[%-50s]%3d%%  %8d/%d", bar.rate, bar.percent, bar.cur, bar.total)
}