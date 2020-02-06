package bgpclient

import (
	"context"
	"encoding/binary"
	"fmt"

	pbapi "github.com/cisco-ie/jalapeno-go-gateway/pkg/apis"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	api "github.com/osrg/gobgp/api"
	"google.golang.org/grpc"
)

// RDType defines type of Route Distinguisher
type RDType uint32

const (
	RouteDistinguisherTwoOctetAS RDType = iota + 1
	RouteDistinguisherIPAddressAS
	RouteDistinguisherFourOctetAS
)

// RTType defines type of Route Target
type RTType uint32

const (
	TwoOctetAsSpecificExtended RTType = iota + 1
	IPv4AddressSpecificExtended
	FourOctetAsSpecificExtended
)

// RDValue defines an easier to consume form of RD. T carries RD type, and Value carries a slice
// of bytes representing RD
type RDValue struct {
	T     RDType
	Value [8]byte
}

// RTValue defines an easier to consume form of RT. T carries RT type, and Value carries a slice
// of bytes representing a single Route Target.
type RTValue struct {
	T     RTType
	Value [8]byte
}

// VPNRequest puts together parameters (RD and RTs) used in search of a matching VPN label
type VPNRequest struct {
	RD RDValue
	RT []RTValue
}

// VPNReply carries values of labels along with corresponding to the label RD and RTs
type VPNReply struct {
	RD    RDValue
	RT    []RTValue
	Label uint32
}

// BGPClient defines the interface for communication with gobgpd process
type BGPClient interface {
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

func (b *bgpclient) Stop() {
	b.conn.Close()
}

// NewVPNRequest creates VPNRequest struct and populates its fields with information
// came in gRPC request.
func NewVPNRequest(reqVPN *pbapi.RequestVPN) (*VPNRequest, error) {
	var rdValue ptypes.DynamicAny
	if err := ptypes.UnmarshalAny(reqVPN.Rd, &rdValue); err != nil {
		return nil, fmt.Errorf("failed to unmarshal route distinguisher with error: %+v", err)
	}
	vpnRequest := VPNRequest{
		RD: RDValue{
			Value: [8]byte{},
		},
		RT: make([]RTValue, len(reqVPN.Rt)),
	}
	// Processing RD found  in the VPN request
	switch v := rdValue.Message.(type) {
	case *pbapi.RouteDistinguisherTwoOctetAS:
		binary.BigEndian.PutUint16(vpnRequest.RD.Value[0:], uint16(v.Admin))
		binary.BigEndian.PutUint32(vpnRequest.RD.Value[4:], v.Assigned)
		vpnRequest.RD.T = RouteDistinguisherTwoOctetAS
	case *pbapi.RouteDistinguisherIPAddress:
		copy(vpnRequest.RD.Value[0:], v.Admin)
		binary.BigEndian.PutUint32(vpnRequest.RD.Value[4:], v.Assigned)
		vpnRequest.RD.T = RouteDistinguisherIPAddressAS
	case *pbapi.RouteDistinguisherFourOctetAS:
		binary.BigEndian.PutUint32(vpnRequest.RD.Value[0:], v.Admin)
		binary.BigEndian.PutUint16(vpnRequest.RD.Value[4:], uint16(v.Assigned))
		vpnRequest.RD.T = RouteDistinguisherFourOctetAS
	}
	// Processing all RTs found in VPN reuqest
	for i := 0; i < len(reqVPN.Rt); i++ {
		var rtValue ptypes.DynamicAny
		if err := ptypes.UnmarshalAny(reqVPN.Rt[i], &rtValue); err != nil {
			return nil, fmt.Errorf("failed to unmarshal route target with error: %+v", err)
		}
		vpnRequest.RT[i].Value = [8]byte{}
		switch v := rtValue.Message.(type) {
		case *pbapi.TwoOctetAsSpecificExtended:
			binary.BigEndian.PutUint16(vpnRequest.RT[i].Value[0:], uint16(v.SubType))
			binary.BigEndian.PutUint16(vpnRequest.RT[i].Value[2:], uint16(v.As))
			binary.BigEndian.PutUint32(vpnRequest.RT[i].Value[4:], v.LocalAdmin)
			vpnRequest.RT[i].T = TwoOctetAsSpecificExtended
		case *pbapi.IPv4AddressSpecificExtended:
			binary.BigEndian.PutUint16(vpnRequest.RT[i].Value[0:], uint16(v.SubType))
			copy(vpnRequest.RT[i].Value[2:], []byte(v.Address))
			binary.BigEndian.PutUint16(vpnRequest.RT[i].Value[6:], uint16(v.LocalAdmin))
			vpnRequest.RT[i].T = IPv4AddressSpecificExtended
		case *pbapi.FourOctetAsSpecificExtended:
			binary.BigEndian.PutUint16(vpnRequest.RT[i].Value[0:], uint16(v.SubType))
			binary.BigEndian.PutUint32(vpnRequest.RT[i].Value[2:], v.As)
			binary.BigEndian.PutUint16(vpnRequest.RT[i].Value[6:], uint16(v.LocalAdmin))
			vpnRequest.RT[i].T = FourOctetAsSpecificExtended
		}
	}

	return &vpnRequest, nil
}

// NewVPNReply creates pbapi.ResponseVPNEntry struct and populates its fields with information
// came in VPNReply.
func NewVPNReply(repl *VPNReply) (*pbapi.ResponseVPNEntry, error) {
	gRepl := pbapi.ResponseVPNEntry{
		Label: repl.Label,
	}
	var err error

	switch repl.RD.T {
	case RouteDistinguisherTwoOctetAS:
		gRepl.Rd, err = ptypes.MarshalAny(&pbapi.RouteDistinguisherTwoOctetAS{
			Admin:    binary.BigEndian.Uint32(repl.RD.Value[0:]),
			Assigned: binary.BigEndian.Uint32(repl.RD.Value[4:]),
		})
		if err != nil {
			return nil, err
		}
	case RouteDistinguisherIPAddressAS:
		gRepl.Rd, err = ptypes.MarshalAny(&pbapi.RouteDistinguisherIPAddress{
			Admin:    string(repl.RD.Value[0:4]),
			Assigned: binary.BigEndian.Uint32(repl.RD.Value[4:]),
		})
		if err != nil {
			return nil, err
		}
	case RouteDistinguisherFourOctetAS:
		gRepl.Rd, err = ptypes.MarshalAny(&pbapi.RouteDistinguisherFourOctetAS{
			Admin:    binary.BigEndian.Uint32(repl.RD.Value[0:]),
			Assigned: binary.BigEndian.Uint32(repl.RD.Value[4:]),
		})
		if err != nil {
			return nil, err
		}
	}

	// Processing all RTs found in VPN reuqest
	gRepl.Rt = make([]*any.Any, 0)
	var a *any.Any
	for i := 0; i < len(repl.RT); i++ {
		switch repl.RT[i].T {
		case TwoOctetAsSpecificExtended:
			a, err = ptypes.MarshalAny(&pbapi.TwoOctetAsSpecificExtended{
				SubType:    binary.BigEndian.Uint32(repl.RT[i].Value[0:]),
				As:         binary.BigEndian.Uint32(repl.RT[i].Value[2:]),
				LocalAdmin: binary.BigEndian.Uint32(repl.RT[i].Value[4:]),
			})
			if err != nil {
				return nil, err
			}
		case IPv4AddressSpecificExtended:
			a, err = ptypes.MarshalAny(&pbapi.IPv4AddressSpecificExtended{
				SubType:    binary.BigEndian.Uint32(repl.RT[i].Value[0:]),
				Address:    string(repl.RT[i].Value[2:6]),
				LocalAdmin: binary.BigEndian.Uint32(repl.RT[i].Value[6:]),
			})
			if err != nil {
				return nil, err
			}
		case FourOctetAsSpecificExtended:
			a, err = ptypes.MarshalAny(&pbapi.FourOctetAsSpecificExtended{
				SubType:    binary.BigEndian.Uint32(repl.RT[i].Value[0:]),
				As:         binary.BigEndian.Uint32(repl.RT[i].Value[2:]),
				LocalAdmin: binary.BigEndian.Uint32(repl.RT[i].Value[6:]),
			})
			if err != nil {
				return nil, err
			}
		}
		gRepl.Rt = append(gRepl.Rt, a)
	}

	return &gRepl, nil
}
