package sync

import (
	"github.com/Auginte/go-monitoring/domain/common"
	"github.com/sfreiberg/simplessh"
	"strconv"
	"strings"
)

type (
	// SSHConfig - data transfer object for SSH credentials
	SSHConfig struct {
		User           string
		Domain         string
		Port           int
		PrivateKeyPath string
	}
)

// MakeSSH - constructor for SSHConfig
func MakeSSH(url string, privateKey string) *SSHConfig {
	parts := strings.SplitN(url, "@", 2)
	return &SSHConfig{
		User:           parts[0],
		Domain:         parts[1],
		Port:           22,
		PrivateKeyPath: privateKey,
	}
}

// Client of simplessh
func (s *SSHConfig) Client() *simplessh.Client {
	result, err := simplessh.ConnectWithKeyFile(s.Domain+":"+strconv.Itoa(s.Port), s.User, s.PrivateKeyPath)
	common.LogError(err)
	return result
}
