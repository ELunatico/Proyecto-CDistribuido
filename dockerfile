# Usar la imagen oficial de Go en su versión más reciente
FROM golang:alpine

# Crear directorio de trabajo
WORKDIR /app

# Copiar todo el código al contenedor
COPY . .

# Descargar dependencias y compilar la aplicación
RUN go mod tidy
RUN go build -o main .

# Exponer el puerto
EXPOSE 8080

# Comando para ejecutar la app
CMD ["./main"]