package weed

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
