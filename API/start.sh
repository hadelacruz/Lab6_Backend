#!/bin/sh

# Esperar la BD
until pg_isready -h db -U postgres; do
  echo "Esperando por la base de datos..."
  sleep 3
done

# Ejecutar la aplicación
echo "BD lista, arrancando la API..."
go run main.go