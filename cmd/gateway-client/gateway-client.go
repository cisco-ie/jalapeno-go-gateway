package main

import (
	"context"
	"io"
	"math/rand"
	"net"
	"os"
	"os/signal"

	pbapi "github.com/cisco-ie/jalapeno-go-gateway/pkg/apis"
	"github.com/golang/glog"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

const (
	jalapenoGateway = "jalapeno-gateway:40040"
)

func main() {
	client := os.Getenv("CLIENT_IP")
	if client == "" {
		glog.Errorf("env variable \"CLIENT_IP\" is not defined, cannor proceed further, exiting...")
		os.Exit(1)
	}
	if net.ParseIP(client).To4() == nil && net.ParseIP(client).To16() == nil {
		glog.Errorf("env variable \"CLIENT_IP\" carries an invalid ip %s", client)
		os.Exit(1)
	}
	ctx := peer.NewContext(context.TODO(), &peer.Peer{
		Addr: &net.IPAddr{
			IP: net.ParseIP(client),
		}})
	conn, err := grpc.DialContext(ctx, jalapenoGateway)
	if err != nil {
		glog.Errorf("failed to connect to Jalapeno Gateway at the address: %s with error: %+v", jalapenoGateway, err)
		os.Exit(1)
	}
	defer conn.Close()
	gwclient := pbapi.NewGatewayServiceClient(conn)

	rd1, _ := ptypes.MarshalAny(&pbapi.RouteDistinguisherTwoOctetAS{Admin: uint32(577), Assigned: uint32(65002)})
	rd2, _ := ptypes.MarshalAny(&pbapi.RouteDistinguisherTwoOctetAS{Admin: uint32(577), Assigned: uint32(65003)})
	rd3, _ := ptypes.MarshalAny(&pbapi.RouteDistinguisherTwoOctetAS{Admin: uint32(577), Assigned: uint32(65004)})
	rt1, _ := ptypes.MarshalAny(&pbapi.TwoOctetAsSpecificExtended{As: uint32(577), LocalAdmin: uint32(65002)})
	rt2, _ := ptypes.MarshalAny(&pbapi.TwoOctetAsSpecificExtended{As: uint32(577), LocalAdmin: uint32(65003)})
	rt3, _ := ptypes.MarshalAny(&pbapi.TwoOctetAsSpecificExtended{As: uint32(577), LocalAdmin: uint32(65004)})
	requests := []*pbapi.RequestVPN{
		{
			Rd:       rd1,
			ImportRt: []*any.Any{rt1},
		},
		{
			Rd:       rd2,
			ImportRt: []*any.Any{rt2},
		},
		{
			Rd:       rd3,
			ImportRt: []*any.Any{rt3},
		},
	}
	stopCh := setupSignalHandler()
	// The client will randomly send requests between three VRFs green, blue and red.
	for {
		i := rand.Intn(3)
		stream, err := gwclient.VPN(ctx, requests[i])
		if err != nil {
			glog.Errorf("failed to request VPN label for request %+v with error: %+v", requests[i])
		}
		for {
			entry, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				glog.Errorf("failed to receive a message from the stream with error: %+v", err)
			}
			glog.Infof("Received message: %+v", *entry)
		}
		select {
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
