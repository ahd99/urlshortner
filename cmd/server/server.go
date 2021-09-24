package main

import (
	"fmt"
	"net/http"
	//"time"

	monitor "github.com/ahd99/urlshortner/internal/monitoring/server"
	"github.com/ahd99/urlshortner/pkg/logger"
	"github.com/ahd99/urlshortner/pkg/logger/zapLogger"
	//"github.com/ahd99/urlshortner/pkg/mongodb"
	"github.com/ahd99/urlshortner/pkg/urlmap"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type ServerHandler struct {
	urlMap urlmap.URLMap
}

func (server ServerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[1:]
	redirectUrl, err := server.urlMap.GetUrl(key)

	defer func() {
		totalReq.WithLabelValues(key, "1").Inc()
		monitor.RequestReceived(key, redirectUrl, r.RemoteAddr)
	}()

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		logger1.Error("Error in finding url", logger.String("err", err.Error()), logger.String("key", key), logger.String("url", redirectUrl))
		return
	}
	logger1.Debug("url found.", logger.String("key", key), logger.String("url", redirectUrl))
	//mongodb.InsertReqLog(key, redirectUrl, r.RemoteAddr, time.Now())

	w.Header().Set("location", redirectUrl)
	w.WriteHeader(http.StatusMovedPermanently)
}

func startServer(urlmap1 urlmap.URLMap, port int) {
	server := ServerHandler{
		urlMap: urlmap1,
	}
	logger1.Info("Starting Server...", logger.Int("port", 8081))
	err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), server)
	logger1.Fatal("ListenAndServe Error.", logger.Int("port", 8081), logger.String("Error", err.Error()))
}

var logger1 logger.Logger

func main() {
	logger1 = initLogger()
	go initPrometheus()
	
	//mongodb.InitMongo(logger1)
	//defer mongodb.CloseMongo()

	go StartMonitoringServer(8091)

	urlmap := urlmap.New()
	urlmap.Add("dig", "https://digiato.com")
	urlmap.Add("asr", "https://asriran.com")
	startServer(urlmap, 8081)
}

func initLogger() logger.Logger {
	logger.SetKeyValuePairFactory(&zapLogger.ZapKeyValFactory{})
	logger1 := logger.NewLogger(zapLogger.ZapLoggerFactory{})
	return logger1
}

var totalReq *prometheus.CounterVec

func initPrometheus() {
	//promauto.NewCounter()
	totalReq = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "Req_Count_Total",
		Help: "counter for received requests",
	}, []string{"key", "res"})
	prometheus.Register(totalReq)
	err := http.ListenAndServe(":2112", promhttp.Handler())
	logger1.Error("Prometheus ListenAndServe Error.", logger.Int("port", 2112), logger.String("Error", err.Error()))
}
