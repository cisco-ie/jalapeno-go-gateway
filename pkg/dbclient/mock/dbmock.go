package mock

import (
	"net"
	"sync"
	"time"

	pbapi "github.com/cisco-ie/jalapeno-go-gateway/pkg/apis"
	"github.com/golang/glog"
)

var (
	// maxDBRequestTimeout defines a maximum wait time for db worker to retrieve QoE information
	maxDBRequestTimeout = time.Millisecond * 2000
)

// DBClient defines public methods to handle QoE related information in DB
type DBClient interface {
	GetQoE(*pbapi.RequestQoE, chan pbapi.ResponseQoE)
}

type dbcMock struct {
	mu            sync.Mutex
	store         map[int32]*pbapi.Qoe
	randomDelay   bool
	randomFailure bool
}

func (dbc *dbcMock) worker(qoe *pbapi.Qoe) *pbapi.Qoe {
	repl := qoe
	repl.Err = pbapi.GatewayErrors_OK

	return repl
}

func (dbc *dbcMock) GetQoE(reqQoes *pbapi.RequestQoE, result chan pbapi.ResponseQoE) {
	// Initializing reply
	replQoEs := pbapi.ResponseQoE{}
	replQoEs.Qoes = make(map[int32]*pbapi.Qoe)
	// Wait group will be used to wait until all workers are done
	// In case worker hangs, grpc client's context will destroy it when
	// timeout expires.
	var wg sync.WaitGroup
	for i, req := range reqQoes.Qoes {
		wg.Add(1)
		go func(i int32, req *pbapi.Qoe) {
			replQoEs.Qoes[i] = dbc.worker(req)
			wg.Done()
		}(i, req)
	}
	wg.Wait()
	glog.Infof("Sending %d QoE back to gRPC server", len(replQoEs.Qoes))
	result <- replQoEs
}

func (dbc *dbcMock) loadTestData() {
	// TODO Consider more flexible way to load test data
	qoe := map[int32]*pbapi.Qoe{
		0: &pbapi.Qoe{
			Src: &pbapi.Endpoint{
				Type:    pbapi.EndpointType_IPV4,
				Address: net.ParseIP("10.0.0.1"),
			},
			Dst: &pbapi.Endpoint{
				Type:    pbapi.EndpointType_IPV4,
				Address: net.ParseIP("20.0.0.1"),
			},
			Qoe: &pbapi.QoeParameters{
				Latency: &pbapi.Latency{
					Value:     200,
					Variation: 10,
				},
			},
			// here is used as "expected" label stack, the server portion will override label stack with actual value
			Label: []uint32{10024, 20024},
		},
		1: &pbapi.Qoe{
			Src: &pbapi.Endpoint{
				Type:    pbapi.EndpointType_IPV6,
				Address: net.ParseIP("2001:1::1"),
			},
			Dst: &pbapi.Endpoint{
				Type:    pbapi.EndpointType_IPV6,
				Address: net.ParseIP("2001:2::1"),
			},
			Qoe: &pbapi.QoeParameters{
				Latency: &pbapi.Latency{
					Value:     200,
					Variation: 10,
				},
			},
			// here is used as "expected" label stack, the server portion will override label stack with actual value
			Label: []uint32{31024, 41024},
		},
		2: &pbapi.Qoe{
			Src: &pbapi.Endpoint{
				Type:    pbapi.EndpointType_IPV4,
				Address: net.ParseIP("10.0.0.1"),
			},
			Dst: &pbapi.Endpoint{
				Type:    pbapi.EndpointType_IPV4,
				Address: net.ParseIP("30.0.0.1"),
			},
			Qoe: &pbapi.QoeParameters{
				Latency: &pbapi.Latency{
					Value:     20,
					Variation: 10,
				},
			},
			// here is used as "expected" label stack, the server portion will override label stack with actual value
			Label: []uint32{12024, 22024},
		},
		3: &pbapi.Qoe{
			Src: &pbapi.Endpoint{
				Type:    pbapi.EndpointType_IPV6,
				Address: net.ParseIP("2001:3::1"),
			},
			Dst: &pbapi.Endpoint{
				Type:    pbapi.EndpointType_IPV6,
				Address: net.ParseIP("2001:4::1"),
			},
			Qoe: &pbapi.QoeParameters{
				Latency: &pbapi.Latency{
					Value:     150,
					Variation: 10,
				},
			},
			// here is used as "expected" label stack, the server portion will override label stack with actual value
			Label: []uint32{33024, 43024},
		},
		4: &pbapi.Qoe{
			Src: &pbapi.Endpoint{
				Type:    pbapi.EndpointType_SID,
				Address: []byte{1, 1, 1, 1},
			},
			Dst: &pbapi.Endpoint{
				Type:    pbapi.EndpointType_SID,
				Address: []byte{2, 2, 2, 2},
			},
			Qoe: &pbapi.QoeParameters{
				Latency: &pbapi.Latency{
					Value:     200,
					Variation: 10,
				},
			},
			// here is used as "expected" label stack, the server portion will override label stack with actual value
			Label: []uint32{34024, 44024},
		},
		5: &pbapi.Qoe{
			Src: &pbapi.Endpoint{
				Type:    pbapi.EndpointType_SID,
				Address: []byte{3, 3, 3, 3},
			},
			Dst: &pbapi.Endpoint{
				Type:    pbapi.EndpointType_SID,
				Address: []byte{4, 4, 4, 4},
			},
			Qoe: &pbapi.QoeParameters{
				Latency: &pbapi.Latency{
					Value:     1500,
					Variation: 10,
				},
			},
			// here is used as "expected" label stack, the server portion will override label stack with actual value
			Label: []uint32{35024, 45024},
		},
		// TODO Add more cases
	}
	dbc.mu.Lock()
	defer dbc.mu.Unlock()
	dbc.store = qoe
}

// NewDBClient return  a new instance of a DB client
func NewMockDBClient(randomDelay, randomFailure bool) DBClient {
	dbc := dbcMock{
		store:         make(map[int32]*pbapi.Qoe),
		randomDelay:   randomDelay,
		randomFailure: randomFailure,
	}
	dbc.loadTestData()

	return &dbc
}
