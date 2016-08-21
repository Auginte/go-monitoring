package parser

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
	"github.com/Auginte/go-monitoring/domain/common"
)



func httpPost(urlPath, body string) string {
	client := http.Client{}
	request, err := http.NewRequest("POST", urlPath, bytes.NewBuffer([]byte(body)))
	common.LogError(err)
	response, err := client.Do(request)
	common.LogError(err)
	responseText, err := ioutil.ReadAll(response.Body)
	common.LogError(err)
	return string(responseText)
}

func httpPut(urlPath, body string) string {
	client := http.Client{}
	request, err := http.NewRequest("PUT", urlPath, bytes.NewBuffer([]byte(body)))
	common.LogError(err)
	response, err := client.Do(request)
	common.LogError(err)
	responseText, err := ioutil.ReadAll(response.Body)
	common.LogError(err)
	return string(responseText)
}

// ReadMappings - reads mapping from json file
func ReadMappings(prefix, container string) string {
	path := prefix + "/" + container + "/config/log-mapping.json"
	if _, err := os.Stat(path); err == nil {
		mapping, err := ioutil.ReadFile(path)
		common.LogError(err)
		return string(mapping)
	}
	return ""
}

// CreateIndex - post to ElasticSearch
func CreateIndex(esEndpoint, indexName string) {
	log.Println("Creating index: " + indexName)
	response := httpPost(esEndpoint+indexName, "")
	log.Println("\tResponse: " + response)
}

func fixInconsistentData(line string) string {
	result := strings.Replace(line, `"upstream_response_time":"-"`, `"upstream_response_time":"-1"`, -1)
	return result
}

// StoreMapping - post to ElasticSearch
func StoreMapping(urlMapping, mapping string) {
	log.Println("Storing mapping: " + urlMapping)
	response := httpPut(urlMapping, string(mapping))
	log.Println("\tResponse: " + response)
}

// StoreDataToES - post data to ElasticSearch
func StoreDataToES(urlStore, dataFile string) {
	// Reading data
	file, err := os.Open(dataFile)
	common.LogError(err)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// Storing data (using hash to not import same data)
		if len(strings.Trim(scanner.Text(), "\t ")) > 0 {
			log.Println("Storing data: " + urlStore)
			id := fmt.Sprintf("%x", md5.Sum(scanner.Bytes()))
			data := fixInconsistentData(scanner.Text())
			response := httpPost(urlStore+"/"+id, data)
			log.Println("\tResponse: " + response)
		}
	}
	common.LogError(scanner.Err())
	common.LogError(file.Close())
}

// GetContainersToBeMonitored - returns [container_name => path]
func GetContainersToBeMonitored(prefix string) map[string]string {
	topDirectories, err := ioutil.ReadDir(prefix)
	common.LogError(err)
	monitoredContainers := map[string]string{}
	for _, directory := range topDirectories {
		if directory.IsDir() {
			subDirs, err := ioutil.ReadDir(prefix + "/" + directory.Name())
			common.LogError(err)
			hasLogs := false
			hasConfig := false
			for _, subDir := range subDirs {
				if subDir.IsDir() && subDir.Name() == "logs" {
					hasLogs = true
				}
				if subDir.IsDir() && subDir.Name() == "config" {
					hasConfig = true
				}
			}
			if hasLogs && hasConfig {
				monitoredContainers[directory.Name()] = prefix + "/" + directory.Name()
			}
		}
	}
	return monitoredContainers
}

// GetJSONLogFiles - returns [path => date]
func GetJSONLogFiles(path string) map[string]string {
	fullPath := path + "/logs"
	files, err := ioutil.ReadDir(fullPath)
	common.LogError(err)
	jsonFiles := map[string]string{}
	date := time.Now().Format("2006.01.02")
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".json") && !file.IsDir() {
			jsonFiles[fullPath+"/"+file.Name()] = date
		}
	}
	return jsonFiles
}

// GetUniqueValues - helper for unique values
func GetUniqueValues(data map[string]string) []string {
	unique := []string{}
	for _, value := range data {
		had := false
		for _, old := range unique {
			if old == value {
				had = true
				break
			}
		}
		if !had {
			unique = append(unique, value)
		}
	}
	return unique
}

// GetEntryTypeFromPath - file name from whole path
func GetEntryTypeFromPath(path string) string {
	re := regexp.MustCompile("^.+\\/logs\\/(.+)\\.json$")
	result := re.FindStringSubmatch(path)
	return result[1]
}
