package gosshtool

import (
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"net"
)

type PtyInfo struct {
	Term  string
	H     int
	W     int
	Modes ssh.TerminalModes
}

type ReadWriteCloser interface {
	io.Reader
	io.WriteCloser
}

type SSHClientConfig struct {
	Host               string
	User               string
	Password           string
	Privatekey         string
	PrivatekeyPassword string
	DialTimeoutSecond  int
	MaxDataThroughput  uint64
}

func makeConfig(user string, password string, privateKey string, privatekeyPassword string) (config *ssh.ClientConfig) {

	if password == "" && privateKey == "" {
		log.Fatal("No password or private key available")
	}
	config = &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
	}
	if privateKey != "" {
		var err error
		var signer ssh.Signer
		if privatekeyPassword == "" {
			signer, err = ssh.ParsePrivateKey([]byte(privateKey))
		} else {
			signer, err = ssh.ParsePrivateKeyWithPassphrase([]byte(privateKey), []byte(privatekeyPassword))
		}
		if err != nil {
			log.Fatalf("ssh.ParsePrivateKey error:%v", err)
		}
		clientkey := ssh.PublicKeys(signer)
		config = &ssh.ClientConfig{
			User: user,
			Auth: []ssh.AuthMethod{
				clientkey,
			},
			HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
				return nil
			},
		}
	}
	return
}
