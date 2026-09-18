package main

import (
	"fmt"
	
	"monolito-go/config"
	"monolito-go/models"
	"monolito-go/routes"
)

func main() {
	fmt.Println("Iniciando")

	config.ConectarBaseDeDatos()

	err := config.DB.AutoMigrate(&models.Hiperparametro{})
	if err != nil {
		fmt.Println("Error al migrar la base de datos: %v", err)
	}
	fmt.Println("Tabla lista.")

	router := routes.ConfigurarRutas()

	err = router.Run(":8080")
	if err != nil {
		fmt.Println("El servidor se murio: %v", err)
	}
}