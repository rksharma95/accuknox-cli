// Copyright 2019-2021 Authors of Hubble
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

package network

import (
	"fmt"
	"net"
	"path"
	"strconv"
	"strings"

	pb "github.com/cilium/cilium/api/v1/flow"
	observerpb "github.com/cilium/cilium/api/v1/observer"
	identity "github.com/cilium/cilium/pkg/identity"
	api "github.com/cilium/cilium/pkg/monitor/api"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

var c *colorer

type FlowData struct {
	Timestamp           string
	Source              string
	Destination         string
	SourceIdentity      string
	DestinationIdentity string
	Node                string
	FlowType            string
	Verdict             string
	Arrow               string
}

func fmtTimestamp(layout string, ts *timestamppb.Timestamp) string {
	if !ts.IsValid() {
		return "N/A"
	}
	return ts.AsTime().Format(layout)
}

func fmtIdentity(i uint32) string {
	numeric := identity.NumericIdentity(i)
	if numeric.IsReservedIdentity() {
		return c.identity(fmt.Sprintf("(%s)", numeric))
	}

	return c.identity(fmt.Sprintf("(identity:%d)", i))
}

func getSecurityIdentities(f *pb.Flow, fd *FlowData) {
	if f == nil {
		fd.SourceIdentity = ""
		fd.DestinationIdentity = ""
	} else {
		fd.SourceIdentity = fmtIdentity(f.GetSource().GetIdentity())
		fd.DestinationIdentity = fmtIdentity(f.GetDestination().GetIdentity())
	}
}

// WriteProtoFlow function
func WriteProtoFlow(res *observerpb.GetFlowsResponse) *FlowData {
	fd := new(FlowData)
	f := res.GetFlow()
	gethostNames(f, fd)
	getSecurityIdentities(f, fd)

	fd.Node = fmt.Sprintf(" [%s]", f.GetNodeName())

	fd.Arrow = "->"
	if f.GetIsReply() == nil {
		// direction is unknown.
		fd.Arrow = "<>"
	} else if f.GetIsReply().Value {
		// flip the arrow and src/dst for reply packets.
		fd.Source, fd.Destination = fd.Destination, fd.Source
		fd.SourceIdentity, fd.DestinationIdentity = fd.DestinationIdentity, fd.SourceIdentity
		fd.Arrow = "<-"
	}

	fd.Timestamp = fmtTimestamp("2006-01-02 15:04:05.000", f.GetTime())
	fd.FlowType = getFlowType(f)
	fd.Verdict = getVerdict(f)

	return fd
}

func getVerdict(f *pb.Flow) string {
	verdict := f.GetVerdict()
	switch verdict {
	case pb.Verdict_FORWARDED:
		return c.verdictForwarded(verdict.String())
	case pb.Verdict_DROPPED, pb.Verdict_ERROR:
		return c.verdictDropped(verdict.String())
	case pb.Verdict_AUDIT:
		return c.verdictAudit(verdict.String())
	default:
		return verdict.String()
	}
}

func gethostNames(f *pb.Flow, fd *FlowData) {
	var srcNamespace, dstNamespace, srcPodName, dstPodName, srcSvcName, dstSvcName string
	if f == nil {
		fd.Source = ""
		fd.Destination = ""
	}

	if f.GetIP() == nil {
		if eth := f.GetEthernet(); eth != nil {
			fd.Source = c.host(eth.GetSource())
			fd.Destination = c.host(eth.GetDestination())
		} else {
			fd.Source = ""
			fd.Destination = ""
		}
	}

	if src := f.GetSource(); src != nil {
		srcNamespace = src.Namespace
		srcPodName = src.PodName
	}
	if dst := f.GetDestination(); dst != nil {
		dstNamespace = dst.Namespace
		dstPodName = dst.PodName
	}
	if svc := f.GetSourceService(); svc != nil {
		srcNamespace = svc.Namespace
		srcSvcName = svc.Name
	}
	if svc := f.GetDestinationService(); svc != nil {
		dstNamespace = svc.Namespace
		dstSvcName = svc.Name
	}
	srcPort, dstPort := getPorts(f)
	src := hostname(f.GetIP().Source, srcPort, srcNamespace, srcPodName, srcSvcName, f.GetSourceNames())
	dst := hostname(f.GetIP().Destination, dstPort, dstNamespace, dstPodName, dstSvcName, f.GetDestinationNames())
	fd.Source = c.host(src)
	fd.Destination = c.host(dst)
}

func getPorts(f *pb.Flow) (string, string) {
	l4 := f.GetL4()
	if l4 == nil {
		return "", ""
	}
	switch l4.Protocol.(type) {
	case *pb.Layer4_TCP:
		return strconv.Itoa(int(l4.GetTCP().SourcePort)), strconv.Itoa(int(l4.GetTCP().DestinationPort))
	case *pb.Layer4_UDP:
		return strconv.Itoa(int(l4.GetUDP().SourcePort)), strconv.Itoa(int(l4.GetUDP().DestinationPort))
	default:
		return "", ""
	}
}

func hostname(ip, port string, ns, pod, svc string, names []string) (host string) {
	host = ip

	if pod != "" {
		// path.Join omits the slash if ns is empty
		host = path.Join(ns, pod)
	} else if svc != "" {
		host = path.Join(ns, svc)
	} else if len(names) != 0 {
		host = strings.Join(names, ",")
	}

	if port != "" && port != "0" {
		return net.JoinHostPort(host, port)
	}

	return host
}

func getFlowType(f *pb.Flow) string {
	if l7 := f.GetL7(); l7 != nil {
		l7Protocol := "l7"
		l7Type := strings.ToLower(l7.Type.String())
		switch l7.GetRecord().(type) {
		case *pb.Layer7_Http:
			l7Protocol = "http"
		case *pb.Layer7_Dns:
			l7Protocol = "dns"
		case *pb.Layer7_Kafka:
			l7Protocol = "kafka"
		}
		return l7Protocol + "-" + l7Type
	}

	switch f.GetEventType().GetType() {
	case api.MessageTypeTrace:
		return api.TraceObservationPoint(uint8(f.GetEventType().GetSubType()))
	case api.MessageTypeDrop:
		return api.DropReason(uint8(f.GetEventType().GetSubType()))
	case api.MessageTypePolicyVerdict:
		switch f.GetVerdict() {
		case pb.Verdict_FORWARDED, pb.Verdict_AUDIT:
			return api.PolicyMatchType(f.GetPolicyMatchType()).String()
		case pb.Verdict_DROPPED:
			return api.DropReason(uint8(f.GetDropReason()))
		case pb.Verdict_ERROR:
			// ERROR should only happen for L7 events.
		}
	case api.MessageTypeCapture:
		return f.GetDebugCapturePoint().String()
	}

	return "UNKNOWN"
}

func init() {
	c = newColorer()
}
