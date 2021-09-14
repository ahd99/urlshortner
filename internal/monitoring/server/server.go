package server

import (
	"strconv"
	"time"

	"github.com/ahd99/urlshortner/internal/monitoring/proto"
	"github.com/ahd99/urlshortner/pkg/logger"
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

var logger1 logger.Logger

type request struct {
	key string
	url string
	ip  string
}

func (server *MonitoringServer) Statistics(req *proto.StatReq, resp proto.Monitoring_StatisticsServer) error {
	logger1.Info("Monitoring -- statistics -- Connection received")
	for {
		<-server.statChannel
		logger1.Debug("Monitoring -- statistics -- event received from channel")
		err := resp.Send(&proto.StatResp{Count: server.requestCount})
		if err != nil {
			logger1.Error("Monitoring -- statistics -- Error sending count.",
				logger.Int64("count", server.requestCount),
				logger.String("err", err.Error()))
			return err
		}
		logger1.Debug("Monitoring -- statistics -- count sent to client successfully")
		time.Sleep(1 * time.Second)
	}
}

var tranListConnId_debug int = 0

func (server *MonitoringServer) TranList(req *proto.TranListReq, resp proto.Monitoring_TranListServer) error {
	tranListConnId_debug++
	connId := tranListConnId_debug
	logger1.Info("Monitoring -- Reqlist(" + strconv.Itoa(connId) + ") -- Connection received.")

	for {
		select {
		case req := <-server.tranChannel:
			logger1.Debug("Monitoring -- Reqlist("+strconv.Itoa(connId)+") -- Req received from channel", logger.String("key", req.key), logger.String("url", req.url), logger.String("ip", req.ip))
			err := resp.Send(&proto.TranListResp{Key: req.key, Url: req.url, Ip: req.ip})
			if err != nil {
				logger1.Error("Monitoring -- Reqlist("+strconv.Itoa(connId)+" ) -- Error sending Req.", logger.String("key", req.key), logger.String("url", req.url), logger.String("ip", req.ip), logger.String("err", err.Error()))
				server.tranChannel <- req
				return err
			}
			logger1.Debug("Monitoring -- Reqlist("+strconv.Itoa(connId)+") -- REq sent to client successfully", logger.String("key", req.key), logger.String("url", req.url), logger.String("ip", req.ip))
		case <-resp.Context().Done():
			logger1.Debug("Monitoring -- Reqlist("+strconv.Itoa(connId)+") -- Channel was done (closed).",
				logger.String("err", resp.Context().Err().Error()))
			return resp.Context().Err()
		}
	}
}

func RequestReceived(key string, url string, ip string) {
	//logger1.Debug("Monitoring. Request received start.", logger.String("key", key), logger.String("url", url), logger.String("ip", ip), logger.Int("count", int(server.requestCount)))
	server.requestCount++ //TODO: needs lock.
	select {
	case server.statChannel <- true:
		logger1.Debug("true sent to statchannel")
	default:
	}

	select {
	case server.tranChannel <- request{key, url, ip}:
		logger1.Debug("tran sent to tranchannel")
	default:
	}
	//logger1.Debug("Monitoring. Request received finish.")

}

func GetServer() *MonitoringServer {
	return server
}

func SetLogger(log logger.Logger) {
	logger1 = log
}
