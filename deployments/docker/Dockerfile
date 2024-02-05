# Usa la imagen oficial de Go como base
FROM golang:1.21.6-alpine

# Establece el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copia el archivo go.mod y go.sum para descargar dependencias
COPY go.mod go.sum ./

# Descarga las dependencias
RUN go mod download

# Copia el resto de los archivos al directorio de trabajo
COPY . .

# Compila la aplicación
# RUN go build -o main cmd/krowi/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/krowi/main.go

# Expone el puerto en el que la aplicación se ejecutará
EXPOSE 8000

# Comando para ejecutar la aplicación cuando el contenedor se inicia
CMD ["./main"]
