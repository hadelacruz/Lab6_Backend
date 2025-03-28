#!/bin/sh

until pg_isready -h db -U postgres; do
  echo "Esperando por la base de datos..."
  sleep 5
done

# Ejecutar la aplicación
echo "Base de datos lista, API arrancando..."
go run main.go
