package main

import (
	"fmt"
	"github.com/Auginte/go-monitoring/domain/parser"
)

func main() {
	// Configuration
	prefix := "../"
	esEndpoint := "http://127.0.0.1:19200/"

	containers := parser.GetContainersToBeMonitored(prefix)
	for containerName, path := range containers {
		fmt.Printf("Exportintg logs of %s\n", containerName)
		jsonFiles := parser.GetJSONLogFiles(path)
		dates := parser.GetUniqueValues(jsonFiles)

		// Creating index
		for _, date := range dates {
			indexName := containerName + "-" + date
			parser.CreateIndex(esEndpoint, indexName)
		}

		// Storing data and mapping for each file
		mapping := parser.ReadMappings(prefix, containerName)
		for file, date := range jsonFiles {
			indexName := containerName + "-" + date
			typeName := containerName + "." + parser.GetEntryTypeFromPath(file)
			urlMapping := esEndpoint + indexName + "/_mapping/" + typeName
			urlStore := esEndpoint + indexName + "/" + typeName

			if len(mapping) > 0 {
				parser.StoreMapping(urlMapping, mapping)
			}
			parser.StoreDataToES(urlStore, file)
		}
	}
}
