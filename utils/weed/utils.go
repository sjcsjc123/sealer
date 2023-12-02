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
	"fmt"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"os"
	"runtime"
)

//type RunOptions struct {
//	Binary  string
//	Name    string
//	args    []string
//	logFile string
//	pidFile string
//}

//func runBinary(ctx context.Context, option *RunOptions, wg *sync.WaitGroup) error {
//	cmd := exec.CommandContext(ctx, option.Binary, option.args...)
//
//	var outputFileWriter *bufio.Writer
//	if option.logFile != "" {
//		// output to binary.
//		outputFile, err := os.Create(option.logFile)
//		if err != nil {
//			return err
//		}
//
//		outputFileWriter = bufio.NewWriter(outputFile)
//		cmd.Stdout = outputFileWriter
//		cmd.Stderr = outputFileWriter
//	}
//
//	if err := cmd.Start(); err != nil {
//		return err
//	}
//
//	if option.pidFile != "" {
//		pid := strconv.Itoa(cmd.Process.Pid)
//
//		f, err := os.Create(option.pidFile)
//		if err != nil {
//			return err
//		}
//
//		_, err = f.Write([]byte(pid))
//		if err != nil {
//			return err
//		}
//	}
//
//	go func() {
//		defer wg.Done()
//		wg.Add(1)
//		if err := cmd.Wait(); err != nil {
//			// Caught signal kill and interrupt error then ignore.
//			var exit *exec.ExitError
//			if errors.As(err, &exit) {
//				if status, ok := exit.Sys().(syscall.WaitStatus); ok {
//					if status.Signaled() &&
//						(status.Signal() == syscall.SIGKILL || status.Signal() == syscall.SIGINT) {
//						return
//					}
//				}
//			}
//			if outputFileWriter != nil {
//				_ = outputFileWriter.Flush()
//			}
//		}
//	}()
//
//	return nil
//}
//
//func runBinaryWithJSONResponse(ctx context.Context, option *RunOptions, wg *sync.WaitGroup) ([]byte, error) {
//	cmd := exec.CommandContext(ctx, option.Binary, option.args...)
//	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
//
//	var jsonOutput bytes.Buffer
//	cmd.Stdout = &jsonOutput
//	if err := cmd.Run(); err != nil {
//		return nil, err
//	}
//	jsonResponse := jsonOutput.Bytes()
//	return jsonResponse, nil
//}
//
//func CreateDirIfNotExists(dir string) (err error) {
//	if err := os.MkdirAll(dir, 0755); err != nil && !os.IsExist(err) {
//		return err
//	}
//	return nil
//}

func weedDownloadURL() (string, error) {
	if runtime.GOOS != "linux" {
		return "", errors.New("unsupported os")
	}
	url := fmt.Sprintf("https://github.com/%s/%s/releases/download/%s/", SeaWeedFsGitHubOrg, SeaWeedFsGithubRepo, SeaWeedFsVersion)
	switch arch := runtime.GOARCH; arch {
	case "amd64":
		url += "linux_amd64.tar.gz"
	case "arm64":
		url += "linux_arm.tar.gz"
	default:
		return "", errors.New("unsupported arch")
	}
	return url, nil
}

func downloadFile(url string, dest string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// check if the destination folder exists
	_, err = os.Stat(dest)
	if err == nil {
		_ = os.RemoveAll(dest)
	}

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
