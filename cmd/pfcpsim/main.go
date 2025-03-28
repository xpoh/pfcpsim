// SPDX-License-Identifier: Apache-2.0
// Copyright 2022-present Open Networking Foundation

package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/pborman/getopt/v2"
	pb "github.com/xpoh/pfcpsim/api"
	"github.com/xpoh/pfcpsim/internal/pfcpsim"
	"github.com/xpoh/pfcpsim/logger"
	"google.golang.org/grpc"
)

const (
	defaultgRPCServerPort = "54321"
)

func startServer(apiDoneChannel chan bool, iFace string, port string, group *sync.WaitGroup) {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", port))
	if err != nil {
		logger.PfcpsimLog.Fatalf("api gRPC Server failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterPFCPSimServer(grpcServer, pfcpsim.NewPFCPSimService(iFace))

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			logger.PfcpsimLog.Fatalf("failed to listed: %v", err)
		}
	}()

	logger.PfcpsimLog.Infoln("server listening on port", port)

	x := <-apiDoneChannel
	if x {
		// if the API channel is closed, stop the gRPC pfcpsim
		grpcServer.Stop()
		logger.PfcpsimLog.Warnln("stopping API gRPC pfcpsim")
	}

	group.Done()
}

func main() {
	port := getopt.StringLong("port", 'p', defaultgRPCServerPort, "the gRPC Server port to listen")
	iFaceName := getopt.StringLong("interface", 'i', "", "Defines the local address. If left blank,"+
		" the IP will be taken from the first non-loopback interface")

	optHelp := getopt.BoolLong("help", 0, "Help")

	getopt.Parse()

	if *optHelp {
		getopt.Usage()
		os.Exit(0)
	}

	// control channels, they are only closed when the goroutine needs to be terminated
	doneChannel := make(chan bool)

	sigs := make(chan os.Signal, 1)
	// stop API servers on SIGTERM
	signal.Notify(sigs, syscall.SIGTERM)

	go func() {
		<-sigs
		close(doneChannel)
	}()

	wg := sync.WaitGroup{}
	wg.Add(1)

	go startServer(doneChannel, *iFaceName, *port, &wg)
	logger.PfcpsimLog.Debugln("started API gRPC Service")

	wg.Wait()

	defer func() {
		logger.PfcpsimLog.Infoln("pfcp Simulator shutting down")
	}()
}
