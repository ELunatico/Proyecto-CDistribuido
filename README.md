# Sistema Distribuido de Optimización de Modelos de Machine Learning

Este proyecto implementa una arquitectura distribuida basada en microservicios y contenedores Docker para la simulación y optimización de hiperparámetros en modelos de Machine Learning. El sistema distribuye cargas de trabajo pesado (entrenamiento) a través de múltiples *workers* independientes de manera asíncrona.

El sistema utiliza un diseño de múltiples capas operando de forma independiente dentro de una red privada de Docker:

1. **Nginx (API Gateway & Reverse Proxy):** 
   - Actúa como la única puerta de entrada pública (Puerto 80).
   - Sirve los archivos estáticos del Frontend (HTML/JS).
   - Implementa un Load Balancer básico mediante un bloque `upstream`, redirigiendo de forma segura el tráfico web (`/api/`) hacia el Middleware.
2. **Frontend:** Panel de control interactivo que consume la API REST y actualiza de manera reactiva el estado de las tareas sin bloquear el navegador.
3. **Middleware / Orquestador (Go):** 
   - Construido con **Gin** y conectado a **PostgreSQL** mediante **GORM**.
   - Administra el estado de las tareas (Procesando, Completado, Error).
   - Utiliza **Goroutines y Channels** para la ejecución concurrente y el manejo asíncrono de los resultados.
4. **Load Balancer Interno (Go):** Intermediario proxy inverso (`httputil.NewSingleHostReverseProxy`) encargado de distribuir equitativamente el tráfico de cómputo intensivo hacia los workers de Python disponibles.
5. **Workers de Procesamiento (Python/FastAPI):** Nueve (9) contenedores independientes que reciben los hiperparámetros, simulan la carga de entrenamiento (listos para integrar modelos reales con `scikit-learn` o `TensorFlow`) y devuelven la precisión final junto con su `Worker ID` (Hostname del contenedor) para fines de trazabilidad.

## Requisitos Previos

Para desplegar esta infraestructura, el equipo host debe contar con:
- **Entorno basado en Linux / Unix**
- **Docker Engine** (o Docker Desktop)
- **Docker Compose**

*(Nota: El script de automatización intentará instalar Docker y Docker Compose si no se detectan en el sistema).*

## Instrucciones de Ejecución

El proyecto incluye un script de automatización que gestiona la construcción y el despliegue de toda la arquitectura.

1. Clona este repositorio y navega a la carpeta principal.
2. Otorga permisos de ejecución al script de instalación (si es necesario):
   ```bash
   chmod +x setup.sh
