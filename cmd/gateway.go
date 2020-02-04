package main

import (
	"flag"
	"fmt"
	"math"
	"net"
	"os"
	"os/signal"
	"strconv"

	"github.com/cisco-ie/jalapeno-go-gateway/pkg/gateway"
	"github.com/golang/glog"
)

const (
	// DefaultGatewayPort defines default port Gateway's gRPC server listens on
	// this port is a container port, not the port used for Gateway Kubernetes Service.
	defaultGatewayPort = "15151"
)

func main() {
	flag.Parse()
	flag.Set("logtostderr", "true")

	fmt.Printf("Gateway...\n")

	// Getting port for gRPC server to listen on, from environment varialbe
	// GATEWAY_PORT
	strPort := os.Getenv("GATEWAY_PORT")
	if strPort == "" {
		// TODO, should it fallback to the default port?
		strPort = defaultGatewayPort
		glog.Warningf("env variable \"GATEWAY_PORT\" is not defined, using default Gateway port: %s", strPort)
	}
	srvPort, err := strconv.Atoi(strPort)
	if err != nil {
		glog.Warningf("env variable \"GATEWAY_PORT\" containes an invalid value %s, using default Gateway port instead: %s", strPort, defaultGatewayPort)
		srvPort, _ = strconv.Atoi(defaultGatewayPort)
	}
	if srvPort == 0 || srvPort > math.MaxUint16 {
		glog.Warningf("env variable \"GATEWAY_PORT\" containes an invalid value %d, using default Gateway port instead: %s\n", srvPort, defaultGatewayPort)
		srvPort, _ = strconv.Atoi(defaultGatewayPort)
	}
	// Initialize DB client

	// Initialize gRPC server
	conn, err := net.Listen("tcp", ":"+strPort)
	if err != nil {
		glog.Errorf("failed to setup listener with with error: %+v", err)
		os.Exit(1)
	}

	gSrv := gateway.NewGateway(conn)
	gSrv.Start()

	// For now just get stuck on stop channel, later add signal processing
	stopCh := setupSignalHandler()
	<-stopCh
	// Clean up section
	// Cleanup DB connection
	// Cleanup gRPC server
	gSrv.Stop()

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
