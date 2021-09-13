package main

import (
	"fmt"
	"net"

	"github.com/ahd99/urlshortner/internal/monitoring/proto"
	"github.com/ahd99/urlshortner/internal/monitoring/server"
	"github.com/ahd99/urlshortner/pkg/logger"
	"google.golang.org/grpc"
)

func StartMonitoringServer(port int) {
	list, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		logger1.Fatal("Error listening on monitoring GRPC port", logger.Int("port", port), logger.String("Err", err.Error()))
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	proto.RegisterMonitoringServer(grpcServer, server.GetServer())
	logger1.Info("Starting monitoring grpc server. ", logger.Int("port", port))
	err = grpcServer.Serve(list)
	logger1.Error("Monitoring server stopped with error.", logger.String("err", err.Error()))
}
