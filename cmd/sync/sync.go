package main

import (
	"fmt"
	"github.com/Auginte/go-monitoring/domain/sync"
	"os"
)

func main() {
	if len(os.Args) >= 6 {
		fmt.Println("Connecting to Server...")
		ssh := sync.MakeSSH(os.Args[1], os.Args[2])
		containersRootPath := os.Args[3]
		stateLogsPath := os.Args[4]
		localBackupPath := os.Args[5]

		client := ssh.Client()
		fmt.Println("Getting list of JSON files...")
		containerLogFiles := sync.GetContainerJSONFiles(client, containersRootPath)
		statsLogFiles := sync.GetStatsJSONFiles(client, stateLogsPath)
		files := append(containerLogFiles, statsLogFiles...)
		for _, file := range files {
			fmt.Printf("Will download: %s\n", file)
		}

		for i, file := range files {
			fmt.Printf("Downloading %d/%d: %s...\n", i+1, len(files), file)
			downloaded := sync.DownloadFile(client, file, localBackupPath)
			fmt.Printf("\tDownloaded to: %s\n", downloaded)
		}
		client.Close()
		fmt.Println("Logs downloaded")
	} else {
		fmt.Println("Tool to sync logs from LIVE to DEV machine")
		fmt.Println("Usage:")
		fmt.Println("\tsync user@example.com /home/user/.ssh/private.pem /docker /var/log/monitoring/state /backup/logs")
		fmt.Println("\tsync USER@SERVER_HOST PRIVATE_KEY CONTAINERS_ROOT_PATH STATE_LOGS_PATH LOCAL_BACKUP_PATH")
	}
}
