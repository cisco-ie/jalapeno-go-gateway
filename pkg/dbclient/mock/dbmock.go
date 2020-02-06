package mock

import (
	"bytes"
	"context"
	"net"
	"sync"
	"time"

	pbapi "github.com/cisco-ie/jalapeno-go-gateway/pkg/apis"
	"github.com/cisco-ie/jalapeno-go-gateway/pkg/bgpclient"
	"github.com/cisco-ie/jalapeno-go-gateway/pkg/dbclient"
	"github.com/golang/glog"
)

var (
	// maxDBRequestTimeout defines a maximum wait time for db worker to retrieve QoE information
	maxDBRequestTimeout = time.Millisecond * 2000
)

type labelInfo struct {
	label uint32
	rt    []bgpclient.RTValue
}
type dbMock struct {
	mu  sync.Mutex
	qoe map[int32]*pbapi.Qoe
	vpn map[bgpclient.RDValue]labelInfo
}

func (db *dbMock) GetVPN(ctx context.Context, r *bgpclient.VPNRequest, ch chan *bgpclient.VPNReply) {
	glog.V(5).Infof("Mocked Data Store received request for RD: %+v RT: %+v", *r.RD, r.RT)
	var repl *bgpclient.VPNReply
	db.mu.Lock()
	defer db.mu.Unlock()

	li, ok := db.vpn[*r.RD]
	if !ok {
		ch <- &bgpclient.VPNReply{}
		glog.V(5).Infof("Mocked Data Store did not find requested RD: %+v", *r.RD)
		return
	}
	repl = &bgpclient.VPNReply{
		RD:    r.RD,
		RT:    li.rt,
		Label: li.label,
	}
	ch <- repl
	glog.V(5).Infof("Mocked Data store returns RD: %+v RT: %+v Label: %d", *repl.RD, repl.RT, repl.Label)
}

// GetQoE is required method by DB interface, it takes requested QoE and searches through
// stored in memory map with mock entries. If found, it returns the value, otherwise return ENOENT error.
func (db *dbMock) GetQoE(ctx context.Context, req *pbapi.Qoe, ch chan *pbapi.Qoe) {
	repl := &pbapi.Qoe{}
	db.mu.Lock()
	defer db.mu.Unlock()

	// Simulating hung DB's Get routine
	// time.Sleep(time.Second *3)

	for _, qoe := range db.qoe {
		if bytes.Compare(req.Src.Address, qoe.Src.Address) == 0 && bytes.Compare(req.Dst.Address, qoe.Dst.Address) == 0 {
			// TODO add more precise calculation of latency with considerring variation.
			if req.Qoe.Latency.Value == qoe.Qoe.Latency.Value {
				repl := qoe
				repl.Err = pbapi.GatewayErrors_OK
				ch <- repl
				return
			}
		}
	}
	repl.Err = pbapi.GatewayErrors_ENOENT
	ch <- repl
}

// loadTestData loads test data into the mock DB.
func (db *dbMock) loadTestData() {
	// TODO Consider more flexible way to load test data
	vpn := map[bgpclient.RDValue]labelInfo{
		{
			T:     bgpclient.RouteDistinguisherTwoOctetAS,
			Value: [8]byte{2, 65, 0, 0, 0, 0, 253, 234},
		}: {
			label: 24000,
			rt: []bgpclient.RTValue{
				{
					T:     bgpclient.TwoOctetAsSpecificExtended,
					Value: [8]byte{2, 65, 0, 0, 0, 0, 253, 234},
				},
			},
		},
		{
			T:     bgpclient.RouteDistinguisherTwoOctetAS,
			Value: [8]byte{2, 65, 0, 0, 0, 0, 253, 235},
		}: {
			label: 24001,
			rt: []bgpclient.RTValue{
				{
					T:     bgpclient.TwoOctetAsSpecificExtended,
					Value: [8]byte{2, 65, 0, 0, 0, 0, 253, 235},
				},
			},
		},
		{
			T:     bgpclient.RouteDistinguisherTwoOctetAS,
			Value: [8]byte{2, 65, 0, 0, 0, 0, 253, 236},
		}: {
			label: 24002,
			rt: []bgpclient.RTValue{
				{
					T:     bgpclient.TwoOctetAsSpecificExtended,
					Value: [8]byte{2, 65, 0, 0, 0, 0, 253, 236},
				},
			},
		},
	}
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
			Label: []uint32{35024, 45024},
		},
	}
	db.mu.Lock()
	defer db.mu.Unlock()
	db.qoe = qoe
	db.vpn = vpn
}

// NewMockDB return  a new instance of a DB client
func NewMockDB() dbclient.DB {
	db := dbMock{
		qoe: make(map[int32]*pbapi.Qoe),
	}
	db.loadTestData()

	return &db
}
