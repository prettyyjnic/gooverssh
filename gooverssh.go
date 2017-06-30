package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"strings"

	"github.com/BurntSushi/toml"
	"github.com/scottkiss/gosshtool"
)

var configs tomlConfig
var runningTurnel map[string]interface{}
var waitgroup sync.WaitGroup
var m *sync.RWMutex

func init() {
	configs.SSH = make(map[string]ssh)
	m = new(sync.RWMutex)
	runningTurnel = make(map[string]interface{})
}

func loadconfig() {
	var configFileName string
	if len(os.Args) <= 1 { // 如果没有传入参数，则默认获取 ./conf下所有的 .toml 文件
		configFileName = "./conf"
	} else {
		configFileName = os.Args[1]
	}
	f, err := os.Stat(configFileName)
	if err != nil {
		panic(err)
	}
	if f.IsDir() { // 如果是文件夹则读取所有 toml结尾的文件
		list, err := filepath.Glob(strings.TrimRight(strings.Replace(configFileName, "\\", "/", -1), "/") + "/*.toml")
		if err != nil {
			panic(err)
		}
		for _, v := range list {
			loadconfigFromFile(v)
		}
	} else {
		loadconfigFromFile(configFileName)
	}
}

func loadconfigFromFile(configFileName string) {
	var config tomlConfig

	if _, err := toml.DecodeFile(configFileName, &config); err != nil {
		panic("load config error" + err.Error())
	}
	for k2, v2 := range config.SSH {
		configs.SSH[k2] = v2
	}
}

func main() {

	loadconfig()

	for serverName, sshConfig := range configs.SSH {
		waitgroup.Add(1)
		go func(serverName string, sshConfig ssh) {
			if sshConfig.Require != "" {
				waitingTimes := 0
				for {
					m.RLock()
					_, ok := runningTurnel[sshConfig.Require]
					m.RUnlock()
					if !ok {
						log.Println(serverName, "waiting for ", sshConfig.Require)
						<-time.After(time.Second * 1)
						waitingTimes++
					} else {
						break
					}
					if waitingTimes > 20 {
						panic(sshConfig.Require + "启动超时！")
					}
				}
			}
			log.Println(serverName, "init... ")
			server := new(gosshtool.LocalForwardServer)
			server.LocalBindAddress = ":" + sshConfig.LocalPort
			server.RemoteAddress = sshConfig.RemoteAddress
			server.SshServerAddress = sshConfig.SshServerAddress
			server.SshUserName = sshConfig.SshUser
			server.SshUserPassword = sshConfig.SshPassword
			server.SshPrivateKeyPassword = sshConfig.SshPrivateKeyPassword

			buf, err := ioutil.ReadFile(sshConfig.SshPriveteKeyPath)
			if err != nil {
				panic(err)
			}
			server.SshPrivateKey = string(buf)

			log.Println(serverName, "starting... ")
			server.Start(func() {
				m.Lock()
				runningTurnel[serverName] = 1
				m.Unlock()
				log.Println(serverName, "is started ! ", sshConfig.LocalPort, " -> ", sshConfig.RemoteAddress)
			})
			server.Stop()
			waitgroup.Done()
		}(serverName, sshConfig)
	}
	go func() {
		for {
			m.RLock()
			if len(runningTurnel) == len(configs.SSH) {
				log.Println("all ssh turnel had started...")
				m.RUnlock()
				break
			}
			m.RUnlock()
			<-time.After(time.Second * 1)
		}
	}()
	waitgroup.Wait()
}

type tomlConfig struct {
	SSH map[string]ssh
}

type ssh struct {
	RemoteAddress         string `toml:"remoteAddress"`
	LocalPort             string `toml:"localPort"`
	SshServerAddress      string `toml:"sshServerAddress"`
	SshUser               string `toml:"sshUser"`
	SshPassword           string `toml:"sshPassword"`
	SshPriveteKeyPath     string `toml:"sshPriveteKeyPath"`
	SshPrivateKeyPassword string `toml:"sshPrivateKeyPassword"`
	Require               string `toml:"require"`
}
