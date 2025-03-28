package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var db *sql.DB

// Estructura de datos para los partidos
type Match struct {
	ID    int    `json:"id"`
	Team1 string `json:"team1"`
	Team2 string `json:"team2"`
	Score string `json:"score"`
}

// Conexión a la base de datos
func initDB() {
	var err error
	// Conectar a la base de datos PostgreSQL
	connStr := "postgres://postgres:password@db:5432/matches?sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
}

// Manejo de errores para las consultas
func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// Endpoint para obtener todos los partidos
func getMatches(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, team1, team2, score FROM matches")
	checkErr(err)
	defer rows.Close()

	var matches []Match
	for rows.Next() {
		var match Match
		if err := rows.Scan(&match.ID, &match.Team1, &match.Team2, &match.Score); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		matches = append(matches, match)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(matches)
}

// Endpoint para obtener un partido por ID
func getMatchByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var match Match
	err := db.QueryRow("SELECT id, team1, team2, score FROM matches WHERE id = $1", id).Scan(&match.ID, &match.Team1, &match.Team2, &match.Score)
	if err != nil {
		http.Error(w, "Match not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(match)
}

// Endpoint para actualizar un partido existente
func updateMatch(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var match Match
	if err := json.NewDecoder(r.Body).Decode(&match); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := db.Exec("UPDATE matches SET team1 = $1, team2 = $2, score = $3 WHERE id = $4", match.Team1, match.Team2, match.Score, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Endpoint para crear un nuevo partido
func createMatch(w http.ResponseWriter, r *http.Request) {
	var match Match
	if err := json.NewDecoder(r.Body).Decode(&match); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := db.QueryRow("INSERT INTO matches (team1, team2, score) VALUES ($1, $2, $3) RETURNING id", match.Team1, match.Team2, match.Score).Scan(&match.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(match)
}

// Endpoint para eliminar un partido
func deleteMatch(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	_, err := db.Exec("DELETE FROM matches WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Configuración del servidor y rutas
func main() {
	// Inicializar la base de datos
	initDB()
	defer db.Close()

	// Crear el enrutador
	router := mux.NewRouter()
	router.HandleFunc("/api/matches", getMatches).Methods("GET")
	router.HandleFunc("/api/matches/{id}", getMatchByID).Methods("GET")
	router.HandleFunc("/api/matches/{id}", updateMatch).Methods("PUT")
	router.HandleFunc("/api/matches", createMatch).Methods("POST")
	router.HandleFunc("/api/matches/{id}", deleteMatch).Methods("DELETE")

	// Iniciar el servidor
	fmt.Println("API escuchando en puerto 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
