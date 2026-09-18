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

func VerDashboard(c *gin.Context) {
	var resultados []models.Hiperparametro
	config.DB.Order("id asc").Find(&resultados)
	c.JSON(http.StatusOK, resultados)
}

func IniciarOptimizacion(c *gin.Context) {
	canal := make(chan models.Hiperparametro, 9)

	for i := 1; i <= 9; i++ {
		tarea := models.Hiperparametro{
			Estado: "Procesando",
		}
		config.DB.Create(&tarea)

		go enviarAWorker(int(tarea.ID), canal)
	}

	go func() {
		for i := 1; i <= 9; i++ {
			resultado := <-canal
			config.DB.Model(&models.Hiperparametro{}).Where("id = ?", resultado.ID).Updates(resultado)
		}
	}()

	c.JSON(200, gin.H{"mensaje": "tareas en proceso, revisa la tabla en unos segundos"})
}

func enviarAWorker(id int, canal chan models.Hiperparametro) {
	datos := map[string]interface{}{
		"tasa_aprendizaje": 0.01,
		"profundidad":      5,
	}
	jsonData, _ := json.Marshal(datos)

	resp, err := http.Post("http://load_balancer:9000/entrenar", "application/json", bytes.NewBuffer(jsonData))
	
	precision := 0.0
	estado := "Completado"
	nombreWorker := "Desconocido"

	if err != nil {
		fmt.Println("Falló la conexión con el worker:", err)
		estado = "Error"
	} else {
		defer resp.Body.Close()
		var resultado map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&resultado)

		val, ok := resultado["precision_final"].(float64)
		if !ok {
			fmt.Println("El worker no mando lo esperado. Llegó:", resultado)
			precision = 0.0 
		} else {
			precision = val
		}

		if nombre, ok := resultado["worker_id"].(string); ok {
			nombreWorker = nombre
		}
	}

	tareaActualizada := models.Hiperparametro{
		ID:             uint(id), 
		Estado:         estado,
		PrecisionFinal: precision,
		WorkerAsignado: nombreWorker,
	}

	canal <- tareaActualizada
}