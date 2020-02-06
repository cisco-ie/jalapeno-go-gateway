package bgpclient

import (
	"context"
	"encoding/binary"
	"fmt"

	pbapi "github.com/cisco-ie/jalapeno-go-gateway/pkg/apis"
	"github.com/golang/glog"
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
	RD *RDValue
	RT []RTValue
}

// VPNReply carries values of labels along with corresponding to the label RD and RTs
type VPNReply struct {
	RD    *RDValue
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
	var err error
	vpnRequest := VPNRequest{}
	// Processing RD found  in the VPN request
	if vpnRequest.RD, err = unmarshalRD(reqVPN.Rd); err != nil {
		return nil, err
	}
	// Processing all RTs found in VPN reuqest
	if vpnRequest.RT, err = unmarshalRT(reqVPN.Rt); err != nil {
		return nil, err
	}
	glog.V(5).Infof("vpn request RD: %+v RT: %+v", *vpnRequest.RD, vpnRequest.RT)

	return &vpnRequest, nil
}

// NewVPNReply creates pbapi.ResponseVPNEntry struct and populates its fields with information
// came in VPNReply.
func NewVPNReply(repl *VPNReply) (*pbapi.ResponseVPNEntry, error) {
	gRepl := pbapi.ResponseVPNEntry{
		Label: repl.Label,
	}
	var err error
	// Marshalling RD from VPNReply
	if gRepl.Rd, err = marshalRD(repl.RD); err != nil {
		return nil, err
	}
	// Marshalling RTs from VPNReply
	if gRepl.Rt, err = marshalRT(repl.RT); err != nil {
		return nil, err
	}

	return &gRepl, nil
}

func unmarshalRD(rd *any.Any) (*RDValue, error) {
	var rdValue ptypes.DynamicAny
	if err := ptypes.UnmarshalAny(rd, &rdValue); err != nil {
		return nil, fmt.Errorf("failed to unmarshal route distinguisher with error: %+v", err)
	}
	repl := &RDValue{
		Value: [8]byte{},
	}

	switch v := rdValue.Message.(type) {
	case *pbapi.RouteDistinguisherTwoOctetAS:
		binary.BigEndian.PutUint16(repl.Value[0:], uint16(v.Admin))
		binary.BigEndian.PutUint32(repl.Value[4:], v.Assigned)
		repl.T = RouteDistinguisherTwoOctetAS
	case *pbapi.RouteDistinguisherIPAddress:
		copy(repl.Value[0:], v.Admin)
		binary.BigEndian.PutUint32(repl.Value[4:], v.Assigned)
		repl.T = RouteDistinguisherIPAddressAS
	case *pbapi.RouteDistinguisherFourOctetAS:
		binary.BigEndian.PutUint32(repl.Value[0:], v.Admin)
		binary.BigEndian.PutUint16(repl.Value[4:], uint16(v.Assigned))
		repl.T = RouteDistinguisherFourOctetAS
	}

	return repl, nil
}

func unmarshalRT(rts []*any.Any) ([]RTValue, error) {
	repl := make([]RTValue, 0)
	for i := 0; i < len(rts); i++ {
		var rtValue ptypes.DynamicAny
		if err := ptypes.UnmarshalAny(rts[i], &rtValue); err != nil {
			return nil, fmt.Errorf("failed to unmarshal route target with error: %+v", err)
		}
		rt := RTValue{
			Value: [8]byte{},
		}
		switch v := rtValue.Message.(type) {
		case *pbapi.TwoOctetAsSpecificExtended:
			binary.BigEndian.PutUint16(rt.Value[0:], uint16(v.SubType))
			binary.BigEndian.PutUint16(rt.Value[2:], uint16(v.As))
			binary.BigEndian.PutUint32(rt.Value[4:], v.LocalAdmin)
			rt.T = TwoOctetAsSpecificExtended
		case *pbapi.IPv4AddressSpecificExtended:
			binary.BigEndian.PutUint16(rt.Value[0:], uint16(v.SubType))
			copy(rt.Value[2:], []byte(v.Address))
			binary.BigEndian.PutUint16(rt.Value[6:], uint16(v.LocalAdmin))
			rt.T = IPv4AddressSpecificExtended
		case *pbapi.FourOctetAsSpecificExtended:
			binary.BigEndian.PutUint16(rt.Value[0:], uint16(v.SubType))
			binary.BigEndian.PutUint32(rt.Value[2:], v.As)
			binary.BigEndian.PutUint16(rt.Value[6:], uint16(v.LocalAdmin))
			rt.T = FourOctetAsSpecificExtended
		}
		repl = append(repl, rt)
	}

	return repl, nil
}

func marshalRD(rd *RDValue) (*any.Any, error) {
	repl := &any.Any{}
	var err error
	if rd == nil {
		return nil, nil
	}
	switch rd.T {
	case RouteDistinguisherTwoOctetAS:
		repl, err = ptypes.MarshalAny(&pbapi.RouteDistinguisherTwoOctetAS{
			Admin:    binary.BigEndian.Uint32(rd.Value[0:]),
			Assigned: binary.BigEndian.Uint32(rd.Value[4:]),
		})
		if err != nil {
			return nil, err
		}
	case RouteDistinguisherIPAddressAS:
		repl, err = ptypes.MarshalAny(&pbapi.RouteDistinguisherIPAddress{
			Admin:    string(rd.Value[0:4]),
			Assigned: binary.BigEndian.Uint32(rd.Value[4:]),
		})
		if err != nil {
			return nil, err
		}
	case RouteDistinguisherFourOctetAS:
		repl, err = ptypes.MarshalAny(&pbapi.RouteDistinguisherFourOctetAS{
			Admin:    binary.BigEndian.Uint32(rd.Value[0:]),
			Assigned: binary.BigEndian.Uint32(rd.Value[4:]),
		})
		if err != nil {
			return nil, err
		}
	}

	return repl, nil
}

func marshalRT(rts []RTValue) ([]*any.Any, error) {
	repl := make([]*any.Any, 0)
	var err error
	var rt *any.Any

	for i := 0; i < len(rts); i++ {
		switch rts[i].T {
		case TwoOctetAsSpecificExtended:
			rt, err = ptypes.MarshalAny(&pbapi.TwoOctetAsSpecificExtended{
				SubType:    binary.BigEndian.Uint32(rts[i].Value[0:]),
				As:         binary.BigEndian.Uint32(rts[i].Value[2:]),
				LocalAdmin: binary.BigEndian.Uint32(rts[i].Value[4:]),
			})
			if err != nil {
				return nil, err
			}
		case IPv4AddressSpecificExtended:
			rt, err = ptypes.MarshalAny(&pbapi.IPv4AddressSpecificExtended{
				SubType:    binary.BigEndian.Uint32(rts[i].Value[0:]),
				Address:    string(rts[i].Value[2:6]),
				LocalAdmin: binary.BigEndian.Uint32(repl[i].Value[6:]),
			})
			if err != nil {
				return nil, err
			}
		case FourOctetAsSpecificExtended:
			rt, err = ptypes.MarshalAny(&pbapi.FourOctetAsSpecificExtended{
				SubType:    binary.BigEndian.Uint32(rts[i].Value[0:]),
				As:         binary.BigEndian.Uint32(rts[i].Value[2:]),
				LocalAdmin: binary.BigEndian.Uint32(rts[i].Value[6:]),
			})
			if err != nil {
				return nil, err
			}
		}
		repl = append(repl, rt)
	}

	return repl, nil
}
