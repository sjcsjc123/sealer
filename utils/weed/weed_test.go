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
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

var config1 = &Config{
	Master: []string{"192.168.0.105"},
	Node:   []string{"192.168.0.105"},
	User:   "root",
	Passwd: "FlyDB2023",
	Port:   "22",
}

func TestNewDeployer(t *testing.T) {
	_, err := NewDeployer(config1)
	assert.Nil(t, err)
	fmt.Printf("weedMasterNode: %v \n", config1.weedMasterNode)
	fmt.Printf("weedVolumeNode: %v \n", config1.weedVolumeNode)
}

func TestDownloadWeedToRemoteIp(t *testing.T) {
	_, err := NewDeployer(config1)
	assert.Nil(t, err)
	d := deployer{
		config: config1,
	}
	cmd, err := d.generateDownloadCmd()
	assert.Nil(t, err)
	client := ssh.NewSSHClient(&v1.SSH{
		User:   "root",
		Passwd: "FlyDB2023",
		Port:   "22",
	}, true)
	for k, v := range cmd {
		out, err := client.Cmd(net.ParseIP(k), nil, v)
		assert.Nil(t, err)
		fmt.Printf("out: %v \n", string(out))
	}
}

func TestStart(t *testing.T) {
	d, err := NewDeployer(config1)
	assert.Nil(t, err)
	err = d.Start(context.Background())
	assert.Nil(t, err)
}
