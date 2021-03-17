package collector

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	getNameCmd        = "nvidia-smi --query-gpu=name --format=csv,noheader,nounits -i "
	getTemperatureCmd = "nvidia-smi --query-gpu=temperature.gpu --format=csv,noheader,nounits -i "
	getMemUsedCmd     = "nvidia-smi --query-gpu=memory.used --format=csv,noheader,nounits -i "
	getMemFreeCmd     = "nvidia-smi --query-gpu=memory.free --format=csv,noheader,nounits -i "
	getMemTotalCmd    = "nvidia-smi --query-gpu=memory.total --format=csv,noheader,nounits -i "
	getUtilizationCmd = "nvidia-smi --query-gpu=utilization.gpu --format=csv,noheader,nounits -i "
	// getFanSpeedCmd       = "nvidia-smi --query-gpu=fan.speed --format=csv,noheader,nounits -i "
	// getPowerCmd          = "nvidia-smi --query-gpu=power.draw --format=csv,noheader,nounits -i "
)

type GPUInfoCollector struct {
	Seq            string
	Temperature    *prometheus.Desc
	MemUsed        *prometheus.Desc
	MemFree        *prometheus.Desc
	MemTotal       *prometheus.Desc
	GPUUtilization *prometheus.Desc
}

func NewGPUInfoCollector(num string, name string) *GPUInfoCollector {
	return &GPUInfoCollector{
		Seq: num,
		Temperature: prometheus.NewDesc(
			"gpu_temperature",
			"Shows temperature about gpu",
			nil,
			prometheus.Labels{"gpu_seq": num, "name": name}),
		MemUsed: prometheus.NewDesc(
			"gpu_memory_used",
			"Shows gpu memory used (MiB)",
			nil,
			prometheus.Labels{"gpu_seq": num, "name": name}),
		MemFree: prometheus.NewDesc(
			"gpu_memory_free",
			"Shows gpu memory free (MiB)",
			nil,
			prometheus.Labels{"gpu_seq": num, "name": name}),
		MemTotal: prometheus.NewDesc(
			"gpu_memory_total",
			"Shows gpu memory total (MiB)",
			nil,
			prometheus.Labels{"gpu_seq": num, "name": name}),
		GPUUtilization: prometheus.NewDesc(
			"gpu_utilization",
			"Shows gpu utilization (%)",
			nil,
			prometheus.Labels{"gpu_seq": num, "name": name}),
	}
}

func (c *GPUInfoCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.Temperature
	ch <- c.MemUsed
	ch <- c.MemFree
	ch <- c.MemTotal
	ch <- c.GPUUtilization
}

func (c *GPUInfoCollector) Collect(ch chan<- prometheus.Metric) {

	tempVal := getGPUTemperature(c.Seq)
	usedVal := getGPUMemUsed(c.Seq)
	freeVal := getGPUMemFree(c.Seq)
	totalVal := getGPUMemTotal(c.Seq)
	utilizationVal := getGPUUtilization(c.Seq)

	ch <- prometheus.MustNewConstMetric(c.Temperature, prometheus.GaugeValue, tempVal)
	ch <- prometheus.MustNewConstMetric(c.MemUsed, prometheus.GaugeValue, usedVal)
	ch <- prometheus.MustNewConstMetric(c.MemFree, prometheus.GaugeValue, freeVal)
	ch <- prometheus.MustNewConstMetric(c.MemTotal, prometheus.GaugeValue, totalVal)
	ch <- prometheus.MustNewConstMetric(c.GPUUtilization, prometheus.GaugeValue, utilizationVal)
}

// 获取显卡温度
func getGPUTemperature(n string) float64 {
	_, s := ExecCommand(getTemperatureCmd + n)
	t, _ := strconv.ParseFloat(s, 64)
	return t
}

// 获取使用的显存 MiB
func getGPUMemUsed(n string) float64 {
	_, s := ExecCommand(getMemUsedCmd + n)
	m, _ := strconv.ParseFloat(s, 64)
	return m
}

// 获取空闲的显存 MiB
func getGPUMemFree(n string) float64 {
	_, s := ExecCommand(getMemFreeCmd + n)
	m, _ := strconv.ParseFloat(s, 64)
	return m
}

// 获取总的显存 MiB
func getGPUMemTotal(n string) float64 {
	_, s := ExecCommand(getMemTotalCmd + n)
	m, _ := strconv.ParseFloat(s, 64)
	return m
}

// 获取 GPU 使用率
func getGPUUtilization(n string) float64 {
	_, s := ExecCommand(getUtilizationCmd + n)
	m, _ := strconv.ParseFloat(s, 64)
	return m
}

// 获取显卡名字
func getGPUName(n string) string {
	_, s := ExecCommand(getNameCmd + n)
	return s
}

func getGPUNums() int {
	_, n := ExecCommand("nvidia-smi -L | wc -l")
	nums, _ := strconv.Atoi(n)
	return nums
}

func init() {
	nums := getGPUNums()
	for i := 0; i < nums; i++ {
		n := strconv.Itoa(i)
		name := getGPUName(n)
		prometheus.MustRegister(NewGPUInfoCollector(n, name))
	}
}
