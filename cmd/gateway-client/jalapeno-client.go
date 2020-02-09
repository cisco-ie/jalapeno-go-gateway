package main

import (
	"context"
	"flag"
	"io"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"time"

	pbapi "github.com/cisco-ie/jalapeno-go-gateway/pkg/apis"
	"github.com/cisco-ie/jalapeno-go-gateway/pkg/bgpclient"
	"github.com/golang/glog"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/osrg/gobgp/pkg/packet/bgp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	jalapenoGateway = "jalapeno-gateway:40040"
)

func main() {
	flag.Parse()
	flag.Set("logtostderr", "true")

	client := os.Getenv("CLIENT_IP")
	if client == "" {
		glog.Errorf("env variable \"CLIENT_IP\" is not defined, cannor proceed further, exiting...")
		os.Exit(1)
	}
	if net.ParseIP(client).To4() == nil && net.ParseIP(client).To16() == nil {
		glog.Errorf("env variable \"CLIENT_IP\" carries an invalid ip %s", client)
		os.Exit(1)
	}
	glog.Errorf("\"CLIENT_IP:\" %s", client)

	conn, err := grpc.DialContext(context.TODO(), jalapenoGateway, grpc.WithInsecure())
	if err != nil {
		glog.Errorf("failed to connect to Jalapeno Gateway at the address: %s with error: %+v", jalapenoGateway, err)
		os.Exit(1)
	}
	defer conn.Close()
	gwclient := pbapi.NewGatewayServiceClient(conn)
	// Mocking RDs and RTs for 3 VRFs
	requests := make([]*pbapi.RequestVPN, 3)
	requests[0] = &pbapi.RequestVPN{
		Rt: make([]*any.Any, 0),
	}
	requests[0].Rd = bgpclient.MarshalRD(bgp.NewRouteDistinguisherTwoOctetAS(577, 65000))
	requests[0].Rt = append(requests[0].Rt, bgpclient.MarshalRT(bgp.NewTwoOctetAsSpecificExtended(bgp.EC_SUBTYPE_ROUTE_TARGET, 577, 65000, false)))

	requests[1] = &pbapi.RequestVPN{
		Rt: make([]*any.Any, 0),
	}
	requests[1].Rd = bgpclient.MarshalRD(bgp.NewRouteDistinguisherIPAddressAS("57.57.57.57", 65001))
	requests[1].Rt = append(requests[0].Rt, bgpclient.MarshalRT(bgp.NewIPv4AddressSpecificExtended(bgp.EC_SUBTYPE_ROUTE_TARGET, "57.57.57.57", 65001, false)))
	requests[2] = &pbapi.RequestVPN{
		Rt: make([]*any.Any, 0),
	}
	requests[2].Rd = bgpclient.MarshalRD(bgp.NewRouteDistinguisherFourOctetAS(456734567, 65002))
	requests[2].Rt = append(requests[0].Rt, bgpclient.MarshalRT(bgp.NewFourOctetAsSpecificExtended(bgp.EC_SUBTYPE_ROUTE_TARGET, 456734567, 65002, false)))

	stopCh := setupSignalHandler()
	// The client will randomly send requests between three VRFs green, blue and red.
	ticker := time.NewTicker(time.Second * 10)
	ctx := metadata.NewOutgoingContext(context.TODO(), metadata.New(map[string]string{
		"CLIENT_IP": net.ParseIP(client).String(),
	}))
	for {
		i := rand.Intn(3)
		stream, err := gwclient.VPN(ctx, requests[i])
		if err != nil {
			glog.Errorf("failed to request VPN label for request %+v with error: %+v", requests[i], err)
		} else {
			for {
				entry, err := stream.Recv()
				if err == io.EOF {
					break
				}
				if err != nil {
					glog.Errorf("failed to receive a message from the stream with error: %+v", err)
					break
				}
				if entry == nil {
					glog.Info("Received empty message\n")
					continue
				}
				rd, err := bgpclient.UnmarshalRD(entry.Rd)
				if err != nil {
					glog.Errorf("failed to unmarshal received rd with error: %+v", err)
				}
				glog.Infof("Recived RD: %+v", rd)
				rts, err := bgpclient.UnmarshalRT(entry.Rt)
				if err != nil {
					glog.Errorf("failed to unmarshal received rt with error: %+v", err)
				}
				glog.Infof("Recived RTs: %+v", rts)
			}
		}
		select {
		case <-ticker.C:
		case <-stopCh:
			glog.Info("Stop signal received, exiting...")
			os.Exit(0)
		}
	}
}

var (
	onlyOneSignalHandler = make(chan struct{})
	shutdownSignals      = []os.Signal{os.Interrupt}
)

func setupSignalHandler() (stopCh <-chan struct{}) {
	close(onlyOneSignalHandler) // panics when called twice

	stop := make(chan struct{})
	c := make(chan os.Signal, 2)
	signal.Notify(c, shutdownSignals...)
	go func() {
		<-c
		close(stop)
		<-c
		os.Exit(1) // second signal. Exit directly.
	}()

	return stop
}
