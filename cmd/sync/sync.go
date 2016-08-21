package main

import (
	"fmt"
	"github.com/Auginte/go-monitoring/domain/sync"
	"os"
)

func main() {
	if len(os.Args) >= 3 {
		ssh := sync.MakeSSH(os.Args[1], os.Args[2])
		session := ssh.NewSession()
		session.PipeStdError()
		data := session.Run("ps aux")
		fmt.Printf("|%s|\n", data)
		ssh.Close()
	} else {
		fmt.Println("Tool to sync logs from LIVE to DEV machine")
		fmt.Println("Usage:")
		fmt.Println("\tsync user@example.com /home/user/.ssh/private.pem")
	}
}
