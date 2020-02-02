package gateway

import (
	"context"
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
	var rdValue ptypes.DynamicAny
	if err := ptypes.UnmarshalAny(reqVPN.Rd, &rdValue); err != nil {
		return fmt.Errorf("failed to unmarshal route distinguisher with error: %+v", err)
	}
	rd := "RD: "
	switch v := rdValue.Message.(type) {
	case *pbapi.RouteDistinguisherTwoOctetAS:
		rd = fmt.Sprintf("%d:%d ", uint16(v.Admin), v.Assigned)
	case *pbapi.RouteDistinguisherIPAddress:
		rd = fmt.Sprintf("%s:%d ", v.Admin, uint16(v.Assigned))
	case *pbapi.RouteDistinguisherFourOctetAS:
		rd = fmt.Sprintf("%d:%d ", v.Admin, uint16(v.Assigned))
	}
	var rt string
	for i := 0; i < len(reqVPN.Rt); i++ {
		var rtValue ptypes.DynamicAny
		if err := ptypes.UnmarshalAny(reqVPN.Rt[i], &rtValue); err != nil {
			return fmt.Errorf("failed to unmarshal route target with error: %+v", err)
		}
		rt += "RT: "
		switch v := rtValue.Message.(type) {
		case *pbapi.TwoOctetAsSpecificExtended:
			rt += fmt.Sprintf("Subtype:%d 2 bytes AS:%d LocalAdmin:%d Transitive:%t", v.SubType, uint16(v.As), v.LocalAdmin, v.IsTransitive)
		case *pbapi.IPv4AddressSpecificExtended:
			rt += fmt.Sprintf("Subtype:%d Address:%d LocalAdmin:%d Transitive:%t", v.SubType, v.Address, uint16(v.LocalAdmin), v.IsTransitive)
		case *pbapi.FourOctetAsSpecificExtended:
			rt += fmt.Sprintf("Subtype:%d 4 bytes AS:%d LocalAdmin:%d Transitive:%t", v.SubType, v.As, uint16(v.LocalAdmin), v.IsTransitive)
		}
	}
	glog.Infof("Request for VPN with %s %s", rd, rt)
	if err := stream.Send(&pbapi.ResponseVPNEntry{}); err != nil {
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
