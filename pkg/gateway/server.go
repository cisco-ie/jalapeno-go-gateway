package gateway

import (
	"context"
	"net"
	"time"

	pbapi "github.com/cisco-ie/jalapeno-go-gateway/pkg/apis"
	"github.com/golang/glog"
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
}

func (g *gateway) Start() {
	glog.Infof("Starting Gateway's gRPC on %s\n", g.conn.Addr().String())
	go g.gSrv.Serve(g.conn)
}

func (g *gateway) Stop() {
	glog.Infof("Stopping Gateway's gRPC server...")
	g.gSrv.Stop()
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
func NewGateway(conn net.Listener) Gateway {
	gSrv := gateway{
		conn: conn,
		gSrv: grpc.NewServer([]grpc.ServerOption{}...),
	}
	pbapi.RegisterQoEServiceServer(gSrv.gSrv, &gSrv)

	return &gSrv

}

// processQoERequest start DB client and wait for either of 2 events, result comming back from a result channel
// or a context timing out.
func (g *gateway) processQoERequest(ctx context.Context, reqQoes *pbapi.RequestQoE) (*pbapi.ResponseQoE, error) {
	var replQoes pbapi.ResponseQoE
	result := make(chan pbapi.ResponseQoE)
	// TODO Once DB Interface is built, switch to using it here.
	go dbClientSim(reqQoes, result)
	select {
	case replQoes = <-result:
		return &replQoes, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func dbClientSim(reqQoes *pbapi.RequestQoE, result chan pbapi.ResponseQoE) {
	replQoes := pbapi.ResponseQoE{
		Qoes: reqQoes.Qoes,
	}
	// Simulating hung process
	// time.Sleep(maxRequestProcessTime)

	result <- replQoes
}
