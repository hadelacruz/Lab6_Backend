package main

import (
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "net/http"
    "time"
)

var db *gorm.DB

// Definición de la estructura de cada partido
type Match struct {
    ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
    HomeTeam  string    `json:"homeTeam" gorm:"type:varchar(255);not null"`
    AwayTeam  string    `json:"awayTeam" gorm:"type:varchar(255);not null"`
    MatchDate time.Time `json:"matchDate" gorm:"type:date;not null"`
}

func main() {
    dsn := "host=db user=postgres password=postgres dbname=matches port=5432 sslmode=disable"
    database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        panic("No se puede conectar a la base de datos")
    }
    db = database
    db.AutoMigrate(&Match{})

    r := gin.Default()

	//CORS
    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"*"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
    }))

    // rutas
    r.GET("/api/matches", GetMatches)
    r.GET("/api/matches/:id", GetMatch)
    r.POST("/api/matches", CreateMatch)
    r.PUT("/api/matches/:id", UpdateMatch)
    r.DELETE("/api/matches/:id", DeleteMatch)

    r.Run(":8080")
}

//http://localhost:8080/api/matches
// Obtener todos los partidos
func GetMatches(c *gin.Context) {
    var partidos []Match
    resultado := db.Find(&partidos)
    if resultado.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener los partidos"})
        return
    }
    c.JSON(http.StatusOK, partidos)
}

///http://localhost:8080/api/matches/:id
// Obtener un partido por ID
func GetMatch(c *gin.Context) {
    var partido Match
    id := c.Param("id")
    
    resultado := db.First(&partido, "id = ?", id)
    if resultado.Error != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Partido no encontrado"})
        return
    }
    c.JSON(http.StatusOK, partido)
}

//http://localhost:8080/api/matches
// Crear un nuevo partido
func CreateMatch(c *gin.Context) {
    var entrada struct {
        HomeTeam  string `json:"homeTeam"`
        AwayTeam  string `json:"awayTeam"`
        MatchDate string `json:"matchDate"`
    }
    
    if err := c.ShouldBindJSON(&entrada); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
        return
    }
    
    if entrada.HomeTeam == "" || entrada.AwayTeam == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Los equipos no pueden estar vacíos"})
        return
    }
        // Convertir string a time.Time

    fecha, err := time.Parse("2006-01-02", entrada.MatchDate)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de fecha inválido. Use YYYY-MM-DD"})
        return
    }
    
    partido := Match{
        HomeTeam:  entrada.HomeTeam,
        AwayTeam:  entrada.AwayTeam,
        MatchDate: fecha,
    }
    
    if resultado := db.Create(&partido); resultado.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo registrar el partido"})
        return
    }
    
    c.JSON(http.StatusCreated, gin.H{"mensaje": "Partido creado exitosamente", "partido": partido})
}

//http://localhost:8080/api/matches/:id
// Actualizar un partido
func UpdateMatch(c *gin.Context) {
    var partido Match
    id := c.Param("id")
    
    if resultado := db.First(&partido, "id = ?", id); resultado.Error != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Partido no encontrado"})
        return
    }
    
    var entrada struct {
        HomeTeam  string `json:"homeTeam"`
        AwayTeam  string `json:"awayTeam"`
        MatchDate string `json:"matchDate"`
    }
    
    if err := c.ShouldBindJSON(&entrada); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
        return
    }
    
    fecha, err := time.Parse("2006-01-02", entrada.MatchDate)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de fecha inválido. Use YYYY-MM-DD"})
        return
    }
    
    actualizar := map[string]interface{}{
        "HomeTeam":  entrada.HomeTeam,
        "AwayTeam":  entrada.AwayTeam,
        "MatchDate": fecha,
    }
    
    if resultado := db.Model(&partido).Updates(actualizar); resultado.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo actualizar el partido"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"mensaje": "Partido actualizado correctamente", "partido": partido})
}

//http://localhost:8080/api/matches/:id
// Eliminar un partido
func DeleteMatch(c *gin.Context) {
    id := c.Param("id")
    
    if resultado := db.Delete(&Match{}, "id = ?", id); resultado.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo eliminar el partido"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"mensaje": "Partido eliminado correctamente"})
}
