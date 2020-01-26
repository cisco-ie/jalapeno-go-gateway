package gateway

import (
	"context"
	"flag"
	"net"
	"os"
	"testing"
	"time"

	pbapi "github.com/cisco-ie/jalapeno-go-gateway/pkg/apis"
	"github.com/cisco-ie/jalapeno-go-gateway/pkg/dbclient"
	"github.com/cisco-ie/jalapeno-go-gateway/pkg/dbclient/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 4096 * 1024

var conn *bufconn.Listener

func TestMain(m *testing.M) {
	flag.Parse()
	flag.Set("logtostderr", "true")
	conn = bufconn.Listen(bufSize)
	db := mock.NewMockDB()
	dbc := dbclient.NewDBClient(db)
	gSrv := NewGateway(conn, dbc)
	gSrv.Start()
	os.Exit(m.Run())
}

func bufDialer(string, time.Duration) (net.Conn, error) {
	return conn.Dial()
}

func TestClientRequestQoE(t *testing.T) {
	tests := []struct {
		name       string
		qoeRequest map[int32]*pbapi.Qoe
		peer       *peer.Peer
	}{
		{
			name: "Single QoE successful request",
			qoeRequest: map[int32]*pbapi.Qoe{
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
					// Err here is used as "expected" error, the server portion will override err with actual value
					Err: pbapi.GatewayErrors_OK,
				},
			},
			peer: &peer.Peer{
				Addr: &net.IPAddr{
					IP: net.ParseIP("1.1.1.1"),
				},
			},
		},
		{
			name: "2 QoEs successful request",
			qoeRequest: map[int32]*pbapi.Qoe{
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
					// Err here is used as "expected" error, the server portion will override err with actual value
					Err: pbapi.GatewayErrors_OK,
				},
				1: &pbapi.Qoe{
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
					// Err here is used as "expected" error, the server portion will override err with actual value
					Err: pbapi.GatewayErrors_OK,
				},
			},
			peer: &peer.Peer{
				Addr: &net.IPAddr{
					IP: net.ParseIP("1.1.1.2"),
				},
			},
		},
	}
	for _, tt := range tests {
		ctx := peer.NewContext(context.TODO(), tt.peer)
		conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithDialer(bufDialer), grpc.WithInsecure())
		if err != nil {
			t.Fatalf("failed to dial bufnet with error: %v", err)
		}
		defer conn.Close()
		client := pbapi.NewQoEServiceClient(conn)
		resp, err := client.QoE(ctx, &pbapi.RequestQoE{Qoes: tt.qoeRequest})
		if err != nil {
			t.Fatalf("QoE request failed with error: %v", err)
		}
		for key, qoe := range resp.Qoes {
			// First check if expected error state matches to the one returned by the server
			if qoe.Err != tt.qoeRequest[key].Err {
				t.Fatalf("test \"%s\" failed, expected error: %+v but got %+v", tt.name, tt.qoeRequest[key].Err, qoe.Err)
			}
			// If error was expected, then nothing else left to do, on to the next item.
			if qoe.Err != pbapi.GatewayErrors_OK {
				continue
			}
			// Now check if returned by the server stack matches to the expected stack.
			for i, label := range tt.qoeRequest[key].Label {
				if !labelExists(qoe.Label, label) {
					t.Fatalf("test \"%s\" failed, expected label %d is not found in server's reply.", tt.name, tt.qoeRequest[key].Label[i])
				}
			}
		}
	}
}

func labelExists(labels []uint32, label uint32) bool {
	for _, l := range labels {
		if l == label {
			return true
		}
	}
	return false
}
