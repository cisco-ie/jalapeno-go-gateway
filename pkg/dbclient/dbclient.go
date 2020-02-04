package dbclient

import (
	"context"
	"sync"
	"time"

	pbapi "github.com/cisco-ie/jalapeno-go-gateway/pkg/apis"
	"github.com/golang/glog"
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

var (
	// maxDBRequestTimeout defines a maximum wait time for db worker to retrieve QoE information
	maxDBRequestTimeout = time.Millisecond * 2000
)

// DB defines methods to access database
type DB interface {
	GetQoE(context.Context, *pbapi.Qoe, chan *pbapi.Qoe)
	GetVPN(context.Context, *VPNRequest, chan *VPNReply)
}

// DBClient defines public method for gRPC server to handle QoE related requests
type DBClient interface {
	GetQoE(context.Context, *pbapi.RequestQoE, chan *pbapi.ResponseQoE)
	GetVPN(context.Context, *VPNRequest, chan *VPNReply)
}

// dbClient defines the database client, it stores the databse interface,
// the database interface can be either a real database or a mock database used in unit testing.
type dbClient struct {
	db DB
}

// GetVPN defines a function to request from DB VPN label based on RD and RTs passed by the client.
func (dbc *dbClient) GetVPN(ctx context.Context, req *VPNRequest, result chan *VPNReply) {
	glog.Infof("dbclient GetVPN rd: %+v ", req.RD)
	// Initializing reply
	var repl *VPNReply
	ch := make(chan *VPNReply)
	// Starting DB's Get as a go routine and wait either for a result
	// received from ch channel or a context timeout event.
	go dbc.db.GetVPN(ctx, req, ch)
	// Wait for all go routine either to complete or the context to timeout
	for {
		select {
		case repl = <-ch:
			result <- repl
			glog.Infof("dbclient GetVPN reply: %+v ", repl.Label)
			return
		case <-ctx.Done():
			// Context was canceled, returning to prevent go routine leaking.
			result <- &VPNReply{}
			glog.Infof("dbclient GetVPN reply with empty reply, error: %+v ", ctx.Err())
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
	glog.Infof("Sending %d QoE back to gRPC server", len(replQoEs.Qoes))
	result <- &replQoEs
}

// NewDBClient return  a new instance of a DB client
func NewDBClient(db DB) DBClient {
	return &dbClient{
		db: db,
	}
}
