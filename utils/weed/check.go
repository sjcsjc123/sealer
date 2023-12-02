// Copyright © 2021 Alibaba Group Holding Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package weed

import (
	"fmt"
	"github.com/pkg/errors"
	"net"
	"os"
	"os/exec"
	"path"
	"strconv"
)

// checkPort checks if the port is available or can be used.
func checkPort(host string, port int) bool {
	// Check TCP
	tcpAddr := fmt.Sprintf("%s:%d", host, port)
	tcpConn, tcpErr := net.Dial("tcp", tcpAddr)
	if tcpConn != nil {
		defer tcpConn.Close()
	}
	if tcpErr == nil {
		return true // Port is in use
	}
	return false // Port is available
}

// checkBinFile checks if the bin file is exist and can be used.
func checkBinFile(fileName string) bool {
	binName := path.Base(fileName)
	_, err := os.Stat(fileName)
	if err != nil {
		return false
	}
	switch binName {
	case "weed":
		err := exec.Command(fileName, "version").Run()
		if err != nil {
			return false
		}
	case "etcd":
		err := exec.Command(fileName, "--version").Run()
		if err != nil {
			return false
		}
	default:
	}
	return true
}

// check checks if the config is valid.
func check(config *Config) error {
	// check config param valid
	err := checkConfigParam(config)
	if err != nil {
		return err
	}
	// check weed master peers, set it to `weedMasterList`
	err = checkWeedMaster(config)
	if err != nil {
		return err
	}
	err = checkWeedVolume(config)
	if err != nil {
		return err
	}
	return createDir(config)
}

func initPort(currentIP string, port int) int {
	for checkPort(currentIP, port) {
		port++
	}
	return port
}

func checkConfigParam(config *Config) error {
	// check config add set default value if not set
	if config.BinDir == "" {
		config.BinDir = "/var/sealer/bin"
	}
	if config.DefaultReplication == "" {
		config.DefaultReplication = "000"
	}
	if len(config.Master) == 0 {
		return errors.New("master list is empty")
	}
	if len(config.Node) == 0 {
		return errors.New("node list is empty")
	}
	return nil
}

func checkWeedMaster(config *Config) error {
	config.weedMasterNode = make([]NodeConfig, 0)
	if len(config.Master) == 1 {
		// k0s, need more local node
		config.masterNeedMoreLocalNode = true
		weedMasterPort := initPort(config.Master[0], 9333)
		for i := 0; i < 3; i++ {
			config.weedMasterNode = append(config.weedMasterNode, NodeConfig{
				IP:   config.Master[0],
				Port: weedMasterPort,
				Dir:  DefaultWeedMasterDir + "-" + strconv.Itoa(i),
			})
			weedMasterPort = initPort(config.Master[0], weedMasterPort+1)
		}
	} else {
		// k8s
		for _, m := range config.Master {
			weedMasterPort := initPort(m, 9333)
			config.weedMasterNode = append(config.weedMasterNode, NodeConfig{
				IP:   m,
				Port: weedMasterPort,
				Dir:  DefaultWeedMasterDir,
			})
		}
	}
	return nil
}

func checkWeedVolume(config *Config) error {
	config.weedVolumeNode = make([]NodeConfig, 0)
	if len(config.Node) == 1 {
		// k0s, need more local node
		config.volumeNeedMoreLocalNode = true
		weedVolumePort := initPort(config.Node[0], 8080)
		for i := 0; i < 3; i++ {
			config.weedVolumeNode = append(config.weedVolumeNode, NodeConfig{
				IP:   config.Node[0],
				Port: weedVolumePort,
				Dir:  DefaultWeedVolumeDir + "-" + strconv.Itoa(i),
			})
			weedVolumePort = initPort(config.Node[0], weedVolumePort+1)
		}
	} else {
		// k8s
		for _, n := range config.Node {
			weedVolumePort := initPort(n, 8080)
			config.weedVolumeNode = append(config.weedVolumeNode, NodeConfig{
				IP:   n,
				Port: weedVolumePort,
				Dir:  DefaultWeedVolumeDir,
			})
		}
	}
	return nil
}

func createDir(config *Config) error {
	for _, m := range config.weedMasterNode {
		_, err := os.Stat(m.Dir)
		if err != nil {
			if os.IsNotExist(err) {
				err = os.MkdirAll(m.Dir, 0755)
				if err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}
	for _, v := range config.weedVolumeNode {
		_, err := os.Stat(v.Dir)
		if err != nil {
			if os.IsNotExist(err) {
				err = os.MkdirAll(v.Dir, 0755)
				if err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}
	_, err := os.Stat(config.BinDir)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(config.BinDir, 0755)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}
