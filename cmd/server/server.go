package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ahd99/urlshortner/pkg/logger"
	"github.com/ahd99/urlshortner/pkg/logger/zapLogger"
	"github.com/ahd99/urlshortner/pkg/urlmap"
)

type ServerHandler struct {
	urlMap urlmap.URLMap
}

func (server ServerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[1:]
	redirectUrl, err := server.urlMap.GetUrl(key)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		logger1.Error("Error in finding url", logger.String("err", err.Error()), logger.String("key", key), logger.String("url", redirectUrl))
		return
	}
	logger1.Debug("url found.", logger.String("key", key), logger.String("url", redirectUrl))
	w.Header().Set("location", redirectUrl)
	w.WriteHeader(http.StatusMovedPermanently)
}

func startServer(urlmap1 urlmap.URLMap, port int) {
	server := ServerHandler{
		urlMap: urlmap1,
	}
	log.Println("Server started on port", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), server))
}

var logger1 logger.Logger

func main() {
	logger1 = initLogger()
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
