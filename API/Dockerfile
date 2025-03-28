# Imagen Go
FROM golang:1.24.1-alpine

# Instalar PostgreSQL Client para tener el comando pg_isready
RUN apk add --no-cache postgresql-client

# Configurar el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copiar los archivos go.mod y go.sum al contenedor
COPY go.mod go.sum ./ 

# Descargar dependencias de Go
RUN go get -d ./... && go mod tidy

# Copiar el resto de los archivos del proyecto
COPY . .

# Copiar el archivo start.sh y dar permisos de ejecución
COPY start.sh /app/start.sh
RUN chmod +x /app/start.sh

EXPOSE 8080

#Correr backend
CMD ["sh", "/app/start.sh"]
