# accuknox-cli

`accuknox` is a CLI client to help manage [KubeArmor](https://github.com/kubearmor/KubeArmor), [Discovery-engine](https://github.com/accuknox/discovery-engine) and [Cilium](https://github.com/cilium/cilium).

KubeArmor is a container-aware runtime security enforcement system that
restricts the behavior (such as process execution, file access, and networking
operation) of containers at the system level.

Discovery Engine discovers the security posture for your workloads and auto-discovers the policy-set required to put the workload in least-permissive mode. The engine leverages the rich visibility provided by KubeArmor and Cilium to auto discover the systems and network security posture.

## Installation

The following sections show how to install the `accuknox` CLI. It can be installed either from source, or from pre-built binary releases.

### From Script

`accuknox` has an installer script that will automatically grab the latest version of accuknox and install it locally.

```
curl -sfL https://raw.githubusercontent.com/accuknox/accuknox-cli/main/install.sh | sudo sh -s -- -b /usr/local/bin
```

The binary will be installed in `/usr/local/bin` folder.

### From Source 

Building accuknox from source provides the latest (pre-release) `accuknox` version.

```
git clone https://github.com/accuknox/accuknox-cli
cd accuknox-cli
make install
```

## Usage

```
CLI Utility to help manage Accuknox security solution
	
accuknox-cli tool helps to install, manage and troubleshoot Accuknox security solution

Usage:
  accuknox [command]

Available Commands:
  completion   Generate the autocompletion script for the specified shell
  discover     Discover applicable policies
  help         Help about any command
  install      Install KubeArmor, Cilium and Discovery-engine in a Kubernetes Cluster
  log          Observe Logs from KubeArmor
  port-forward port-forward KubeArmor, Cilium and Discovery-engine in a Kubernetes Cluster
  selfupdate   selfupdate this cli tool
  summary      Policy summary from discovery engine
  sysdump      Collect system dump information for troubleshooting and error report
  uninstall    Uninstall KubeArmor, Cilium and Discovery-engine from a Kubernetes Cluster
  version      Display version information
  vm           Vm commands for kvmservice


Flags:
  -h, --help   help for accuknox

Use "accuknox [command] --help" for more information about a command.
```
