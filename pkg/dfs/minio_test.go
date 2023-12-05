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

//
//import (
//	"context"
//	"github.com/stretchr/testify/assert"
//	"os"
//	"testing"
//)
//
//var config = &Config{
//	Master: []string{"172.19.0.2", "172.19.0.3"},
//	Node:   []string{"172.19.0.4", "172.19.0.5"},
//	User:   "root",
//	Passwd: "123456",
//}
//
//var ctx = context.Background()
//
//func TestNewDeployer(t *testing.T) {
//	_, err := NewDeployer(config)
//	assert.Nil(t, err)
//}
//
//func TestStart(t *testing.T) {
//	d, err := NewDeployer(config)
//	assert.Nil(t, err)
//	err = d.Start(ctx)
//	assert.Nil(t, err)
//	if err != nil {
//		panic(err)
//	}
//	ok := d.IsRunning(ctx)
//	assert.True(t, ok)
//}
//
//func TestUploadFilename(t *testing.T) {
//	d, err := NewDeployer(config)
//	assert.Nil(t, err)
//	dir, err := os.Getwd()
//	assert.Nil(t, err)
//	err = d.UploadFile(ctx, dir+"/test/hello.txt")
//	assert.Nil(t, err)
//	listDir, err := d.ListDir(ctx, dir)
//	assert.Nil(t, err)
//	assert.Equal(t, 1, len(listDir))
//	err = d.RemoveDir(ctx, dir+"/test")
//	assert.Nil(t, err)
//}
//
//func TestUploadDir(t *testing.T) {
//	d, err := NewDeployer(config)
//	assert.Nil(t, err)
//	dir, err := os.Getwd()
//	assert.Nil(t, err)
//	err = d.UploadDir(ctx, dir+"/test")
//	assert.Nil(t, err)
//	listDir, err := d.ListDir(ctx, dir)
//	assert.Nil(t, err)
//	// 3 files:
//	//      dir1/hello.txt
//	//      dir2/hello.txt
//	//      hello.txt
//	assert.Equal(t, 3, len(listDir))
//	err = d.RemoveDir(ctx, dir+"/test")
//	assert.Nil(t, err)
//}
//
//func TestUploadDir1(t *testing.T) {
//	d, err := NewDeployer(config)
//	assert.Nil(t, err)
//	dir, err := os.Getwd()
//	assert.Nil(t, err)
//	err = d.UploadDir(ctx, dir+"/test")
//	assert.Nil(t, err)
//	listDir, err := d.ListDir(ctx, dir+"/test/dir1")
//	assert.Nil(t, err)
//	// 1 files:
//	//      dir1/hello.txt
//	assert.Equal(t, 1, len(listDir))
//	err = d.RemoveDir(ctx, dir)
//	assert.Nil(t, err)
//}
//
//func TestDownloadFile(t *testing.T) {
//	d, err := NewDeployer(config)
//	assert.Nil(t, err)
//	dir, err := os.Getwd()
//	assert.Nil(t, err)
//	err = d.UploadDir(ctx, dir+"/test")
//	assert.Nil(t, err)
//	// you can see the file in /tmp
//	// root@loubens:/tmp/home/loubens/project/sealer/pkg/dfs/test# pwd
//	// /tmp/home/loubens/project/sealer/pkg/dfs/test
//	// root@loubens:/tmp/home/loubens/project/sealer/pkg/dfs/test# cat hello.txt
//	// hello sealer
//	count, err := d.DownloadFile(ctx, dir+"/test", "/tmp/")
//	assert.Nil(t, err)
//	assert.Equal(t, 3, count)
//	err = d.RemoveDir(ctx, dir+"/test")
//	assert.Nil(t, err)
//}
//
//func TestRemoveFile(t *testing.T) {
//	d, err := NewDeployer(config)
//	assert.Nil(t, err)
//	dir, err := os.Getwd()
//	assert.Nil(t, err)
//	err = d.UploadDir(ctx, dir+"/test")
//	assert.Nil(t, err)
//	list, err := d.ListDir(ctx, dir+"/test")
//	assert.Nil(t, err)
//	assert.Equal(t, 3, len(list))
//	for _, l := range list {
//		err = d.RemoveFile(ctx, "/"+l)
//		assert.Nil(t, err)
//	}
//	list, err = d.ListDir(ctx, dir+"/test")
//	assert.Nil(t, err)
//	assert.Equal(t, 0, len(list))
//}
//
//func TestRemoveDir(t *testing.T) {
//	d, err := NewDeployer(config)
//	assert.Nil(t, err)
//	dir, err := os.Getwd()
//	assert.Nil(t, err)
//	err = d.UploadDir(ctx, dir+"/test")
//	assert.Nil(t, err)
//	list, err := d.ListDir(ctx, dir+"/test")
//	assert.Nil(t, err)
//	assert.Equal(t, 3, len(list))
//	err = d.RemoveDir(ctx, dir+"/test")
//	assert.Nil(t, err)
//	list, err = d.ListDir(ctx, dir+"/test")
//	assert.Nil(t, err)
//	assert.Equal(t, 0, len(list))
//}
//
//func TestUploadDir2(t *testing.T) {
//	dir1 := "./test/dir1"
//	d, err := NewDeployer(config)
//	assert.Nil(t, err)
//	err = d.UploadDir(ctx, dir1)
//	assert.Nil(t, err)
//	list, err := d.ListDir(ctx, dir1)
//	assert.Nil(t, err)
//	assert.Equal(t, 1, len(list))
//	err = d.RemoveDir(ctx, dir1)
//	assert.Nil(t, err)
//}
//
//func TestDownloadFile1(t *testing.T) {
//	dir1 := "./test/dir1"
//	d, err := NewDeployer(config)
//	assert.Nil(t, err)
//	err = d.UploadDir(ctx, dir1)
//	assert.Nil(t, err)
//	list, err := d.ListDir(ctx, dir1)
//	assert.Nil(t, err)
//	assert.Equal(t, 1, len(list))
//	count, err := d.DownloadFile(ctx, dir1, "/tmp/")
//	assert.Nil(t, err)
//	assert.Equal(t, 1, count)
//	err = d.RemoveDir(ctx, dir1)
//	assert.Nil(t, err)
//}
//
//func TestStop(t *testing.T) {
//	d, err := NewDeployer(config)
//	assert.Nil(t, err)
//	err = d.Stop(ctx)
//	assert.Nil(t, err)
//	ok := d.IsRunning(ctx)
//	assert.False(t, ok)
//}
