package bgpclient

import (
	"context"

	api "github.com/osrg/gobgp/api"
	"google.golang.org/grpc"
)

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
