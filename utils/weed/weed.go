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
	"context"
	"fmt"
	v1 "github.com/sealerio/sealer/types/api/v1"
	"github.com/sealerio/sealer/utils/ssh"
	"net"
	"os"
	"path/filepath"
	"strconv"
)

const (
	SeaWeedFsGitHubOrg   = "seaweedfs"
	SeaWeedFsGithubRepo  = "seaweedfs"
	SeaWeedFsVersion     = "3.54"
	GOOSLinux            = "linux"
	DefaultWeedMasterDir = "/var/sealer/weed/master"
	DefaultWeedVolumeDir = "/var/sealer/weed/volume"
)

type Config struct {
	// Master is the master node config.
	Master []string
	// Node is the slave node.
	Node []string
	// BinDir is the directory of binary.
	BinDir string
	// DefaultReplication is the default replication.
	DefaultReplication string
	// Encrypted is the flag of encrypted.
	Encrypted bool
	// User is the user of ssh.
	User string
	// Passwd is the password of ssh.
	Passwd string
	// Pk is the private key of ssh.
	Pk string
	// PkPasswd is the password of private key.
	PkPasswd string
	// Port is the port of ssh.
	Port string
	// masterNeedMoreLocalNode is the flag of master node need more local node.
	masterNeedMoreLocalNode bool
	// volumeNeedMoreLocalNode is the flag of volume node need more local node.
	volumeNeedMoreLocalNode bool
	// weedMasterNode is the config of master node.
	weedMasterNode []NodeConfig
	// weedVolumeNode is the config of volume node.
	weedVolumeNode []NodeConfig
}

type Deployer interface {
	// Start starts deployer cluster.
	Start(ctx context.Context) error
	// Stop stops deployer cluster.
	Stop(ctx context.Context) error
	// IsRunning returns the status of deployer cluster.
	IsRunning(ctx context.Context) bool
	// UploadFile uploads file to deployer cluster.
	UploadFile(ctx context.Context, filename string) error
	// UploadDir uploads directory to deployer cluster.
	UploadDir(ctx context.Context, dir string) error
	// DownloadFile downloads file from deployer cluster.
	DownloadFile(ctx context.Context, fid string, outputDir string) error
}

type deployer struct {
	master Master
	volume Volume
	config *Config
	ssh    ssh.Interface
}

type IpPort struct {
	IP   string
	Port int
}

func NewDeployer(config *Config) (Deployer, error) {
	err := check(config)
	if err != nil {
		return nil, err
	}
	return &deployer{
		config: config,
		master: NewWeedMaster(config),
		volume: NewWeedVolume(config, config.weedMasterNode),
		ssh: ssh.NewSSHClient(&v1.SSH{
			User:      config.User,
			Passwd:    config.Passwd,
			Port:      config.Port,
			Pk:        config.Pk,
			PkPasswd:  config.PkPasswd,
			Encrypted: config.Encrypted,
		}, true),
	}, nil
}

func (d *deployer) Start(ctx context.Context) error {
	downloadCmd, err := d.generateDownloadCmd()
	if err != nil {
		return err
	}
	// download weed
	for k, v := range downloadCmd {
		err = d.ssh.CmdAsync(net.ParseIP(k), nil, v)
		if err != nil {
			return err
		}
	}
	startMasterCmd := d.generateStartMasterCmd()
	for k, v := range startMasterCmd {
		ip := net.ParseIP(k.IP)
		cmd := v
		go func() {
			err = d.ssh.CmdAsync(ip, nil, cmd)
			if err != nil {
				fmt.Printf("err: %v \n", err)
				os.Exit(1)
			}
		}()
		for {
			_, err = d.ssh.Cmd(ip, nil, "lsof -i:"+strconv.Itoa(k.Port))
			if err != nil {
				continue
			}
			break
		}
	}
	startVolumeCmd := d.generateStartVolumeCmd()
	for k, v := range startVolumeCmd {
		ip := net.ParseIP(k.IP)
		cmd := v
		go func() {
			err = d.ssh.CmdAsync(ip, nil, cmd)
			if err != nil {
				fmt.Printf("err: %v \n", err)
				os.Exit(1)
			}
		}()
		for {
			_, err = d.ssh.Cmd(ip, nil, "lsof -i:"+strconv.Itoa(k.Port))
			if err != nil {
				continue
			}
			break
		}
	}
	return nil
}

func (d *deployer) Stop(ctx context.Context) error {
	return nil
}

func (d *deployer) IsRunning(ctx context.Context) bool {
	return false
}

func (d *deployer) UploadFile(ctx context.Context, filename string) error {
	return nil
}

func (d *deployer) UploadDir(ctx context.Context, dir string) error {
	return nil
}

func (d *deployer) DownloadFile(ctx context.Context, fid string, outputDir string) error {
	return nil
}

func (d *deployer) generateStartMasterCmd() map[IpPort]string {
	res := make(map[IpPort]string)
	for _, m := range d.config.weedMasterNode {
		script := generateMasterCmd(MasterShellConfig{
			MDir:               m.Dir,
			IP:                 m.IP,
			Port:               m.Port,
			Peers:              GetServerList(d.config.weedMasterNode),
			DefaultReplication: d.config.DefaultReplication,
			LogFile:            m.Dir + "/master.log",
			BinFile:            d.config.BinDir + "/weed",
		})
		res[IpPort{
			IP:   m.IP,
			Port: m.Port,
		}] = script
	}
	return res
}

func (d *deployer) generateStartVolumeCmd() map[IpPort]string {
	res := make(map[IpPort]string)
	for _, v := range d.config.weedVolumeNode {
		script := generateVolumeCmd(VolumeShellConfig{
			MServer: GetServerList(d.config.weedMasterNode),
			IP:      v.IP,
			Port:    v.Port,
			Dir:     v.Dir,
			LogFile: v.Dir + "/volume.log",
			BinFile: d.config.BinDir + "/weed",
		})
		res[IpPort{
			IP:   v.IP,
			Port: v.Port,
		}] = script
	}
	return res
}

func (d *deployer) generateDownloadCmd() (map[string]string, error) {
	res := make(map[string]string)
	url, err := weedDownloadURL()
	if err != nil {
		return nil, err
	}
	tarFileName := filepath.Join(d.config.BinDir, "weed.tar.gz")
	cmd := fmt.Sprintf(`
        binDir=%s
        weedCmd=%s

        if [ -d $binDir ] && [ -x $weedCmd ]; then
            echo "weed command is present and executable. No need to download."
        else
            url=%s
            tarFileName=%s
            curl -L -o $tarFileName $url && tar -xvf $tarFileName -C %s && rm -f $tarFileName
        fi
    `, d.config.BinDir, filepath.Join(d.config.BinDir, "weed"), url, tarFileName, d.config.BinDir)
	for _, node := range d.config.weedMasterNode {
		res[node.IP] = cmd
	}
	for _, node := range d.config.weedVolumeNode {
		res[node.IP] = cmd
	}
	return res, nil
}
