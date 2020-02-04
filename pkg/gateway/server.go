package gateway

import (
	"context"
	"encoding/binary"
	"fmt"
	"net"
	"time"

	pbapi "github.com/cisco-ie/jalapeno-go-gateway/pkg/apis"
	"github.com/cisco-ie/jalapeno-go-gateway/pkg/bgpclient"
	"github.com/cisco-ie/jalapeno-go-gateway/pkg/dbclient"
	"github.com/golang/glog"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

var (
	// maxRequestProcessTime defines a maximum wait time for a client request processing by DB client
	// and returning results
	maxRequestProcessTime = time.Millisecond * 2000
)

// Gateway defines interface to Gateway gRPC server
type Gateway interface {
	Start()
	Stop()
}
type gateway struct {
	gSrv *grpc.Server
	conn net.Listener
	dbc  dbclient.DBClient
	bgp  bgpclient.BGPClient
}

func (g *gateway) Start() {
	glog.Infof("Starting Gateway's gRPC on %s\n", g.conn.Addr().String())
	go g.gSrv.Serve(g.conn)
}

func (g *gateway) Stop() {
	glog.Infof("Stopping Gateway's gRPC server...")
	g.gSrv.Stop()
}

func (g *gateway) VPN(reqVPN *pbapi.RequestVPN, stream pbapi.GatewayService_VPNServer) error {
	ctx := stream.Context()
	peer, ok := peer.FromContext(ctx)
	if ok {
		glog.V(5).Infof("VPN request from client: %+v", *peer)
	}
	var rdValue ptypes.DynamicAny
	if err := ptypes.UnmarshalAny(reqVPN.Rd, &rdValue); err != nil {
		return fmt.Errorf("failed to unmarshal route distinguisher with error: %+v", err)
	}
	// rdbytes stores route distinguisher from the request in a byte slice form.
	var rdbytes [8]byte
	var rdType dbclient.RDType
	switch v := rdValue.Message.(type) {
	case *pbapi.RouteDistinguisherTwoOctetAS:
		binary.BigEndian.PutUint16(rdbytes[0:], uint16(v.Admin))
		binary.BigEndian.PutUint32(rdbytes[4:], v.Assigned)
		rdType = dbclient.RouteDistinguisherTwoOctetAS
	case *pbapi.RouteDistinguisherIPAddress:
		copy(rdbytes[0:], v.Admin)
		binary.BigEndian.PutUint32(rdbytes[4:], v.Assigned)
		rdType = dbclient.RouteDistinguisherIPAddressAS
	case *pbapi.RouteDistinguisherFourOctetAS:
		binary.BigEndian.PutUint32(rdbytes[0:], v.Admin)
		binary.BigEndian.PutUint16(rdbytes[4:], uint16(v.Assigned))
		rdType = dbclient.RouteDistinguisherFourOctetAS
	}
	// rtbytes stores all route targets received in the request, each route target is represented
	// by a slice of bytes.
	rtbytes := make([][8]byte, len(reqVPN.Rt))
	for i := 0; i < len(reqVPN.Rt); i++ {
		var rtValue ptypes.DynamicAny
		if err := ptypes.UnmarshalAny(reqVPN.Rt[i], &rtValue); err != nil {
			return fmt.Errorf("failed to unmarshal route target with error: %+v", err)
		}
		switch v := rtValue.Message.(type) {
		case *pbapi.TwoOctetAsSpecificExtended:
			binary.BigEndian.PutUint16(rtbytes[i][0:], uint16(v.SubType))
			binary.BigEndian.PutUint16(rtbytes[i][2:], uint16(v.As))
			binary.BigEndian.PutUint32(rtbytes[i][4:], v.LocalAdmin)
		case *pbapi.IPv4AddressSpecificExtended:
			binary.BigEndian.PutUint16(rtbytes[i][0:], uint16(v.SubType))
			copy(rtbytes[i][2:], []byte(v.Address))
			binary.BigEndian.PutUint16(rtbytes[i][6:], uint16(v.LocalAdmin))
		case *pbapi.FourOctetAsSpecificExtended:
			binary.BigEndian.PutUint16(rtbytes[i][0:], uint16(v.SubType))
			binary.BigEndian.PutUint32(rtbytes[i][2:], v.As)
			binary.BigEndian.PutUint16(rtbytes[i][6:], uint16(v.LocalAdmin))
		}
	}
	glog.Infof("Request in bytes, rd: %+v rt: %+v", rdbytes, rtbytes)

	repl, err := g.processVPNRequest(ctx, &dbclient.VPNRequest{
		RD: dbclient.RDValue{
			T:     rdType,
			Value: rdbytes,
		},
	})
	if err != nil {
		return err
	}
	if err := stream.Send(&pbapi.ResponseVPNEntry{
		Rd:    reqVPN.Rd,
		Label: repl.Label,
	}); err != nil {
		return err
	}

	return nil
}

func (g *gateway) QoE(ctx context.Context, reqQoes *pbapi.RequestQoE) (*pbapi.ResponseQoE, error) {
	peer, ok := peer.FromContext(ctx)
	if ok {
		glog.V(5).Infof("QoE request from client: %+v", *peer)
	}
	// Gateway has DB Interface which is used to get requested by a client QoE information.
	// To protect from hung connections, the context passed to DB client interface will be canceled
	// up on reaching a timeout.
	dbctx, cancel := context.WithTimeout(ctx, maxRequestProcessTime)
	defer cancel()
	replQoes, err := g.processQoERequest(dbctx, reqQoes)
	if err != nil {
		return nil, err
	}
	return replQoes, nil
}

// NewGateway return an instance of Gateway interface
func NewGateway(conn net.Listener, dbc dbclient.DBClient, bgp bgpclient.BGPClient) Gateway {
	gSrv := gateway{
		conn: conn,
		gSrv: grpc.NewServer([]grpc.ServerOption{}...),
		dbc:  dbc,
		bgp:  bgp,
	}
	pbapi.RegisterGatewayServiceServer(gSrv.gSrv, &gSrv)

	return &gSrv

}

// processQoERequest start DB client and wait for either of 2 events, result comming back from a result channel
// or a context timing out.
func (g *gateway) processQoERequest(ctx context.Context, reqQoEs *pbapi.RequestQoE) (*pbapi.ResponseQoE, error) {
	var replQoEs *pbapi.ResponseQoE
	result := make(chan *pbapi.ResponseQoE)
	// Requesting DB client to retrieve requested infotmation
	go g.dbc.GetQoE(ctx, reqQoEs, result)
	select {
	case replQoEs = <-result:
		return replQoEs, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// processQoERequest start DB client and wait for either of 2 events, result comming back from a result channel
// or a context timing out.
func (g *gateway) processVPNRequest(ctx context.Context, req *dbclient.VPNRequest) (*dbclient.VPNReply, error) {
	glog.Infof("processVPNRequest rd: %+v ", req.RD)
	var repl *dbclient.VPNReply
	result := make(chan *dbclient.VPNReply)
	// Requesting DB client to retrieve requested infotmation
	go g.dbc.GetVPN(ctx, req, result)
	select {
	case repl = <-result:
		glog.Infof("processVPNRequest reply: %+v ", repl.Label)
		return repl, nil
	case <-ctx.Done():
		glog.Infof("processVPNRequest reply with error: %+v ", ctx.Err())
		return nil, ctx.Err()
	}
}
