package state

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"
	"github.com/Auginte/go-monitoring/domain/common"
)

func currentTime() string {
	return time.Now().Format(time.RFC3339Nano)
}

// CPU related

type procStat struct {
	Type              string `json:"type"`
	Time              string `json:"time"`
	CPUUser           int    `json:"CpuUser"`
	CPUNice           int    `json:"CpuNice"`
	CPUSystem         int    `json:"CpuSystem"`
	CPUIdle           int    `json:"CpuIdle"`
	CPUIOwait         int    `json:"CpuIOwait"`
	CPUInterrupts     int    `json:"CpuInterrupts"`
	CPUSoftInterrupts int    `json:"CpuSoftInterrupts"`
	Processes         int
	ProcessesRunning  int
	ProcessesBlocked  int
}

func parseProcStat(line string, data *procStat) {
	var tmp int
	line = strings.Replace(line, "  ", " ", -1)
	parts := strings.SplitN(line, " ", 2)
	group := parts[0]
	rest := parts[1]
	switch group {
	case "cpu":
		fmt.Sscanf(
			rest,
			"%d %d %d %d %d %d %d %d %d %d",
			&data.CPUUser,
			&data.CPUNice,
			&data.CPUSystem,
			&data.CPUIdle,
			&data.CPUIOwait,
			&data.CPUInterrupts,
			&data.CPUSoftInterrupts,
			&tmp,
			&tmp,
			&tmp,
		)
	case "processes":
		fmt.Sscanf(rest, "%d", &data.Processes)
	case "procs_running":
		fmt.Sscanf(rest, "%d", &data.ProcessesRunning)
	case "procs_blocked":
		fmt.Sscanf(rest, "%d", &data.ProcessesBlocked)
	}
}

// GetProcStatJSON - data from /proc/stat
func GetProcStatJSON() string {
	file, err := os.Open("/proc/stat")
	common.LogError(err)
	scanner := bufio.NewScanner(file)
	data := &procStat{}
	data.Type = "proc.stat"
	data.Time = currentTime()
	for scanner.Scan() {
		line := scanner.Text()
		parseProcStat(line, data)
	}
	common.LogError(scanner.Err())
	common.LogError(file.Close())

	serialised, err := json.Marshal(data)
	common.LogError(err)
	return string(serialised)
}

// Memory related

type procMemInfo struct {
	Type         string `json:"type"`
	Time         string `json:"time"`
	MemTotal     int
	MemFree      int
	MemAvailable int
	SwapTotal    int
	SwapFree     int
}

func parseProcMemInfo(line string, data *procMemInfo) {
	parts := strings.SplitN(line, ":", 2)
	group := parts[0]
	rest := strings.Trim(parts[1], " kB")
	switch group {
	case "MemTotal":
		fmt.Sscanf(rest, "%d", &data.MemTotal)
	case "MemFree":
		fmt.Sscanf(rest, "%d", &data.MemFree)
	case "MemAvailable":
		fmt.Sscanf(rest, "%d", &data.MemAvailable)
	case "SwapTotal":
		fmt.Sscanf(rest, "%d", &data.SwapTotal)
	case "SwapFree":
		fmt.Sscanf(rest, "%d", &data.SwapFree)
	}
}

// GetProcMemInfoJSON - data from "/proc/meminfo"
func GetProcMemInfoJSON() string {
	file, err := os.Open("/proc/meminfo")
	common.LogError(err)
	scanner := bufio.NewScanner(file)
	data := &procMemInfo{}
	data.Type = "proc.meminfo"
	data.Time = currentTime()
	for scanner.Scan() {
		line := scanner.Text()
		parseProcMemInfo(line, data)
	}
	common.LogError(scanner.Err())
	common.LogError(file.Close())

	serialised, err := json.Marshal(data)
	common.LogError(err)
	return string(serialised)
}

// HDD related

type procDiskStats struct {
	Type            string `json:"type"`
	Time            string `json:"time"`
	Device          string
	SizeTotal       uint64
	SizeFree        uint64
	INodesTotal     uint64
	INodesFree      uint64
	IOInProgress    int
	ReadsCompleted  int
	WritesCompleted int
	ReadsTime       int
	WritesTime      int
}

func parseProcDiskStats(line string, data *procDiskStats) {
	var tmp int
	re := regexp.MustCompile("\\s\\s+")
	line = strings.Trim(string(re.ReplaceAll([]byte(line), []byte(" "))), " ")
	parts := strings.SplitN(line, " ", 4)
	switch parts[2] {
	case "sda":
		fallthrough
	case "xvda":
		fmt.Sscanf(line, "%d %d %s %d %d %d %d %d %d %d %d %d %d %d", &tmp, &tmp, &data.Device, &data.ReadsCompleted, &tmp, &tmp, &data.ReadsTime, &data.WritesCompleted, &tmp, &tmp, &data.WritesTime, &data.IOInProgress, &tmp, &tmp)
	}
}

// GetProcDiskStatsJSON - data from /proc/diskstats"
func GetProcDiskStatsJSON() string {
	file, err := os.Open("/proc/diskstats")
	common.LogError(err)
	scanner := bufio.NewScanner(file)
	data := &procDiskStats{}
	data.Type = "proc.diskstats"
	data.Time = currentTime()
	for scanner.Scan() {
		line := scanner.Text()
		parseProcDiskStats(line, data)
	}
	common.LogError(scanner.Err())
	common.LogError(file.Close())

	var stat syscall.Statfs_t
	wd, err := os.Getwd()
	common.LogError(err)
	syscall.Statfs(wd, &stat)
	data.SizeTotal = uint64(stat.Bsize) * stat.Blocks
	data.SizeFree = uint64(stat.Bsize) * stat.Bavail
	data.INodesTotal = stat.Files
	data.INodesFree = stat.Ffree

	serialised, err := json.Marshal(data)
	common.LogError(err)
	return string(serialised)
}

// Storage helpers

// AppendData - store to file
func AppendData(data, file string) {
	f, err := os.OpenFile(file, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	common.LogError(err)
	_, err = f.WriteString(data + "\n")
	common.LogError(err)
	err = f.Close()
	common.LogError(err)
}

// Docker API

type (
	labels struct {
		Service string `json:"com.docker.compose.service"`
	}

	network struct {
		IPAddress string
	}

	container struct {
		ID              string `json:"Id"`
		Image           string
		Names           []string
		State           string
		Created         int64
		HostConfig      map[string]string
		Labels          labels
		NetworkSettings map[string]map[string]network
	}

	cpuUsage struct {
		TotalUsage uint64 `json:"total_usage"`
	}

	cpuStats struct {
		CPUUsage    cpuUsage `json:"cpu_usage"`
		SystemUsage uint64   `json:"system_cpu_usage"`
	}

	memoryStats struct {
		Usage    uint64 `json:"usage"`
		MaxUsage uint64 `json:"max_usage"`
		Limit    uint64 `json:"limit"`
	}

	ioStat struct {
		Operation string `json:"op"`
		Value     int64  `json:"value"`
	}

	networkStat struct {
		ReadBytes    int64 `json:"rx_bytes"`
		ReadPackets  int64 `json:"rx_packets"`
		WriteBytes   int64 `json:"tx_bytes"`
		WritePackets int64 `json:"tx_packets"`
	}

	dockerStats struct {
		CPUStats    cpuStats               `json:"cpu_stats"`
		MemoryStats memoryStats            `json:"memory_stats"`
		IoStats     map[string][]ioStat    `json:"blkio_stats"`
		Networks    map[string]networkStat `json:"networks"`
	}

	inspectedState struct {
		Pid          int
		FinishedAt   string
		LogPath      string
		RestartCount int
	}

	inspectedContainer struct {
		State        inspectedState
		RestartCount int
	}

	aggregatedContainer struct {
		Type         string `json:"type"`
		Time         string `json:"time"`
		Container    string
		ContainerID  string `json:"ContainerId"`
		Image        string
		Created      string
		Finished     string
		RestartCount int
		MainPid      int
		Pids         []int
		State        string
		IPAddress    string

		CPUUser   int `json:"CpuUser"`
		CPUSystem int `json:"CpuSystem"`

		MemoryResident int64
		MemoryCache    int64

		ReadsCompleted  int64
		WritesCompleted int64

		ProcessStatuses string
	}
)

// NetReadBytes - helper for ReadBytes
func (d dockerStats) NetReadBytes() int64 {
	for _, data := range d.Networks {
		return data.ReadBytes
	}
	return 0
}

// NetReadBytes - helper for WriteBytes
func (d dockerStats) NetWriteBytes() int64 {
	for _, data := range d.Networks {
		return data.WriteBytes
	}
	return 0
}

// IoReadBytes - helper for IoReadBytes
func (d dockerStats) IoReadBytes() int64 {
	var sum int64
	for _, data := range d.IoStats {
		for _, stat := range data {
			if stat.Operation == "Read" {
				sum += stat.Value
			}
		}
	}
	return sum
}

// IoWriteBytes - helper for IoWriteBytes
func (d dockerStats) IoWriteBytes() int64 {
	var sum int64
	for _, data := range d.IoStats {
		for _, stat := range data {
			if stat.Operation == "Write" {
				sum += stat.Value
			}
		}
	}
	return sum
}

func (c container) getIPAddress() string {
	for _, networkSettings := range c.NetworkSettings {
		for _, network := range networkSettings {
			return network.IPAddress
		}
	}
	return ""
}

func dockerSocketDial(_, _ string) (net.Conn, error) {
	return net.Dial("unix", "/var/run/docker.sock")
}

func connectToDocker(path string) io.ReadCloser {
	tr := &http.Transport{
		Dial: dockerSocketDial,
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get("http://127.0.0.1/" + path)
	common.LogError(err)
	return resp.Body
}

func parseDockerContainersJSON(reader io.ReadCloser) []aggregatedContainer {
	// Reading data from Docker API
	data, err := ioutil.ReadAll(reader)
	common.LogError(err)

	// Unmarshaling
	result := []container{}
	err = json.Unmarshal([]byte(data), &result)
	common.LogError(err)

	// Converting to stats structures
	stats := []aggregatedContainer{}
	for _, container := range result {
		stat := aggregatedContainer{}
		stat.Type = "docker.stats"
		stat.Time = currentTime()
		stat.Container = container.Labels.Service
		stat.ContainerID = container.ID
		stat.Image = container.Image
		stat.Created = time.Unix(container.Created, 0).Format(time.RFC3339)
		stat.State = container.State
		stat.IPAddress = container.getIPAddress()

		stats = append(stats, stat)
	}

	return stats
}

func includeContainerInspect(container aggregatedContainer) aggregatedContainer {
	reader := connectToDocker("containers/" + container.ContainerID + "/json")
	rawJSON, err := ioutil.ReadAll(reader)
	common.LogError(err)
	common.LogError(reader.Close())
	data := &inspectedContainer{}
	err = json.Unmarshal(rawJSON, data)
	common.LogError(err)
	container.MainPid = data.State.Pid
	container.Finished = data.State.FinishedAt
	container.RestartCount = data.RestartCount

	return container
}

func includeContainerPids(container aggregatedContainer) aggregatedContainer {
	command := "/sys/fs/cgroup/cpuacct/docker/" + container.ContainerID + "/cgroup.procs"
	command2 := strings.Replace(command, "/sys/fs/cgroup/", "/cgroup/", 1)
	lines := readLines(command, command2)
	for _, line := range lines {
		pid, err := strconv.ParseInt(line, 10, 0)
		common.LogError(err)
		container.Pids = append(container.Pids, int(pid))
	}
	return container
}

func includeContainerCPUStats(container aggregatedContainer) aggregatedContainer {
	command := "/sys/fs/cgroup/cpuacct/docker/" + container.ContainerID + "/cpuacct.stat"
	command2 := strings.Replace(command, "/sys/fs/cgroup/", "/cgroup/", 1)
	lines := readLines(command, command2)
	for _, line := range lines {
		var key string
		var value int
		fmt.Sscanf(line, "%s %d", &key, &value)
		switch key {
		case "user":
			container.CPUUser = value
		case "system":
			container.CPUSystem = value
		}
	}
	return container
}

func includeContainerMemoryStats(container aggregatedContainer) aggregatedContainer {
	command := "/sys/fs/cgroup/memory/docker/" + container.ContainerID + "/memory.stat"
	command2 := strings.Replace(command, "/sys/fs/cgroup/", "/cgroup/", 1)
	lines := readLines(command, command2)
	for _, line := range lines {
		var key string
		var value int64
		fmt.Sscanf(line, "%s %d", &key, &value)
		switch key {
		case "total_cache":
			container.MemoryCache = value
		case "total_rss":
			container.MemoryResident = value
		}
	}
	return container
}

func includeContainerIOStats(container aggregatedContainer) aggregatedContainer {
	command := "/sys/fs/cgroup/blkio/docker/" + container.ContainerID + "/blkio.throttle.io_serviced"
	command2 := strings.Replace(command, "/sys/fs/cgroup/", "/cgroup/", 1)
	lines := readLines(command, command2)
	container.ReadsCompleted = 0
	container.WritesCompleted = 0
	for _, line := range lines {
		parts := strings.SplitN(line, " ", 3)
		if len(parts) == 3 {
			value, err := strconv.ParseInt(parts[2], 10, 64)
			common.LogError(err)
			switch parts[1] {
			case "Read":
				container.ReadsCompleted += value
			case "Write":
				container.WritesCompleted += value
			}
		}
	}
	return container
}

func readLines(command string, alternativeCommand string) []string {
	result := []string{}
	file, err := os.Open(command)
	if err != nil {
		file, err = os.Open(alternativeCommand)
	}
	common.LogError(err)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		result = append(result, scanner.Text())
	}
	common.LogError(scanner.Err())
	common.LogError(file.Close())
	return result
}

func includeContainerProcessStatuses(container aggregatedContainer) aggregatedContainer {
	container.ProcessStatuses = ""
	for _, pid := range container.Pids {
		content, err := ioutil.ReadFile("/proc/" + strconv.Itoa(pid) + "/stat")
		common.LogError(err)
		parts := strings.SplitN(string(content), ") ", 2)
		if len(parts) == 2 {
			state := processState(parts[1][0])
			container.ProcessStatuses += " " + state
		}
	}
	container.ProcessStatuses = strings.TrimSpace(container.ProcessStatuses)
	return container
}

func processState(state byte) string {
	switch string(state) {
	case "R":
		return "Running"
	case "S":
		return "Sleeping"
	case "D":
		return "DiskSleeping"
	case "T":
		return "Stopped"
	case "t":
		return "TracingStop"
	case "W":
		return "Paging"
	case "X":
		fallthrough
	case "x":
		return "Dead"
	case "K":
		return "Wakekill"
	case "P":
		return "Parked"
	case "Z":
		return "Zombie"
	default:
		return string(state)
	}
}

func getDockerContainersJSON(elements []aggregatedContainer) string {
	result := ""
	for key, stat := range elements {
		serialised, err := json.Marshal(stat)
		common.LogError(err)
		if key != 0 {
			result = result + "\n"
		}
		result = result + string(serialised)
	}
	return result
}

// GetAllDockerStats - all docker stats
func GetAllDockerStats() string {
	reader := connectToDocker("containers/json")
	containers := parseDockerContainersJSON(reader)
	common.LogError(reader.Close())
	for key, container := range containers {
		container = includeContainerInspect(container)
		container = includeContainerPids(container)
		container = includeContainerCPUStats(container)
		container = includeContainerMemoryStats(container)
		container = includeContainerIOStats(container)
		container = includeContainerProcessStatuses(container)
		containers[key] = container
	}
	containerStats := getDockerContainersJSON(containers)
	return containerStats
}
