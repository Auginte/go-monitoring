package sync

import (
	"fmt"
	"github.com/Auginte/go-monitoring/domain/common"
	"github.com/sfreiberg/simplessh"
	"strings"
)

// GetContainerJSONFiles - returns list of Json log files inside each container folder
func GetContainerJSONFiles(client *simplessh.Client, prefix string) []string {
	command := fmt.Sprintf("cd %s && find */logs/*.json", prefix)
	lines, err := client.Exec(command)
	common.LogError(err)
	files := strings.Split(strings.TrimSpace(string(lines)), "\n")
	for key, file := range files {
		files[key] = prefix + "/" + file
	}
	return files
}

// GetStatsJSONFiles - returns list of files for current state monitoring
func GetStatsJSONFiles(client *simplessh.Client, prefix string) []string {
	command := fmt.Sprintf("find %s/*.json", prefix)
	lines, err := client.Exec(command)
	common.LogError(err)
	return strings.Split(strings.TrimSpace(string(lines)), "\n")
}

// DownloadFile - escapes file name, downloads and returns target file
func DownloadFile(client *simplessh.Client, file string, targetPath string) string {
	name := strings.Replace(strings.Trim(file, "/"), "/", "--", -1)
	path := targetPath + "/" + name
	err := client.Download(file, path)
	common.LogError(err)
	return path
}
