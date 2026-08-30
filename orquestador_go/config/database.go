package config

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConectarBaseDeDatos() {
	dsn := "host=db_postgres user=postgres password=root dbname=churn_db port=5432 sslmode=disable TimeZone=America/Mexico_City"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Fallo al conectar con la base de datos: \n", err)
	}
	log.Println("Conexión a PostgreSQL.")
	DB = db
}