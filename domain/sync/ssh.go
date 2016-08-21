package sync

import (
	"github.com/Auginte/go-monitoring/domain/common"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type (
	// SSHClient - wrapper for GoLang ssh connection (connection can have sessions)
	SSHClient struct {
		User           string
		Domain         string
		Port           int
		PrivateKeyPath string
		connection     *ssh.Client
	}

	// SSHSession - sesion, where you can run commands or transfer files
	SSHSession struct {
		session *ssh.Session
	}
)

// MakeSSH - constructor for SshClient
func MakeSSH(url string, privateKey string) *SSHClient {
	parts := strings.SplitN(url, "@", 2)
	return &SSHClient{
		User:           parts[0],
		Domain:         parts[1],
		Port:           22,
		PrivateKeyPath: privateKey,
	}
}

func (s *SSHClient) publicKeyFile(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	common.LogError(err)
	key, err := ssh.ParsePrivateKey(buffer)
	common.LogError(err)
	return ssh.PublicKeys(key)
}

// NewSession - creates xterm session in server; creates connection, if not existed
func (s *SSHClient) NewSession() *SSHSession {
	if s.connection == nil {
		sshConfig := &ssh.ClientConfig{
			User: s.User,
			Auth: []ssh.AuthMethod{
				s.publicKeyFile(s.PrivateKeyPath),
			},
		}
		connection, err := ssh.Dial("tcp", s.Domain+":"+strconv.Itoa(s.Port), sshConfig)
		s.connection = connection
		common.LogError(err)
	}
	session, err := s.connection.NewSession()

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		session.Close()
		common.LogError(err)
	}

	common.LogError(err)
	return &SSHSession{
		session: session,
	}
}

// Close connection for ssh
func (s *SSHClient) Close() {
	if s.connection != nil {
		s.connection.Close()
	}
}

// PipeStdError - all errors in server should go to client error stream
func (s *SSHSession) PipeStdError() {
	stderr, err := s.session.StderrPipe()
	common.LogError(err)
	go io.Copy(os.Stderr, stderr)
}

// Run command on server
func (s *SSHSession) Run(command string) string {
	stdout, err := s.session.StdoutPipe()
	common.LogError(err)
	err = s.session.Run(command)
	data, err := ioutil.ReadAll(stdout)
	common.LogError(err)
	return string(data)
}
