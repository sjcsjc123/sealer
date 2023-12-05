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
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/pkg/errors"
	v1 "github.com/sealerio/sealer/types/api/v1"
	"github.com/sealerio/sealer/utils/ssh"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type Config struct {
	// Master is the master node config.
	Master []string
	// Node is the slave node.
	Node []string
	// BinDir is the directory of binary.
	BinDir string
	// MinioDir is the directory of minio.
	MinioDir string
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
	// MinioUser is the user of minio.
	MinioUser string
	// MinioPasswd is the password of minio.
	MinioPasswd string
	// minioNode is the config of master node.
	minioNode []NodeConfig
}

type Deployer interface {
	// Start starts deployer cluster.
	Start(ctx context.Context) error
	// Stop stops deployer cluster.
	Stop(ctx context.Context) error
	// IsRunning returns the status of deployer cluster.
	// Notice: IsRunning method could cost a lot of time, suggest set context timeout.
	// like this: ctx, _ := context.WithTimeout(context.Background(), time.Second*1)
	IsRunning(ctx context.Context) bool
	// UploadFile uploads file to deployer cluster.
	UploadFile(ctx context.Context, filename string) error
	// UploadDir uploads directory to deployer cluster.
	UploadDir(ctx context.Context, dir string) error
	// DownloadFile downloads file from deployer cluster.
	DownloadFile(ctx context.Context, prefix string, outPutDir string) (int, error)
	// ListDir lists directory from deployer cluster.
	ListDir(ctx context.Context, dir string) ([]string, error)
	// RemoveFile removes file from deployer cluster.
	RemoveFile(ctx context.Context, filename string) error
	// RemoveDir removes directory from deployer cluster.
	RemoveDir(ctx context.Context, dir string) error
}

type deployer struct {
	config *Config
	ssh    ssh.Interface
	client []*minio.Client
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
	for _, node := range d.config.minioNode {
		// check port is in use
		if checkPort(node.IP, node.Port) {
			return errors.Errorf("port %d of %s is in use", node.Port, node.IP)
		}
		// check command
		checkCommandScript := getCheckCommandScript()
		err := d.ssh.CmdAsync(net.ParseIP(node.IP), nil, checkCommandScript)
		if err != nil {
			return err
		}
	}
	downloadCmd, err := d.generateDownloadCmd()
	if err != nil {
		return err
	}
	// download minio
	for k, v := range downloadCmd {
		err = d.ssh.CmdAsync(net.ParseIP(k), nil, v)
		if err != nil {
			return err
		}
	}
	// start minio
	startCmd, err := d.generateStartCmd()
	if err != nil {
		return err
	}
	for k, v := range startCmd {
		ip := net.ParseIP(k)
		cmd := v
		go func() {
			err = d.ssh.CmdAsync(ip, map[string]string{
				"MINIO_ROOT_USER":     d.config.MinioUser,
				"MINIO_ROOT_PASSWORD": d.config.MinioPasswd,
			}, cmd)
			if err != nil {
				return
			}
		}()
	}
	for {
		ok := d.IsRunning(ctx)
		if ok {
			break
		}
	}
	return d.setMinioClient()
}

func (d *deployer) setMinioClient() error {
	endpoints := GetServerList(d.config.minioNode)
	clientList := make([]*minio.Client, 0)
	for _, endpoint := range endpoints {
		client, err := minio.New(endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(d.config.MinioUser, d.config.MinioPasswd, ""),
			Secure: false,
		})
		if err != nil {
			return err
		}
		clientList = append(clientList, client)
	}
	d.client = clientList
	return nil
}

func (d *deployer) Stop(ctx context.Context) error {
	for _, node := range d.config.minioNode {
		ip := net.ParseIP(node.IP)
		go func() {
			// kill application which is using port
			_ = d.ssh.CmdAsync(ip, nil,
				"kill -15 $(lsof -i:"+fmt.Sprintf("%d", 9000)+" | grep -v PID | awk '{print $2}')")
		}()
	}
	for {
		ok := d.IsRunning(ctx)
		if !ok {
			break
		}
	}
	return nil
}

func (d *deployer) IsRunning(ctx context.Context) bool {
	if d.client == nil {
		err := d.setMinioClient()
		if err != nil || d.client == nil {
			return false
		}
	}
	for _, client := range d.client {
		ctx, _ = context.WithTimeout(ctx, time.Second*1)
		_, err := client.ListBuckets(ctx)
		if err != nil {
			continue
		}
		return true
	}
	return false
}

func (d *deployer) UploadFile(ctx context.Context, filename string) error {
	if d.client == nil {
		err := d.setMinioClient()
		if err != nil || d.client == nil {
			return err
		}
	}
	for _, client := range d.client {
		exists, err := client.BucketExists(ctx, "sealer")
		if err != nil {
			continue
		}
		if !exists {
			err = client.MakeBucket(ctx, "sealer", minio.MakeBucketOptions{})
			if err != nil {
				continue
			}
		}
		filename, err = filepath.Abs(filename)
		if err != nil {
			return err
		}
		// upload file to bucket
		_, err = client.FPutObject(ctx, "sealer", filename, filename, minio.PutObjectOptions{})
		if err != nil {
			continue
		}
		return nil
	}
	return errors.New("there is no minio node can work")
}

func (d *deployer) UploadDir(ctx context.Context, dir string) error {
	if d.client == nil {
		err := d.setMinioClient()
		if err != nil || d.client == nil {
			return err
		}
	}
	for _, client := range d.client {
		exists, err := client.BucketExists(ctx, "sealer")
		if err != nil {
			continue
		}
		if !exists {
			err = client.MakeBucket(ctx, "sealer", minio.MakeBucketOptions{})
			if err != nil {
				continue
			}
		}
		dir, err = filepath.Abs(dir)
		if err != nil {
			return err
		}
		// upload dir to bucket
		err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				_, err = client.FPutObject(ctx, "sealer", path, path, minio.PutObjectOptions{})
				if err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			continue
		}
		return nil
	}
	return errors.New("there is no minio node can work")
}

func (d *deployer) DownloadFile(ctx context.Context, prefix string, out string) (int, error) {
	if d.client == nil {
		err := d.setMinioClient()
		if err != nil || d.client == nil {
			return 0, err
		}
	}
	prefix, err := filepath.Abs(prefix)
	if err != nil {
		return 0, err
	}
	size := 0
	for _, client := range d.client {
		// check client is alive
		_, err = client.ListBuckets(ctx)
		if err != nil {
			continue
		}
		// list objects
		listObjects := client.ListObjects(ctx, "sealer", minio.ListObjectsOptions{
			Prefix:    modDir(prefix),
			Recursive: true,
		})
		for object := range listObjects {
			if object.Err != nil {
				return 0, object.Err
			}
			// download object
			outPutFile := filepath.Join(out, object.Key)
			err = client.FGetObject(ctx, "sealer", object.Key, outPutFile, minio.GetObjectOptions{})
			if err != nil {
				return 0, err
			}
			size++
		}
		return size, nil
	}
	return 0, errors.New("there is no minio node can work")
}

func (d *deployer) ListDir(ctx context.Context, dir string) ([]string, error) {
	if d.client == nil {
		err := d.setMinioClient()
		if err != nil || d.client == nil {
			return nil, err
		}
	}
	dir, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	for _, client := range d.client {
		// check client is alive
		_, err = client.ListBuckets(ctx)
		if err != nil {
			continue
		}
		// list objects
		listObjects := client.ListObjects(ctx, "sealer", minio.ListObjectsOptions{
			Prefix:    modDir(dir),
			Recursive: true,
		})
		res := make([]string, 0)
		for object := range listObjects {
			if object.Err != nil {
				return nil, object.Err
			}
			res = append(res, object.Key)
		}
		return res, nil
	}
	return nil, errors.New("there is no minio node can work")
}

func (d *deployer) RemoveFile(ctx context.Context, filename string) error {
	if d.client == nil {
		err := d.setMinioClient()
		if err != nil || d.client == nil {
			return err
		}
	}
	filename, err := filepath.Abs(filename)
	if err != nil {
		return err
	}
	for _, client := range d.client {
		// remove object
		err = client.RemoveObject(ctx, "sealer", modFile(filename), minio.RemoveObjectOptions{})
		if err != nil {
			continue
		}
		return nil
	}
	return errors.New("there is no minio node can work")
}

func (d *deployer) RemoveDir(ctx context.Context, dir string) error {
	if d.client == nil {
		err := d.setMinioClient()
		if err != nil || d.client == nil {
			return err
		}
	}
	dir, err := filepath.Abs(dir)
	if err != nil {
		return err
	}
	for _, client := range d.client {
		// check client is alive
		_, err = client.ListBuckets(ctx)
		if err != nil {
			continue
		}
		// list objects
		listObjects := client.ListObjects(ctx, "sealer", minio.ListObjectsOptions{
			Prefix:    modDir(dir),
			Recursive: true,
		})
		for object := range listObjects {
			if object.Err != nil {
				return object.Err
			}
			// remove object
			err = client.RemoveObject(ctx, "sealer", object.Key, minio.RemoveObjectOptions{})
			if err != nil {
				return err
			}
		}
		return nil
	}
	return errors.New("there is no minio node can work")
}

func (d *deployer) generateDownloadCmd() (map[string]string, error) {
	res := make(map[string]string)
	url, err := minioDownloadURL()
	if err != nil {
		return nil, err
	}
	binFileName := filepath.Join(d.config.BinDir, "minio")
	cmd := fmt.Sprintf(`
        binDir=%s
        cmd=%s
		datadir=%s

		if [ ! -d $datadir ]; then
  			mkdir -p $datadir
		fi

		if [ ! -d $binDir ]; then
            mkdir -p $binDir
        fi

        if [ -x $cmd ]; then
            echo "minio command is present and executable. No need to download."
        else
            url=%s
            curl -L -o $cmd $url
            chmod +x $cmd
        fi
    `, d.config.BinDir, binFileName, d.config.MinioDir, url)
	for _, ip := range d.config.Master {
		res[ip] = cmd
	}
	for _, ip := range d.config.Node {
		res[ip] = cmd
	}
	return res, nil
}

func (d *deployer) generateStartCmd() (map[string]string, error) {
	res := make(map[string]string)
	for _, node := range d.config.minioNode {
		res[node.IP] = ""
	}
	if len(res) == 1 {
		singleCmd := d.generateSingleCmd()
		for ip := range res {
			res[ip] = singleCmd
		}
	} else {
		clusterCmd := d.generateClusterCmd(res)
		for ip := range res {
			res[ip] = clusterCmd
		}
	}
	return res, nil
}

func (d *deployer) generateSingleCmd() string {
	return fmt.Sprintf(`setsid %s server --address %s:%v %s`,
		d.config.BinDir+"/minio", d.config.minioNode[0].IP, d.config.minioNode[0].Port, d.config.minioNode[0].Dir)
}

func (d *deployer) generateClusterCmd(res map[string]string) string {
	cmd := fmt.Sprintf("setsid %s server ", d.config.BinDir+"/minio")
	for ip := range res {
		cmd += fmt.Sprintf("http://%s%s ", ip, d.config.MinioDir)
	}
	return cmd
}

// https://dl.min.io/server/minio/release/linux-amd64/minio
func minioDownloadURL() (string, error) {
	switch runtime.GOOS {
	case "linux":
		switch runtime.GOARCH {
		case "amd64":
			return "https://dl.min.io/server/minio/release/linux-amd64/minio", nil
		case "arm64":
			return "https://dl.min.io/server/minio/release/linux-arm64/minio", nil
		default:
			return "", errors.New("unsupported arch")
		}
	default:
		return "", errors.New("unsupported os")
	}
}

func modDir(dir string) string {
	if strings.HasPrefix(dir, "/") {
		dir = dir[1:]
	}
	if strings.HasSuffix(dir, "/") {
		return dir
	}
	return dir + "/"
}

func modFile(filename string) string {
	if strings.HasPrefix(filename, "/") {
		filename = filename[1:]
	}
	if strings.HasSuffix(filename, "/") {
		filename = filename[:len(filename)-1]
	}
	return filename
}
