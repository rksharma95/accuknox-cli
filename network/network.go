package network

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/cilium/cilium/api/v1/flow"
	"github.com/cilium/cilium/api/v1/observer"
	"google.golang.org/grpc"
)

// Options Structure
type Options struct {
	Follow    bool
	whitelist []*flow.FlowFilter
	blacklist []*flow.FlowFilter
}

// StopChan Channel
var (
	StopChan chan struct{}
)

// ConnectHubbleRelay Function
func ConnectHubbleRelay() (*grpc.ClientConn, error) {
	addr := "localhost:4245"

	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {

		return nil, err
	}

	return conn, err
}

func UpdateBlackList(o *Options, flag string, value string) {
	ff := &flow.FlowFilter{}

	switch flag {
	// ip
	case "from-ip":
		ff.SourceIp = append(ff.SourceIp, value)
	case "to-ip":
		ff.DestinationIp = append(ff.DestinationIp, value)

	// pod
	case "from-pod":
		ff.SourcePod = append(ff.SourcePod, value)

	case "to-pod":
		ff.DestinationPod = append(ff.DestinationPod, value)

	// fqdn
	case "from-fqdn":
		ff.SourceFqdn = append(ff.SourceFqdn, value)

	case "to-fqdn":
		ff.DestinationFqdn = append(ff.DestinationFqdn, value)

	// label
	case "from-label":
		ff.SourceLabel = append(ff.SourceLabel, value)

	case "to-label":
		ff.DestinationLabel = append(ff.DestinationLabel, value)

	// service
	case "from-service":

		ff.SourceService = append(ff.SourceService, value)

	case "to-service":
		ff.DestinationService = append(ff.DestinationService, value)

	// port
	case "from-port":
		ff.SourcePort = append(ff.SourcePort, value)

	case "to-port":
		ff.DestinationPort = append(ff.DestinationPort, value)

	// verdict
	case "verdict":
		v, ok := flow.Verdict_value[value]

		if !ok {
			fmt.Printf("invalid --verdict: %v", value)
			os.Exit(1)
		}
		ff.Verdict = append(ff.Verdict, flow.Verdict(v))

		// TODO --http-status, --http-method, --http-path, --tcp-flag, --ip-version, --nodename
	}

	o.blacklist = append(o.blacklist, ff)
}

func UpdateWhiteList(o *Options, flag string, value string) {
	ff := &flow.FlowFilter{}

	switch flag {
	// ip
	case "from-ip":
		ff.SourceIp = append(ff.SourceIp, value)
	case "to-ip":
		ff.DestinationIp = append(ff.DestinationIp, value)

	// pod
	case "from-pod":
		ff.SourcePod = append(ff.SourcePod, value)

	case "to-pod":
		ff.DestinationPod = append(ff.DestinationPod, value)

	// fqdn
	case "from-fqdn":
		ff.SourceFqdn = append(ff.SourceFqdn, value)

	case "to-fqdn":
		ff.DestinationFqdn = append(ff.DestinationFqdn, value)

	// label
	case "from-label":
		ff.SourceLabel = append(ff.SourceLabel, value)

	case "to-label":
		ff.DestinationLabel = append(ff.DestinationLabel, value)

	// service
	case "from-service":

		ff.SourceService = append(ff.SourceService, value)

	case "to-service":
		ff.DestinationService = append(ff.DestinationService, value)

	// port
	case "from-port":
		ff.SourcePort = append(ff.SourcePort, value)

	case "to-port":
		ff.DestinationPort = append(ff.DestinationPort, value)

	// verdict
	case "verdict":
		v, ok := flow.Verdict_value[value]

		if !ok {
			fmt.Printf("invalid --verdict: %v", value)
			os.Exit(1)
		}
		ff.Verdict = append(ff.Verdict, flow.Verdict(v))
	}
	// TODO --http-status, --http-method, --http-path, --tcp-flag, --ip-version, --nodename

	o.whitelist = append(o.whitelist, ff)

}

// StartHubbleRelay Function
func StartHubbleRelay(o Options) error {

	conn, err := ConnectHubbleRelay()
	if err != nil {
		return err
	}

	defer func() {
		_ = conn.Close()
	}()

	client := observer.NewObserverClient(conn)

	var num uint64
	switch {
	case o.Follow:
		num = ^uint64(0)
	default:
		num = 20
	}

	req := &observer.GetFlowsRequest{
		Number:    num,
		Follow:    o.Follow,
		Whitelist: o.whitelist,
		Blacklist: o.blacklist,
		Since:     nil,
		Until:     nil,
	}

	stream, err := client.GetFlows(context.Background(), req)
	if err != nil {
		err = errors.New("failed to connect to the gRPC server\nPossible troubleshooting:\n- Check if Hubble relay is running\n- Create a portforward to hubble relay service using\n\t\033[1maccuknox port-forward cilium\033[0m")
		return err
	}
	for {
		select {
		case <-StopChan:
			return nil

		default:
			res, err := stream.Recv()

			if err == io.EOF {
				return nil
			}

			if err != nil {
				return err
			}

			switch res.ResponseTypes.(type) {
			case *observer.GetFlowsResponse_Flow:
				fd := WriteProtoFlow(res)
				_, err := fmt.Fprintf(os.Stdout,
					"%s%s: %s %s %s %s %s %s %s \n",
					fd.Timestamp,
					fd.Node,
					fd.Source,
					fd.SourceIdentity,
					fd.Arrow,
					fd.Destination,
					fd.DestinationIdentity,
					fd.FlowType,
					fd.Verdict)

				if err != nil {
					return fmt.Errorf("failed to write out packet: %v", err)
				}
			}
		}
	}

}
