package weed

import (
	"fmt"
	"strings"
)

type MasterShellConfig struct {
	MDir               string
	IP                 string
	Port               int
	DefaultReplication string
	Peers              []string
	LogFile            string
	BinFile            string
}

type VolumeShellConfig struct {
	MServer []string
	Port    int
	IP      string
	Dir     string
	LogFile string
	BinFile string
}

func generateMasterCmd(conf MasterShellConfig) string {
	return fmt.Sprintf(
		"%s master -mdir=%s -ip=%s -port=%d -defaultReplication=%s -peers=%s >> %s 2>&1",
		conf.BinFile, conf.MDir, conf.IP, conf.Port, conf.DefaultReplication,
		strings.Join(conf.Peers, ","), conf.LogFile)
}

func generateVolumeCmd(conf VolumeShellConfig) string {
	return fmt.Sprintf(
		"%s volume -port=%d -ip=%s -dir=%s -mserver=%s >> %s 2>&1",
		conf.BinFile, conf.Port, conf.IP, conf.Dir, strings.Join(conf.MServer, ","), conf.LogFile)
}
