package models

import (
	"time"
)

type Hiperparametro struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	TasaAprendizaje  float64 `gorm:"not null"`
	ProfundidadArbol int     `gorm:"not null"`
	Estimadores      int     `gorm:"not null"`
	WorkerAsignado string
	Estado         string
	PrecisionFinal float64 `gorm:"default:0.0"`
}