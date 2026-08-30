package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	
	"monolito-go/config"
	"monolito-go/models"
)

type PayloadWorker struct {
	TasaAprendizaje  float64 `json:"tasa_aprendizaje"`
	ProfundidadArbol int     `json:"profundidad_arbol"`
	Estimadores      int     `json:"estimadores"`
}

func enviarAWorker(idWorker int, tareas <-chan models.Hiperparametro) {
	for tarea := range tareas {
		config.DB.Model(&tarea).Updates(map[string]interface{}{
			"Estado": "procesando",
			"WorkerAsignado": fmt.Sprintf("Worker-%d", idWorker),
		})

		payload := PayloadWorker{
			TasaAprendizaje:  tarea.TasaAprendizaje,
			ProfundidadArbol: tarea.ProfundidadArbol,
			Estimadores:      tarea.Estimadores,
		}
		jsonData, _ := json.Marshal(payload)

		resp, err := http.Post("http://worker_python:8000/entrenar", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Println("Error conectando al worker a Python:", err)
			config.DB.Model(&tarea).Update("Estado", "error")
			continue
		}
		
		var resultado map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&resultado)
		resp.Body.Close()

		precision := resultado["precision_final"].(float64)
		config.DB.Model(&tarea).Updates(map[string]interface{}{
			"Estado": "completado",
			"PrecisionFinal": precision,
		})
	}
}

func IniciarOptimizacion(c *gin.Context) {
	combinaciones := []models.Hiperparametro{
		{TasaAprendizaje: 0.01, ProfundidadArbol: 5, Estimadores: 100},
		{TasaAprendizaje: 0.05, ProfundidadArbol: 10, Estimadores: 200},
		{TasaAprendizaje: 0.1, ProfundidadArbol: 15, Estimadores: 300},
	}

	for i := range combinaciones {
		config.DB.Create(&combinaciones[i])
	}

	canalTareas := make(chan models.Hiperparametro, len(combinaciones))

	for w := 1; w <= 2; w++ {
		go enviarAWorker(w, canalTareas)
	}

	for _, tarea := range combinaciones {
		canalTareas <- tarea
	}
	close(canalTareas)

	c.Redirect(http.StatusFound, "/")
}

func VerDashboard(c *gin.Context) {
	var resultados []models.Hiperparametro
	
	config.DB.Order("id asc").Find(&resultados)

	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"resultados": resultados,
	})
}