package network

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/cilium/cilium/api/v1/observer"
	"google.golang.org/grpc"
)

// Options Structure
type Options struct {
	Follow bool
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
		Whitelist: nil,
		Blacklist: nil,
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
