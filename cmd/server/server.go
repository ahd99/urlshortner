package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ahd99/urlshortner/pkg/urlmap"
)


type ServerHandler struct {
	urlMap urlmap.URLMap
}

func (server ServerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[1:]
	redirectUrl, err := server.urlMap.GetUrl(key)
	log.Println(key, redirectUrl, err)
	if  err != nil{
		w.WriteHeader(http.StatusNotFound)
		log.Printf("%v\n", err)
		return
	}
	w.Header().Set("location", redirectUrl)
	w.WriteHeader(http.StatusMovedPermanently)
}

func startServer(urlmap1 urlmap.URLMap, port int) {
	server := ServerHandler{
		urlMap:	urlmap1,
	}
	log.Println("Server started on port", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), server))
}

func main() {
	urlmap := urlmap.New()
	urlmap.Add("dig", "https://digiato.com")
	urlmap.Add("asr", "https://asriran.com")
	startServer(urlmap, 8081)
}

