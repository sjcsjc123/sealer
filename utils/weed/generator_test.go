package weed

import "testing"

func TestGenerateMasterCmd(t *testing.T) {
	cmd := generateMasterCmd(MasterShellConfig{
		MDir:               "/home/m1",
		IP:                 "127.0.0.1",
		Port:               9333,
		Peers:              []string{"127.0.0.1:9333"},
		DefaultReplication: "000",
		LogFile:            "/home/m1/master.log",
		BinFile:            "/home/weed",
	})
	t.Log(cmd)
}
