package bgpclient

import (
	"context"
	"fmt"
	"io"

	pbapi "github.com/cisco-ie/jalapeno-go-gateway/pkg/apis"
	"github.com/gogo/protobuf/proto"
	"github.com/golang/glog"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	api "github.com/osrg/gobgp/api"
	"github.com/osrg/gobgp/pkg/packet/bgp"
	"google.golang.org/grpc"
)

type VPNRequest struct {
	RD bgp.RouteDistinguisherInterface
	RT []bgp.ExtendedCommunityInterface
}

// VPNReply carries values of labels along with corresponding to the label RD and RTs
type VPNReply struct {
	RD    bgp.RouteDistinguisherInterface
	RT    []bgp.ExtendedCommunityInterface
	Label uint32
}

// BGPClient defines the interface for communication with gobgpd process
type BGPClient interface {
	GetPrefix() error
	AddPrefix() error
}

type bgpclient struct {
	conn *grpc.ClientConn
	bgp  api.GobgpApiClient
}

// NewBGPClient creates a new instance of BGP client, addr variable carries
// the address of gobgp process, it could in ip:port or dns-name:port formats.
func NewBGPClient(addr string) (BGPClient, error) {
	conn, err := grpc.DialContext(context.TODO(), addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	client := api.NewGobgpApiClient(conn)
	// Testing connection to gobgp by requesting its global config
	if _, err := client.GetBgp(context.TODO(), &api.GetBgpRequest{}); err != nil {
		return nil, err
	}
	return &bgpclient{
		conn: conn,
		bgp:  client,
	}, nil
}

func (c *bgpclient) Stop() {
	c.conn.Close()
}

// NewVPNRequest creates VPNRequest struct and populates its fields with information
// came in gRPC request.
func NewVPNRequest(reqVPN *pbapi.RequestVPN) (*VPNRequest, error) {
	var err error
	vpnRequest := VPNRequest{}
	// Processing RD found  in the VPN request
	if vpnRequest.RD, err = UnmarshalRD(reqVPN.Rd); err != nil {
		return nil, err
	}
	// Processing all RTs found in VPN reuqest
	if vpnRequest.RT, err = UnmarshalRT(reqVPN.Rt); err != nil {
		return nil, err
	}
	glog.V(5).Infof("vpn request RD: %+v RT: %+v", vpnRequest.RD, vpnRequest.RT)

	return &vpnRequest, nil
}

// NewVPNReply creates pbapi.ResponseVPNEntry struct and populates its fields with information
// came in VPNReply.
func NewVPNReply(repl *VPNReply) (*pbapi.ResponseVPNEntry, error) {
	gRepl := pbapi.ResponseVPNEntry{
		Label: repl.Label,
	}
	// Marshalling RD from VPNReply
	gRepl.Rd = MarshalRD(repl.RD)
	// Marshalling RTs from VPNReply
	gRepl.Rt = append(gRepl.Rt, MarshalRTs(repl.RT)...)

	return &gRepl, nil
}

func UnmarshalRD(rd *any.Any) (bgp.RouteDistinguisherInterface, error) {
	var rdValue ptypes.DynamicAny
	if err := ptypes.UnmarshalAny(rd, &rdValue); err != nil {
		return nil, fmt.Errorf("failed to unmarshal route distinguisher with error: %+v", err)
	}

	switch v := rdValue.Message.(type) {
	case *api.RouteDistinguisherTwoOctetAS:
		return bgp.NewRouteDistinguisherTwoOctetAS(uint16(v.Admin), v.Assigned), nil
	case *api.RouteDistinguisherIPAddress:
		return bgp.NewRouteDistinguisherIPAddressAS(v.Admin, uint16(v.Assigned)), nil
	case *api.RouteDistinguisherFourOctetAS:
		return bgp.NewRouteDistinguisherFourOctetAS(v.Admin, uint16(v.Assigned)), nil
	default:
		return nil, fmt.Errorf("Unknown route distinguisher type: %+v", v)
	}
}

func UnmarshalRT(rts []*any.Any) ([]bgp.ExtendedCommunityInterface, error) {
	repl := make([]bgp.ExtendedCommunityInterface, 0)
	for i := 0; i < len(rts); i++ {
		var rtValue ptypes.DynamicAny
		if err := ptypes.UnmarshalAny(rts[i], &rtValue); err != nil {
			return nil, fmt.Errorf("failed to unmarshal route target with error: %+v", err)
		}
		switch v := rtValue.Message.(type) {
		case *api.TwoOctetAsSpecificExtended:
			repl = append(repl, bgp.NewTwoOctetAsSpecificExtended(bgp.ExtendedCommunityAttrSubType(v.SubType), uint16(v.As), v.LocalAdmin, v.IsTransitive))
		case *api.IPv4AddressSpecificExtended:
			repl = append(repl, bgp.NewIPv4AddressSpecificExtended(bgp.ExtendedCommunityAttrSubType(v.SubType), v.Address, uint16(v.LocalAdmin), v.IsTransitive))
		case *api.FourOctetAsSpecificExtended:
			repl = append(repl, bgp.NewFourOctetAsSpecificExtended(bgp.ExtendedCommunityAttrSubType(v.SubType), v.As, uint16(v.LocalAdmin), v.IsTransitive))
		default:
			return nil, fmt.Errorf("Unknown route target type: %+v", v)
		}
	}

	return repl, nil
}

func MarshalRD(rd bgp.RouteDistinguisherInterface) *any.Any {
	var r proto.Message
	switch v := rd.(type) {
	case *bgp.RouteDistinguisherTwoOctetAS:
		glog.Infof("Admin: %+v Assigned: %+v", v.Admin, v.Assigned)
		r = &api.RouteDistinguisherTwoOctetAS{
			Admin:    uint32(v.Admin),
			Assigned: v.Assigned,
		}
	case *bgp.RouteDistinguisherIPAddressAS:
		glog.Infof("Admin: %+v Assigned: %+v", v.Admin, v.Assigned)
		r = &api.RouteDistinguisherIPAddress{
			Admin:    v.Admin.String(),
			Assigned: uint32(v.Assigned),
		}
	case *bgp.RouteDistinguisherFourOctetAS:
		glog.Infof("Admin: %+v Assigned: %+v", v.Admin, v.Assigned)
		r = &api.RouteDistinguisherFourOctetAS{
			Admin:    v.Admin,
			Assigned: uint32(v.Assigned),
		}
	default:
		glog.Infof("Unknown type: %+v", v)
		return nil
	}
	a, _ := ptypes.MarshalAny(r)
	return a
}

func MarshalRT(rt bgp.ExtendedCommunityInterface) *any.Any {
	var r proto.Message
	switch v := rt.(type) {
	case *bgp.TwoOctetAsSpecificExtended:
		r = &api.TwoOctetAsSpecificExtended{
			IsTransitive: true,
			SubType:      uint32(bgp.EC_SUBTYPE_ROUTE_TARGET),
			As:           uint32(v.AS),
			LocalAdmin:   uint32(v.LocalAdmin),
		}
	case *bgp.IPv4AddressSpecificExtended:
		r = &api.IPv4AddressSpecificExtended{
			IsTransitive: true,
			SubType:      uint32(bgp.EC_SUBTYPE_ROUTE_TARGET),
			Address:      v.IPv4.String(),
			LocalAdmin:   uint32(v.LocalAdmin),
		}
	case *bgp.FourOctetAsSpecificExtended:
		r = &api.FourOctetAsSpecificExtended{
			IsTransitive: true,
			SubType:      uint32(bgp.EC_SUBTYPE_ROUTE_TARGET),
			As:           uint32(v.AS),
			LocalAdmin:   uint32(v.LocalAdmin),
		}

	default:
		glog.Infof("Marshal RT Unknown type: %+v", v)
		return nil
	}
	a, _ := ptypes.MarshalAny(r)
	return a
}

func MarshalRTs(rts []bgp.ExtendedCommunityInterface) []*any.Any {
	a := make([]*any.Any, len(rts))
	for i := 0; i < len(rts); i++ {
		a[i] = MarshalRT(rts[i])
	}

	return a
}

func (c *bgpclient) GetPrefix() error {
	return c.getPrefix()
}

func (c *bgpclient) getPrefix() error {
	list, err := c.bgp.ListPath(context.TODO(), &api.ListPathRequest{
		TableType: api.TableType_GLOBAL,
		Family:    &api.Family{Afi: api.Family_AFI_IP, Safi: api.Family_SAFI_UNICAST},
		Prefixes: []*api.TableLookupPrefix{
			{
				Prefix:       "2.2.2.2/32",
				LookupOption: api.TableLookupOption_LOOKUP_EXACT,
			},
		},
	})
	if err != nil {
		return err
	}
	for {
		entry, err := list.Recv()
		if err == io.EOF {
			glog.Info("Received EOF message.")
			return nil
		}
		if err != nil {
			glog.Errorf("failed to receive a message from the stream with error: %+v", err)
			return err
		}

		if entry != nil {
			glog.Infof("Received message: %+v", *entry.Destination)
		} else {
			glog.Info("Received empty message.")
		}
	}
}

func (c *bgpclient) AddPrefix() error {
	return c.addPrefix()
}

func (c *bgpclient) addPrefix() error {
	rd := bgp.NewRouteDistinguisherTwoOctetAS(577, 65002)
	nlrivpn, _ := ptypes.MarshalAny(&api.LabeledVPNIPAddressPrefix{
		Labels:    []uint32{2000},
		Rd:        MarshalRD(rd),
		PrefixLen: 32,
		Prefix:    "5.5.5.5",
	})

	a1, _ := ptypes.MarshalAny(&api.OriginAttribute{
		Origin: 0,
	})
	a2, _ := ptypes.MarshalAny(&api.NextHopAttribute{
		NextHop: "9.9.9.9",
	})
	a3, _ := ptypes.MarshalAny(&api.ExtendedCommunitiesAttribute{
		Communities: []*any.Any{MarshalRT(bgp.NewTwoOctetAsSpecificExtended(bgp.EC_SUBTYPE_ROUTE_TARGET, 577, 65002, false))},
	})
	attrs := []*any.Any{a1, a2, a3}
	_, err := c.bgp.AddPath(context.TODO(), &api.AddPathRequest{
		TableType: api.TableType_GLOBAL,
		Path: &api.Path{
			Nlri:      nlrivpn,
			Pattrs:    attrs,
			Family:    &api.Family{Afi: api.Family_AFI_IP, Safi: api.Family_SAFI_MPLS_VPN},
			Best:      true,
			SourceAsn: 65070,
		},
	})
	if err != nil {
		return err
	}

	nlrilu, _ := ptypes.MarshalAny(&api.LabeledIPAddressPrefix{
		Labels:    []uint32{3000},
		PrefixLen: 32,
		Prefix:    "9.9.9.9",
	})
	a2, _ = ptypes.MarshalAny(&api.NextHopAttribute{
		NextHop: "10.0.0.1",
	})
	attrs = []*any.Any{a1, a2}
	_, err = c.bgp.AddPath(context.TODO(), &api.AddPathRequest{
		TableType: api.TableType_GLOBAL,
		Path: &api.Path{
			Nlri:   nlrilu,
			Pattrs: attrs,
			Family: &api.Family{Afi: api.Family_AFI_IP, Safi: api.Family_SAFI_MPLS_LABEL},
			Best:   true,
		},
	})
	if err != nil {
		return err
	}

	return nil
}
