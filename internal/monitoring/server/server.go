package server

import (
	"time"

	"github.com/ahd99/urlshortner/internal/monitoring/proto"
)

type MonitoringServer struct {
	requestCount int64
	tranChannel  chan request
	statChannel  chan bool

	proto.UnimplementedMonitoringServer
}

var server *MonitoringServer = &MonitoringServer{
	requestCount: 0,
	statChannel:  make(chan bool, 10),
	tranChannel:  make(chan request, 1000),
}

type request struct {
	key string
	url string
	ip  string
}

func (server *MonitoringServer) Statistics(req *proto.StatReq, resp proto.Monitoring_StatisticsServer) error {
	for {
		<-server.statChannel
		err := resp.Send(&proto.StatResp{Count: server.requestCount})
		if err != nil {
			return err
		}
		time.Sleep(1 * time.Second)
	}
}

func (server *MonitoringServer) TranList(req *proto.TranListReq, resp proto.Monitoring_TranListServer) error {
	for req := range server.tranChannel {
		err := resp.Send(&proto.TranListResp{Key: req.key, Url: req.url, Ip: req.ip})
		if err != nil {
			return err
		}
	}
	return nil
}

func GetServer() *MonitoringServer {
	return server
}

func RequestReceived(key string, url string, ip string) {
	server.requestCount++ //TODO: needs lock.
	select {
	case server.statChannel <- true:
	default:
	}

	select {
	case server.tranChannel <- request{key, url, ip}:
	default:
	}

}
