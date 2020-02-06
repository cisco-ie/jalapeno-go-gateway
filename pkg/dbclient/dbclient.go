package dbclient

import (
	"context"
	"sync"
	"time"

	pbapi "github.com/cisco-ie/jalapeno-go-gateway/pkg/apis"
	"github.com/cisco-ie/jalapeno-go-gateway/pkg/bgpclient"
	"github.com/golang/glog"
)

var (
	// maxDBRequestTimeout defines a maximum wait time for db worker to retrieve QoE information
	maxDBRequestTimeout = time.Millisecond * 2000
)

// DB defines methods to access database
type DB interface {
	GetQoE(context.Context, *pbapi.Qoe, chan *pbapi.Qoe)
	GetVPN(context.Context, *bgpclient.VPNRequest, chan *bgpclient.VPNReply)
}

// DBClient defines public method for gRPC server to handle QoE related requests
type DBClient interface {
	GetQoE(context.Context, *pbapi.RequestQoE, chan *pbapi.ResponseQoE)
	GetVPN(context.Context, *bgpclient.VPNRequest, chan *bgpclient.VPNReply)
}

// dbClient defines the database client, it stores the databse interface,
// the database interface can be either a real database or a mock database used in unit testing.
type dbClient struct {
	db DB
}

// GetVPN defines a function to request from DB VPN label based on RD and RTs passed by the client.
func (dbc *dbClient) GetVPN(ctx context.Context, req *bgpclient.VPNRequest, result chan *bgpclient.VPNReply) {
	glog.V(5).Infof("DB Client received request for RD: %+v RT: %+v", *req.RD, req.RT)
	// Initializing reply
	var repl *bgpclient.VPNReply
	ch := make(chan *bgpclient.VPNReply)
	// Starting DB's Get as a go routine and wait either for a result
	// received from ch channel or a context timeout event.
	go dbc.db.GetVPN(ctx, req, ch)
	// Wait for all go routine either to complete or the context to timeout
	for {
		select {
		case repl = <-ch:
			result <- repl
			if repl != nil && repl.RD != nil {
				glog.V(5).Infof("Data store returned RD: %+v RT: %+v Label: %d", *repl.RD, repl.RT, repl.Label)
			} else {
				glog.V(5).Infof("Data store returned nil")
			}
			return
		case <-ctx.Done():
			// Context was canceled, returning to prevent go routine leaking.
			result <- &bgpclient.VPNReply{}
			glog.V(5).Infof("Data store returned error: %+v ", ctx.Err())
			return
		}
	}
}

// GetQoE is function called by gRPC client to retrieve QoE related information from the database.
// The request can carry multiple QoE entries. FOr each entry a go routine function is invoked
// which calls DB interface Get.
func (dbc *dbClient) GetQoE(ctx context.Context, reqQoes *pbapi.RequestQoE, result chan *pbapi.ResponseQoE) {
	// Initializing reply
	replQoEs := pbapi.ResponseQoE{}
	replQoEs.Qoes = make(map[int32]*pbapi.Qoe)
	// Wait group will be used to wait until all workers are done
	// In case worker hangs, grpc client's context will destroy it when
	// timeout expires.
	var wg sync.WaitGroup
	for i, req := range reqQoes.Qoes {
		wg.Add(1)
		go func(ctx context.Context, i int32, req *pbapi.Qoe) {
			defer wg.Done()
			// Make channel for DB's Get method to return found data
			ch := make(chan *pbapi.Qoe)
			// Starting DB's Get as a go routine and wait either for a result
			// received from ch channel or a context timeout event.
			go dbc.db.GetQoE(ctx, req, ch)
			for {
				select {
				case repl := <-ch:
					// Result received, storing it in the reply map with the right key
					replQoEs.Qoes[i] = repl
					return
				case <-ctx.Done():
					// Context was canceled, returning to prevent go routine leaking.
					return
				}
			}

		}(ctx, i, req)
	}
	// Wait for all go routine either to complete or the context to timeout
	wg.Wait()
	glog.V(5).Infof("Sending %d QoE back to gRPC server", len(replQoEs.Qoes))
	result <- &replQoEs
}

// NewDBClient return  a new instance of a DB client
func NewDBClient(db DB) DBClient {
	return &dbClient{
		db: db,
	}
}
