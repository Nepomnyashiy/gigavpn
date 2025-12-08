package sshm

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"net"
	"time"
)

// SSHClient представляет собой клиент для взаимодействия с удаленным сервером по SSH.
type SSHClient struct {
	config *ssh.ClientConfig
	host   string
	port   string
}

// NewSSHClient создает новый экземпляр SSHClient.
// host: IP-адрес или доменное имя сервера.
// port: SSH-порт.
// user: Имя пользователя для подключения.
// privateKeyPath: Путь к файлу приватного ключа.
func NewSSHClient(host, port, user, privateKeyPath string) (*SSHClient, error) {
	key, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать приватный ключ: %w", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("не удалось распарсить приватный ключ: %w", err)
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // ВНИМАНИЕ: небезопасно для продакшена!
		Timeout:         5 * time.Second,
	}

	return &SSHClient{
		config: config,
		host:   host,
		port:   port,
	}, nil
}

// RunCommand выполняет команду на удаленном сервере и возвращает ее вывод.
func (c *SSHClient) RunCommand(cmd string) (string, error) {
	address := net.JoinHostPort(c.host, c.port)
	conn, err := ssh.Dial("tcp", address, c.config)
	if err != nil {
		return "", fmt.Errorf("не удалось подключиться к серверу: %w", err)
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		return "", fmt.Errorf("не удалось создать сессию: %w", err)
	}
	defer session.Close()

	output, err := session.CombinedOutput(cmd)
	if err != nil {
		return string(output), fmt.Errorf("не удалось выполнить команду: %w", err)
	}

	return string(output), nil
}
