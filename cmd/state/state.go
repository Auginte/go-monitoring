package main

import (
	"fmt"
	"github.com/Auginte/go-monitoring/domain/state"
	"os"
)

func main() {
	dryRun := len(os.Args) > 1 && os.Args[1] == "-d"
	storeLogs := len(os.Args) > 1 && os.Args[1] == "-store"

	directory := "./logs"
	if len(os.Args) > 2 {
		directory = os.Args[2]
		fmt.Println("Using directory: " + directory)
	}

	if dryRun {
		fmt.Println(state.GetProcStatJSON())
		fmt.Println(state.GetProcMemInfoJSON())
		fmt.Println(state.GetProcDiskStatsJSON())
		fmt.Println(state.GetAllDockerStats())
	} else if storeLogs {
		fmt.Println("Storing global stats...")
		os.Mkdir(directory, 0744)
		state.AppendData(state.GetProcStatJSON(), directory+"/proc.stats.json")
		state.AppendData(state.GetProcMemInfoJSON(), directory+"/proc.mem.info.json")
		state.AppendData(state.GetProcDiskStatsJSON(), directory+"/proc.disk.stats.json")
		state.AppendData(state.GetAllDockerStats(), directory+"/docker.fast.stats.json")
		fmt.Println("state stored to: " + directory)
	} else {
		fmt.Println("Simple current state monitoring tool")
		fmt.Println("Usage:")
		fmt.Println("\t./state -d\t\tDry run and shows generated JSONS")
		fmt.Println("\t./state -store\t\tAppends statistics to logs")
		fmt.Println("\t./state -store /var/logs\t\tStore to /var/logs")
	}
}
