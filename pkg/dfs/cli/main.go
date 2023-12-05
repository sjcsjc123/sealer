// Copyright © 2023 Alibaba Group Holding Ltd.
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

package main

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"sealer/pkg/dfs"
)

func main() {
	rootCmd := cobra.Command{
		Use:   "dfs-cli",
		Short: "A tool to build, share and run any distributed applications.",
	}
	rootCmd.AddCommand(startCmd())
	rootCmd.AddCommand(uploadCmd())
	rootCmd.AddCommand(downloadCmd())
	rootCmd.AddCommand(listCmd())
	rootCmd.AddCommand(removeCmd())
	rootCmd.AddCommand(stopCmd())
	rootCmd.AddCommand(isRunningCmd())
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

var (
	config   = &dfs.Config{}
	filename string
	dir      string
	prefix   string
	out      string
)

func startCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "start to run a dfs cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			deployer, err := dfs.NewDeployer(config)
			if err != nil {
				return err
			}
			err = deployer.Start(context.Background())
			if err != nil {
				return err
			}
			return nil
		},
	}
	cmd.Flags().StringSliceVar(&config.Master, "master", []string{}, "master ip list")
	cmd.Flags().StringSliceVar(&config.Node, "node", []string{}, "node ip list")
	cmd.Flags().StringVar(&config.Passwd, "passwd", "", "ssh password")
	return cmd
}

func uploadCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upload",
		Short: "upload file to dfs",
		RunE: func(cmd *cobra.Command, args []string) error {
			deployer, err := dfs.NewDeployer(config)
			if err != nil {
				return err
			}
			if filename == "" && dir == "" {
				return errors.New("filename and dir can not be empty at the same time")
			}
			if filename != "" {
				err = deployer.UploadFile(context.Background(), filename)
				if err != nil {
					return err
				}
				return nil
			}
			return deployer.UploadDir(context.Background(), dir)
		},
	}
	cmd.Flags().StringSliceVar(&config.Master, "master", []string{}, "master ip list")
	cmd.Flags().StringSliceVar(&config.Node, "node", []string{}, "node ip list")
	cmd.Flags().StringVar(&config.Passwd, "passwd", "", "ssh password")
	cmd.Flags().StringVar(&filename, "filename", "", "filename to upload")
	cmd.Flags().StringVar(&dir, "dir", "", "directory to upload")
	return cmd
}

func downloadCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "download",
		Short: "download file from dfs",
		RunE: func(cmd *cobra.Command, args []string) error {
			deployer, err := dfs.NewDeployer(config)
			if err != nil {
				return err
			}
			if prefix == "" || out == "" {
				return errors.New("prefix and out can not be empty at the same time")
			}
			_, err = deployer.DownloadFile(context.Background(), prefix, out)
			if err != nil {
				return err
			}
			return nil
		},
	}
	cmd.Flags().StringSliceVar(&config.Master, "master", []string{}, "master ip list")
	cmd.Flags().StringSliceVar(&config.Node, "node", []string{}, "node ip list")
	cmd.Flags().StringVar(&config.Passwd, "passwd", "", "ssh password")
	cmd.Flags().StringVar(&prefix, "prefix", "", "prefix to download")
	cmd.Flags().StringVar(&out, "out", "", "output directory")
	return cmd
}

func listCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list directory from dfs",
		RunE: func(cmd *cobra.Command, args []string) error {
			deployer, err := dfs.NewDeployer(config)
			if err != nil {
				return err
			}
			files, err := deployer.ListDir(context.Background(), dir)
			if err != nil {
				return err
			}
			for _, file := range files {
				logrus.Info(file)
			}
			return nil
		},
	}
	cmd.Flags().StringSliceVar(&config.Master, "master", []string{}, "master ip list")
	cmd.Flags().StringSliceVar(&config.Node, "node", []string{}, "node ip list")
	cmd.Flags().StringVar(&config.Passwd, "passwd", "", "ssh password")
	cmd.Flags().StringVar(&dir, "dir", "", "directory to list")
	return cmd
}

func removeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove",
		Short: "remove file from dfs",
		RunE: func(cmd *cobra.Command, args []string) error {
			deployer, err := dfs.NewDeployer(config)
			if err != nil {
				return err
			}
			if filename == "" && dir == "" {
				return errors.New("filename and dir can not be empty at the same time")
			}
			if filename != "" {
				err = deployer.RemoveFile(context.Background(), filename)
				if err != nil {
					return err
				}
				return nil
			}
			return deployer.RemoveDir(context.Background(), dir)
		},
	}
	cmd.Flags().StringSliceVar(&config.Master, "master", []string{}, "master ip list")
	cmd.Flags().StringSliceVar(&config.Node, "node", []string{}, "node ip list")
	cmd.Flags().StringVar(&config.Passwd, "passwd", "", "ssh password")
	cmd.Flags().StringVar(&filename, "filename", "", "filename to remove")
	cmd.Flags().StringVar(&dir, "dir", "", "directory to remove")
	return cmd
}

func stopCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "stop dfs",
		RunE: func(cmd *cobra.Command, args []string) error {
			deployer, err := dfs.NewDeployer(config)
			if err != nil {
				return err
			}
			err = deployer.Stop(context.Background())
			if err != nil {
				return err
			}
			return nil
		},
	}
	cmd.Flags().StringSliceVar(&config.Master, "master", []string{}, "master ip list")
	cmd.Flags().StringSliceVar(&config.Node, "node", []string{}, "node ip list")
	cmd.Flags().StringVar(&config.Passwd, "passwd", "", "ssh password")
	return cmd
}

func isRunningCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "isRunning",
		Short: "check dfs is running or not",
		RunE: func(cmd *cobra.Command, args []string) error {
			deployer, err := dfs.NewDeployer(config)
			if err != nil {
				return err
			}
			ok := deployer.IsRunning(context.Background())
			logrus.Infof("dfs is running: %v", ok)
			return nil
		},
	}
	cmd.Flags().StringSliceVar(&config.Master, "master", []string{}, "master ip list")
	cmd.Flags().StringSliceVar(&config.Node, "node", []string{}, "node ip list")
	cmd.Flags().StringVar(&config.Passwd, "passwd", "", "ssh password")
	return cmd
}
