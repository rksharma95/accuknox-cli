// SPDX-License-Identifier: Apache-2.0
// Copyright 2021 Authors of KubeArmor

package cmd

import (
	"errors"
	"fmt"
	"net"

	"github.com/kubearmor/kubearmor-client/vm"
	"github.com/spf13/cobra"
)

var (
	vmScriptOptions vm.ScriptOptions
	vmLabelOptions  vm.LabelOptions
	HttpIP          string
	HttpPort        string
	IsKvmsEnv       bool
)

// vmCmd represents the vm command
var vmCmd = &cobra.Command{
	Use:   "vm",
	Short: "Vm commands for kvmservice",
	Long:  `Vm commands for kvmservice`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := cmd.Help(); err != nil {
			return err
		}
		return nil
	},
}

// vmScriptCmd represents the vm command for script download
var vmScriptCmd = &cobra.Command{
	Use:   "getscript",
	Short: "download vm installation script for kvms control plane",
	Long:  `download vm installation script for kvms control plane`,
	RunE: func(cmd *cobra.Command, args []string) error {

		if err := vm.GetScript(client, vmScriptOptions, HttpIP, IsKvmsEnv); err != nil {
			return err
		}
		return nil
	},
}

// ================================ //
// == vm label [add|delete|list] == //
// ================================ //

// vmLabelCmd represents the vm command for label management
var vmLabelCmd = &cobra.Command{
	Use:   "label",
	Short: "label handling for kvms control plane vm",
	Long:  `label handling for kvms control plane vm`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := cmd.Help(); err != nil {
			return err
		}
		return nil
	},
}

// vmLabelAddCmd represents the vm add label command for label management
var vmLabelAddCmd = &cobra.Command{
	Use:   "add",
	Short: "add label for kvms control plane vm",
	Long:  `add label for kvms control plane vm`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Create http address
		httpAddress := "http://" + net.JoinHostPort(HttpIP, HttpPort)

		if err := vm.LabelHandling("ADD", vmLabelOptions, httpAddress, IsKvmsEnv); err != nil {
			return err
		}
		return nil
	},
}

// vmLabelDeleteCmd represents the vm add label command for label management
var vmLabelDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete label for kvms control plane vm",
	Long:  `delete label for kvms control plane vm`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Create http address
		httpAddress := "http://" + net.JoinHostPort(HttpIP, HttpPort)

		if err := vm.LabelHandling("DELETE", vmLabelOptions, httpAddress, IsKvmsEnv); err != nil {
			return err
		}
		return nil
	},
}

// vmLabelListCmd represents the vm list label command for label management
var vmLabelListCmd = &cobra.Command{
	Use:   "list",
	Short: "list labels for vm in k8s/nonk8s control plane",
	Long:  `list labels for vm in k8s/nonk8s control plane`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Create http address
		httpAddress := "http://" + net.JoinHostPort(HttpIP, HttpPort)

		if err := vm.LabelHandling("LIST", vmLabelOptions, httpAddress, IsKvmsEnv); err != nil {
			return err
		}
		return nil
	},
}

// ========================= //
// == vm [add|delete|list] = //
// ========================= //

// vmOnboardAddCmd represents the command for vm onboarding
var vmOnboardAddCmd = &cobra.Command{
	Use:   "add",
	Short: "onboard new VM onto kvms control plane vm",
	Long:  `onboard new VM onto kvms control plane vm`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires a path to valid vm YAML as argument")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		httpAddress := "http://" + net.JoinHostPort(HttpIP, HttpPort)
		if err := vm.Onboarding("ADDED", args[0], httpAddress); err != nil {
			return err
		}
		return nil
	},
}

// vmOnboardDeleteCmd represents the command for vm offboarding
var vmOnboardDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "offboard existing VM from kvms control plane vm",
	Long:  `offboard existing VM from kvms control plane vm`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires a path to valid vm YAML as argument")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		httpAddress := "http://" + net.JoinHostPort(HttpIP, HttpPort)
		if err := vm.Onboarding("DELETED", args[0], httpAddress); err != nil {
			return err
		}
		return nil
	},
}

// vmListCmd represents the command for vm listing
var vmListCmd = &cobra.Command{
	Use:   "list",
	Short: "list configured VMs",
	Long:  `list configured VMs`,
	RunE: func(cmd *cobra.Command, args []string) error {
		httpAddress := "http://" + net.JoinHostPort(HttpIP, HttpPort)
		if err := vm.List(httpAddress); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(vmCmd)

	// subcommands for vm command
	vmCmd.AddCommand(vmScriptCmd)
	vmCmd.AddCommand(vmLabelCmd)
	vmCmd.AddCommand(vmOnboardAddCmd)
	vmCmd.AddCommand(vmOnboardDeleteCmd)
	vmCmd.AddCommand(vmListCmd)

	// subcommands for vm label command
	vmLabelCmd.AddCommand(vmLabelAddCmd)
	vmLabelCmd.AddCommand(vmLabelDeleteCmd)
	vmLabelCmd.AddCommand(vmLabelListCmd)

	// vm script download Options
	vmScriptCmd.Flags().StringVarP(&vmScriptOptions.Port, "port", "p", "32770", "Port of kvmservice")
	vmScriptCmd.Flags().StringVarP(&vmScriptOptions.VMName, "kvm", "v", "", "Name of configured vm")
	vmScriptCmd.Flags().StringVarP(&vmScriptOptions.File, "file", "f", "none", "Filename with path to store the configured vm installation script")

	// Marking this flag as markedFlag and mandatory
	err := vmScriptCmd.MarkFlagRequired("kvm")
	if err != nil {
		_ = fmt.Errorf("kvm option not supplied")
	}

	// options for vm generic commands related to HTTP Request
	vmCmd.PersistentFlags().StringVar(&HttpIP, "http-ip", "127.0.0.1", "IP of kvm-service")
	vmCmd.PersistentFlags().StringVar(&HttpPort, "http-port", "8000", "Port of kvm-service")
	vmCmd.PersistentFlags().BoolVar(&IsKvmsEnv, "kvms", false, "Enable if kvms environment/control-plane")

	// options for vm label command
	vmLabelCmd.PersistentFlags().StringVar(&vmLabelOptions.VMName, "vm", "", "VM name")
	vmLabelCmd.PersistentFlags().StringVar(&vmLabelOptions.VMLabels, "label", "", "list of labels")

}
