package main

import (
	"fmt"
	"github.com/Auginte/go-monitoring/domain/common"
	"gopkg.in/hypersleep/easyssh.v0"
	"os"
	"strings"
)

func testShh(ssh *easyssh.MakeConfig) {
	response, err := ssh.Run("ps ax")
	common.LogError(err)
	fmt.Println(response)
}

func main() {
	if len(os.Args) >= 3 {
		parts := strings.SplitN(os.Args[1], "@", 2)
		ssh := &easyssh.MakeConfig{
			User:   parts[0],
			Server: parts[1],
			Key:    os.Args[2],
			Port:   "22",
		}
		testShh(ssh)
	} else {
		fmt.Println("Tool to sync logs from LIVE to DEV machine")
		fmt.Println("Usage:")
		fmt.Println("\tsync user@example.com /.ssh/private.pem")
	}
}
