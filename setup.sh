#!/bin/bash

echo "Comprobando requisitos del sistema..."

if ! command -v docker &> /dev/null
then
    echo "Docker no detectado. Descargando e instalando Docker..."
    curl -fsSL https://get.docker.com -o get-docker.sh
    sudo sh get-docker.sh
    rm get-docker.sh
    echo "Docker instalado correctamente."
else
    echo "Docker ya está instalado."
fi

if ! command -v docker-compose &> /dev/null
then
    echo "Docker Compose no detectado. Instalando..."
    sudo apt-get update
    sudo apt-get install -y docker-compose
    echo "Docker Compose instalado."
else
    echo "Docker Compose ya está instalado."
fi

echo "Construyendo y levantando los contenedores..."
sudo docker-compose down
sudo docker-compose up -d --build

echo "Entra a http://localhost"