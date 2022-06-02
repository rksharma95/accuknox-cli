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

		handleFilterFlags(cmd)

		if err := network.StartHubbleRelay(networkOptions); err != nil {
			return err
		}
		return nil
	},
}

func handleFilterFlags(cmd *cobra.Command) {
	// not
	var isBlacklist bool = false
	if flag, _ := cmd.Flags().GetBool("not"); flag {
		isBlacklist = true
	}

	// ip
	if flag, _ := cmd.Flags().GetString("from-ip"); flag != "" {
		if isBlacklist {
			isBlacklist = false
			network.UpdateBlackList(&networkOptions, "from-ip", flag)
		} else {
			network.UpdateWhiteList(&networkOptions, "from-ip", flag)
		}
	}
	if flag, _ := cmd.Flags().GetString("to-ip"); flag != "" {
		if isBlacklist {
			isBlacklist = false
			network.UpdateBlackList(&networkOptions, "to-ip", flag)
		} else {
			network.UpdateWhiteList(&networkOptions, "to-ip", flag)
		}
	}

	// pod
	if flag, _ := cmd.Flags().GetString("from-pod"); flag != "" {
		if isBlacklist {
			isBlacklist = false
			network.UpdateBlackList(&networkOptions, "from-pod", flag)
		} else {
			network.UpdateWhiteList(&networkOptions, "from-pod", flag)
		}

	}
	if flag, _ := cmd.Flags().GetString("to-pod"); flag != "" {
		if isBlacklist {
			isBlacklist = false
			network.UpdateBlackList(&networkOptions, "to-pod", flag)
		} else {
			network.UpdateWhiteList(&networkOptions, "to-pod", flag)
		}
	}

	// fqdn
	if flag, _ := cmd.Flags().GetString("from-fqdn"); flag != "" {
		if isBlacklist {
			isBlacklist = false
			network.UpdateBlackList(&networkOptions, "from-fdqn", flag)
		} else {
			network.UpdateWhiteList(&networkOptions, "from-fqdn", flag)
		}
	}
	if flag, _ := cmd.Flags().GetString("to-fqdn"); flag != "" {
		if isBlacklist {
			isBlacklist = false
			network.UpdateBlackList(&networkOptions, "to-fdqn", flag)
		} else {
			network.UpdateWhiteList(&networkOptions, "to-fqdn", flag)
		}
	}

	// label
	if flag, _ := cmd.Flags().GetString("from-label"); flag != "" {
		if isBlacklist {
			isBlacklist = false
			network.UpdateBlackList(&networkOptions, "from-label", flag)
		} else {
			network.UpdateWhiteList(&networkOptions, "from-label", flag)
		}
	}
	if flag, _ := cmd.Flags().GetString("to-label"); flag != "" {
		if isBlacklist {
			isBlacklist = false
			network.UpdateBlackList(&networkOptions, "to-label", flag)
		} else {
			network.UpdateWhiteList(&networkOptions, "to-label", flag)
		}
	}

	// port
	if flag, _ := cmd.Flags().GetString("from-port"); flag != "" {
		if isBlacklist {
			isBlacklist = false
			network.UpdateBlackList(&networkOptions, "from-port", flag)
		} else {
			network.UpdateWhiteList(&networkOptions, "from-port", flag)
		}
	}
	if flag, _ := cmd.Flags().GetString("to-port"); flag != "" {
		if isBlacklist {
			isBlacklist = false
			network.UpdateBlackList(&networkOptions, "to-port", flag)
		} else {
			network.UpdateWhiteList(&networkOptions, "to-port", flag)
		}
	}
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

	// ip
	networkCmd.Flags().String("from-ip", "", "Show all flows originating at the given IP address.")
	networkCmd.Flags().String("to-ip", "", "Show all flows destined to the given IP address.")

	// pod
	networkCmd.Flags().String("from-pod", "", "Show all flows originating at the given pod name (e.g. \"/*.kubearmor.io\").")
	networkCmd.Flags().String("to-pod", "", "Show all flows destined to the given fully qualified domain name (e.g. \"/*.kubearmor.io\").")

	// fqdn
	networkCmd.Flags().String("from-fqdn", "", "Show all flows originating at the given fully qualified domain name (e.g. \"/*.kubearmor.io\").")
	networkCmd.Flags().String("to-fqdn", "", "Show all flows destined to the given fully qualified domain name (e.g. \"/*.kubearmor.io\").")

	// label
	networkCmd.Flags().String("from-label", "", "Show all flows originating at the given lebel.")
	networkCmd.Flags().String("to-label", "", "Show all flows destined to the given label")

	// port
	networkCmd.Flags().String("from-port", "", "Show all flows originating at the given port.")
	networkCmd.Flags().String("to-port", "", "Show all flows destined to the given port.")

	// not
	networkCmd.Flags().Bool("not", false, "reverse the effect of a flag.")

}
