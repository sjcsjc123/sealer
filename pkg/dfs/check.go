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

package dfs

import (
	"fmt"
	"github.com/pkg/errors"
	"net"
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

// getCheckCommandScript check lsof command is installed.
func getCheckCommandScript() string {
	return fmt.Sprintf(`
if ! command -v lsof &> /dev/null
then
	echo "lsof command could not be found"
	exit
fi
`)
}

// check checks if the config is valid.
func check(config *Config) error {
	// check config param valid
	err := checkConfigParam(config)
	if err != nil {
		return err
	}
	// set minioNode
	return setMinioNode(config)
}

// setMinioNode sets minio node.
func setMinioNode(config *Config) error {
	minioNodes := make([]NodeConfig, 0)
	ipMap := make(map[string]bool)
	for _, m := range config.Master {
		if _, ok := ipMap[m]; !ok {
			ipMap[m] = true
			minioNodes = append(minioNodes, NodeConfig{
				IP:   m,
				Port: 9000,
				Dir:  config.MinioDir,
			})
		}
	}
	for _, n := range config.Node {
		if _, ok := ipMap[n]; !ok {
			ipMap[n] = true
			minioNodes = append(minioNodes, NodeConfig{
				IP:   n,
				Port: 9000,
				Dir:  config.MinioDir,
			})
		}
	}
	config.minioNode = minioNodes
	return nil
}

//func initPort(currentIP string, port int) int {
//	for checkPort(currentIP, port) {
//		if checkPortIsMinio(currentIP, port) {
//			return port
//		}
//		port++
//	}
//	return port
//}

// notice：
//
//	if port is minio, the start method of deployer maybe will print error msg like this:
//	"port xxxx is in use"
//func checkPortIsMinio(currentIP string, port int) bool {
//	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", currentIP, port), time.Second)
//	if err != nil {
//		return false
//	}
//	defer conn.Close()
//	_, err = conn.Write([]byte("GET /minio/health/live HTTP/1.1\r\nHost: localhost\r\n\r\n"))
//	if err != nil {
//		return false
//	}
//
//	buf := make([]byte, 512)
//	_, err = conn.Read(buf)
//	if err != nil {
//		return false
//	}
//
//	response := string(buf)
//	return containsMinioInfo(response)
//}
//
//func containsMinioInfo(response string) bool {
//	response = strings.ToLower(response)
//	return strings.Contains(response, "minio") && strings.Contains(response, "200 ok")
//}

func checkConfigParam(config *Config) error {
	// check config add set default value if not set
	if config.BinDir == "" {
		config.BinDir = "/var/lib/sealer/bin"
	}
	if config.MinioDir == "" {
		config.MinioDir = "/var/lib/sealer/minio"
	}
	if config.User == "" {
		config.User = "root"
	}
	if config.Port == "" {
		config.Port = "22"
	}
	if len(config.Master) == 0 {
		return errors.New("master list is empty")
	}
	if len(config.Node) == 0 {
		return errors.New("node list is empty")
	}
	if config.MinioUser == "" {
		config.MinioUser = "sealerminio"
	}
	if config.MinioPasswd == "" {
		config.MinioPasswd = "sealerminio"
	}
	if len(config.MinioUser) < 8 {
		return errors.New("minio user length must be greater than 8")
	}
	if len(config.MinioPasswd) < 8 {
		return errors.New("minio password length must be greater than 8")
	}
	return nil
}
