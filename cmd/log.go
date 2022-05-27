// SPDX-License-Identifier: Apache-2.0
// Copyright 2021 Authors of KubeArmor

package cmd

import (
	"github.com/accuknox/accuknox-cli/network"

	"github.com/kubearmor/kubearmor-client/log"
	"github.com/spf13/cobra"
)

var logOptions log.Options
var networkOptions network.Options

// logCmd represents the log command
var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Observe Logs from KubeArmor",
	Long:  `Observe Logs from KubeArmor`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := cmd.Help(); err != nil {
			return err
		}
		return nil
	},
}

var appCmd = &cobra.Command{
	Use:   "application",
	Short: "Observe Logs from KubeArmor",
	Long:  `Observe Logs from KubeArmor`,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.StopChan = make(chan struct{})
		if err := log.StartObserver(logOptions); err != nil {
			return err
		}
		return nil
	},
}

var networkCmd = &cobra.Command{

	Use:   "network",
	Short: "Observe Logs from hubble relay",
	Long:  `Observe Logs from hubble relay`,
	RunE: func(cmd *cobra.Command, args []string) error {

		network.StopChan = make(chan struct{})
		if err := network.StartHubbleRelay(networkOptions); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(logCmd)

	logCmd.AddCommand(networkCmd)
	logCmd.AddCommand(appCmd)

	appCmd.Flags().StringVar(&logOptions.GRPC, "gRPC", "", "gRPC server information")
	appCmd.Flags().StringVar(&logOptions.MsgPath, "msgPath", "none", "Output location for messages, {path|stdout|none}")
	appCmd.Flags().StringVar(&logOptions.LogPath, "logPath", "stdout", "Output location for alerts and logs, {path|stdout|none}")
	appCmd.Flags().StringVar(&logOptions.LogFilter, "logFilter", "policy", "Filter for what kinds of alerts and logs to receive, {policy|system|all}")
	appCmd.Flags().BoolVar(&logOptions.JSON, "json", false, "Flag to print alerts and logs in the JSON format")
	appCmd.Flags().StringVar(&logOptions.Namespace, "namespace", "", "Specify the namespace")
	appCmd.Flags().StringVar(&logOptions.Operation, "operation", "", "Give the type of the operation (Eg:Process/File/Network)")
	appCmd.Flags().StringVar(&logOptions.LogType, "logType", "", "Log type you want (Eg:ContainerLog/HostLog) ")
	appCmd.Flags().StringVar(&logOptions.ContainerName, "container", "", "name of the container ")
	appCmd.Flags().StringVar(&logOptions.PodName, "pod", "", "name of the pod ")
	appCmd.Flags().StringVar(&logOptions.Resource, "resource", "", "command used by the user")
	appCmd.Flags().StringVar(&logOptions.Source, "source", "", "binary used by the system ")
	appCmd.Flags().Uint32Var(&logOptions.Limit, "limit", 0, "number of logs you want to see")

	networkCmd.Flags().BoolVarP(&networkOptions.Follow, "follow", "f", false, "Follow flows output")

	// filter flags
	// networkCmd.Flags().String()
}
