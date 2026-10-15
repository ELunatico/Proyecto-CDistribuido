package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type ConsulService struct {
	Service struct {
		Address string `json:"Address"`
		Port    int    `json:"Port"`
	} `json:"Service"`
}

func getHealthyWorker() string {
	resp, err := http.Get("http://consul:8500/v1/health/service/python-worker?passing=true")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	var nodos []ConsulService
	json.NewDecoder(resp.Body).Decode(&nodos)

	if len(nodos) == 0 {
		return ""
	}

	nodo := nodos[rand.Intn(len(nodos))]

	return fmt.Sprintf("http://%s:%d", nodo.Service.Address, nodo.Service.Port)
}

func handlerBalanceador(w http.ResponseWriter, r *http.Request) {
	workerURL := getHealthyWorker()
	
	if workerURL == "" {
		http.Error(w, "No hay workers disponibles o Consul está caído", http.StatusServiceUnavailable)
		return
	}

	urlDestino, _ := url.Parse(workerURL)
	proxy := httputil.NewSingleHostReverseProxy(urlDestino)
	proxy.ServeHTTP(w, r)
}

func main() {
	http.HandleFunc("/entrenar", handlerBalanceador)
	fmt.Println("Load Balancer con Service Discovery Consul escuchando en 9000")
	http.ListenAndServe(":9000", nil)
}