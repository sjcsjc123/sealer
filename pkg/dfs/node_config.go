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

import "fmt"

type NodeConfig struct {
	// Dir is the directory of node.
	Dir string
	// Ip is the ip of node.
	IP string
	// Port is the port of node.
	Port int
}

// GetIpList returns the ip list of nodes.
func GetIpList(nodes []NodeConfig) []string {
	var ips []string
	for _, node := range nodes {
		ips = append(ips, node.IP)
	}
	return ips
}

// GetPortList returns the port list of nodes.
func GetPortList(nodes []NodeConfig) []int {
	var ports []int
	for _, node := range nodes {
		ports = append(ports, node.Port)
	}
	return ports
}

// GetDirList returns the dir list of nodes.
func GetDirList(nodes []NodeConfig) []string {
	var dirs []string
	for _, node := range nodes {
		dirs = append(dirs, node.Dir)
	}
	return dirs
}

// GetServerList returns the server list of nodes.
func GetServerList(nodes []NodeConfig) []string {
	var servers []string
	for _, node := range nodes {
		servers = append(servers, fmt.Sprintf("%s:%d", node.IP, node.Port))
	}
	return servers
}
