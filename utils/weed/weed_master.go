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
	"sync"
)

type Master interface {
	UploadFile(ctx context.Context, master string, dir string) ([]UploadFileResponse, error)
	DownloadFile(ctx context.Context, master string, fid string, outputDir string) error
	RemoveFile(ctx context.Context, master string, dir string) error
}

type master struct {
	ip                 string
	port               int
	mDir               string
	defaultReplication string
	peers              []string
	needMoreLocalNode  bool
	portList           []int
	mDirList           []string
	wg                 *sync.WaitGroup
	logDir             string
	pidDir             string
}

type UploadFileResponse struct {
	Fid      string `json:"fid"`
	URL      string `json:"url"`
	FileName string `json:"fileName"`
	Size     int64  `json:"size"`
}

func (m *master) UploadFile(ctx context.Context, master string, dir string) ([]UploadFileResponse, error) {
	//runOptions := RunOptions{
	//	Binary: "weed",
	//	Name:   "upload",
	//	args:   m.buildUploadFileArgs(ctx, master, dir),
	//}
	//jsonResponse, err := runBinaryWithJSONResponse(ctx, &runOptions, m.wg)
	//if err != nil {
	//	return []UploadFileResponse{}, err
	//}
	//var uploadFileResponse []UploadFileResponse
	//err = json.Unmarshal(jsonResponse, &uploadFileResponse)
	//fmt.Println(string(jsonResponse))
	//if err != nil {
	//	return []UploadFileResponse{}, err
	//}
	//return uploadFileResponse, nil
	return nil, nil
}

func (m *master) buildUploadFileArgs(ctx context.Context, params ...interface{}) []string {
	_ = ctx
	return []string{
		"upload",
		"-master=" + params[0].(string),
		"-dir=" + params[1].(string),
	}
}

func (m *master) buildDownloadFileArgs(ctx context.Context, params ...interface{}) []string {
	_ = ctx
	return []string{
		"-server=" + params[0].(string),
		"--dir=" + params[2].(string),
		params[1].(string),
	}
}

func (m *master) DownloadFile(ctx context.Context, master string, fid string, outputDir string) error {
	//runOptions := RunOptions{
	//	Binary: "weed",
	//	Name:   "download",
	//	args:   m.buildDownloadFileArgs(ctx, master, fid, outputDir),
	//}
	//err := runBinary(ctx, &runOptions, m.wg)
	//if err != nil {
	//	return err
	//}
	//return nil
	return nil
}

func (m *master) RemoveFile(ctx context.Context, master string, fid string) error {
	//TODO weed may not support remove file, may be should consider to use other file system
	panic("implement me")
}

func NewWeedMaster(config *Config) Master {
	return &master{}
}
