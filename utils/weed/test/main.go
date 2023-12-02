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

func main() {
	//rootCmd := cobra.Command{
	//	Use:   "weed",
	//	Short: "A tool to build, share and run any distributed applications.",
	//}
	//rootCmd.AddCommand(startCmd())
	//rootCmd.AddCommand(writeCmd())
	//rootCmd.AddCommand(downloadFileCmd())
	//if err := rootCmd.Execute(); err != nil {
	//	logrus.Errorf("sealer-%s: %v", version.GetSingleVersion(), err)
	//	os.Exit(1)
	//}

}

//var config = &weed.Config{}
//
//func startCmd() *cobra.Command {
//	//cmd := &cobra.Command{
//	//	Use:   "start",
//	//	Short: "start to run a weed cluster",
//	//	RunE: func(cmd *cobra.Command, args []string) error {
//	//		deployer := weed.NewDeployer(config)
//	//		err := deployer.CreateEtcdCluster(context.Background())
//	//		if err != nil {
//	//			return err
//	//		}
//	//		err = deployer.CreateWeedCluster(context.Background())
//	//		if err != nil {
//	//			return err
//	//		}
//	//		return nil
//	//	},
//	//}
//	//cmd.Flags().StringSliceVar(&config.MasterIP, "master-ip", []string{}, "master ip list")
//	//cmd.Flags().StringSliceVar(&config.VolumeIP, "volume-ip", []string{}, "volume ip list")
//	//cmd.Flags().StringVar(&config.BinDir, "bin-dir", "", "bin dir")
//	//cmd.Flags().StringVar(&config.EtcdConfigPath, "etcd-config-path", "", "etcd config path")
//	//cmd.Flags().StringVar(&config.CurrentIP, "current-ip", "", "current ip")
//	//cmd.Flags().StringVar(&config.WeedMasterDir, "weed-master-dir", "", "weed master dir")
//	//cmd.Flags().StringVar(&config.WeedVolumeDir, "weed-volume-dir", "", "weed volume dir")
//	//cmd.Flags().StringVar(&config.DefaultReplication, "default-replication", "", "default replication")
//	//return cmd
//}
//
//var dir string
//
//func writeCmd() *cobra.Command {
//	cmd := &cobra.Command{
//		Use:   "write",
//		Short: "write data to weed cluster",
//		RunE: func(cmd *cobra.Command, args []string) error {
//			deployer := weed.NewDeployer(config)
//			err := deployer.UploadFile(context.Background(), dir)
//			if err != nil {
//				return err
//			}
//			return nil
//		},
//	}
//	cmd.Flags().StringVar(&dir, "dir", "", "dir")
//	return cmd
//}
//
//var out string
//
//func downloadFileCmd() *cobra.Command {
//	cmd := &cobra.Command{
//		Use:   "download",
//		Short: "download data from weed cluster",
//		RunE: func(cmd *cobra.Command, args []string) error {
//			deployer := weed.NewDeployer(config)
//			err := deployer.DownloadFile(context.Background(), dir, out)
//			if err != nil {
//				return err
//			}
//			return nil
//		},
//	}
//	cmd.Flags().StringVar(&dir, "dir", "", "dir")
//	cmd.Flags().StringVar(&out, "out", "", "out")
//	return cmd
//}
