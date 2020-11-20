package gateway

import (
	"context"
	"net"
	"time"

	pbapi "github.com/cisco-ie/jalapeno-go-gateway/pkg/apis"
	"github.com/cisco-ie/jalapeno-go-gateway/pkg/bgpclient"
	"github.com/cisco-ie/jalapeno-go-gateway/pkg/dbclient"
	"github.com/golang/glog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
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
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		client := md.Get("CLIENT_IP")
		if len(client) != 0 {
			glog.Infof("VPN request from client: %+v", client)
		}
	}
	// Building new vpn reuqest from info recieved in grpc packet
	vpnRequest, err := bgpclient.NewVPNRequest(reqVPN)
	if err != nil {
		return err
	}
	glog.Infof("request for Route Distinguisher: %+v", vpnRequest.RD.String())
	repl, err := g.processVPNRequest(ctx, vpnRequest)
	if err != nil {
		return err
	}
	r, err := bgpclient.NewVPNReply(repl)
	if err != nil {
		return err
	}
	if err := stream.Send(r); err != nil {
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
func (g *gateway) processVPNRequest(ctx context.Context, req *bgpclient.VPNRequest) (*bgpclient.VPNReply, error) {
	glog.V(5).Infof("processing VPN Request for RD: %+v RT: %+v", req.RD.String(), req.RT)
	var repl *bgpclient.VPNReply
	result := make(chan *bgpclient.VPNReply)
	// Requesting DB client to retrieve requested infotmation
	go g.dbc.GetVPN(ctx, req, result)
	select {
	case repl = <-result:
		if repl != nil && repl.RD != nil {
			glog.V(5).Infof("DB client returned RD: %+v RT: %+v Label: %d", repl.RD.String(), repl.RT, repl.Label)
		} else {
			glog.V(5).Infof("DB client returned nil")
		}
		return repl, nil
	case <-ctx.Done():
		glog.V(5).Infof("DB client returned error: %+v ", ctx.Err())
		return nil, ctx.Err()
	}
}
