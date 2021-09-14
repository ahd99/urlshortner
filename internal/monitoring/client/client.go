package client

import (
	"context"
	"io"

	"github.com/ahd99/urlshortner/internal/monitoring/proto"
	"github.com/ahd99/urlshortner/internal/monitoring/server"
	"github.com/ahd99/urlshortner/pkg/logger"
	"google.golang.org/grpc"
)

var logger1 logger.Logger
var conn *grpc.ClientConn

// target format: 100.10.10.1:8080
func Init(taget string, log logger.Logger) {
	logger1 = log
	//var opts []grpc.DialOption
	conn, err := grpc.Dial(taget, grpc.WithInsecure())
	if err != nil {
		logger1.Error("Error connecting to monitoring server via grpc.", logger.String("server", taget), logger.String("err", err.Error()))
		return
	}
	client := proto.NewMonitoringClient(conn)

	tranChannel := make(chan *server.Request, 1000)

	go servServerTranList(client, tranChannel)
	go servServerStat(client)
	showTransactions(tranChannel)
}

func Cleanup() {
	if conn != nil {
		conn.Close()
	}
}

func servServerTranList(client proto.MonitoringClient, tranchannel chan<- *server.Request) {
	ctx := context.Background()
	tranClient, err := client.TranList(ctx, &proto.TranListReq{})
	if err != nil {
		logger1.Error("Client -- ReqList -- Error getting stream", logger.String("err", err.Error()))
		close(tranchannel)
		return
	}
	logger1.Debug("Client -- ReqList -- stream created.")

	for {
		req, err := tranClient.Recv()
		if err == io.EOF {
			logger1.Debug("Client -- ReqList -- Channel closed (EOF)")
			break
		} else if err != nil {
			logger1.Error("Client -- ReqList -- Error Received from server.", logger.String("err", err.Error()))
			break
		}
		logger1.Debug("Client -- Req received.", logger.String("key", req.Key), logger.String("url", req.Url), logger.String("ip", req.Ip))
		tranchannel <- &server.Request{Key: req.Key, Url: req.Url, Ip: req.Ip}
		logger1.Debug("Client -- Req sent top channel.", logger.String("key", req.Key), logger.String("url", req.Url), logger.String("ip", req.Ip))
	}
	close(tranchannel)
	logger1.Error("Client -- ReqList -- Finish")
}

func servServerStat(client proto.MonitoringClient) {

}

func showTransactions(tranchannel <-chan *server.Request) {
	for req := range tranchannel {
		logger1.Debug("Client -- Req read from channel.", logger.String("key", req.Key), logger.String("url", req.Url), logger.String("ip", req.Ip))
	}
}
