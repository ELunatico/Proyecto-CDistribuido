package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var listaWorkers []string
var turno int = 0

type Discovery struct {
	Workers []string `json:"workers"`
}

func main() {
	archivo, err := ioutil.ReadFile("discovery.json")
	if err != nil {
		fmt.Println("No se pudo leer el discovery.json:", err)
		return
	}

	var datos Discovery
	json.Unmarshal(archivo, &datos)
	listaWorkers = datos.Workers

	fmt.Println("Workers encontrados:", len(listaWorkers))

	http.HandleFunc("/", repartirTranajo)

	fmt.Println("Load Balancer arrancando en el puerto 9000")
	http.ListenAndServe(":9000", nil)
}

func repartirTranajo(w http.ResponseWriter, r *http.Request) {
	if len(listaWorkers) == 0 {
		fmt.Println("No hay workers disponibles")
		return
	}

	workerAsignado := listaWorkers[turno]
	fmt.Println("Redirigiendo petición hacia ->", workerAsignado)

	urlDestino, _ := url.Parse(workerAsignado)

	proxy := httputil.NewSingleHostReverseProxy(urlDestino)
	proxy.ServeHTTP(w, r)

	turno = turno + 1
	
	if turno >= len(listaWorkers) {
		turno = 0
	}
}