// SPDX-License-Identifier: Apache-2.0
// Copyright 2022-present Open Networking Foundation

package commands

import (
	pb "github.com/xpoh/pfcpsim/api"
	"github.com/xpoh/pfcpsim/internal/pfcpctl/config"
	"github.com/xpoh/pfcpsim/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var conn *grpc.ClientConn

func connect() pb.PFCPSimClient {
	// Create an insecure gRPC Channel
	connection, err := grpc.NewClient(config.GlobalConfig.Server, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.PfcpsimLog.Fatalf("error dialing %v: %v", config.GlobalConfig.Server, err)
	}

	return pb.NewPFCPSimClient(connection)
}

func disconnect() {
	if conn != nil {
		err := conn.Close()
		if err != nil {
			logger.PfcpsimLog.Warnln(err)
		}
	}
}
