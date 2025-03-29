package main

import (
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "net/http"
    "time"
	"github.com/swaggo/gin-swagger"
    //"github.com/swaggo/gin-swagger/swaggerFiles"
    "github.com/swaggo/files"

)

var db *gorm.DB

// Definición de la estructura de cada partido
type Match struct {
    ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
    HomeTeam  string    `json:"homeTeam" gorm:"type:varchar(255);not null"`
    AwayTeam  string    `json:"awayTeam" gorm:"type:varchar(255);not null"`
    MatchDate time.Time `json:"matchDate" gorm:"type:date;not null"`
	Goals            int       `json:"goals" gorm:"default:0"`           
    YellowCards      int       `json:"yellowCards" gorm:"default:0"`    
    RedCards         int       `json:"redCards" gorm:"default:0"`         
    ExtraTime        int       `json:"extraTime" gorm:"default:0"`        
}

// @title Matches API
// @version 1.0
// @description API para gestionar partidos de fútbol con estadísticas.
// @host localhost:8080
// @BasePath /api
func main() {
    dsn := "host=db user=postgres password=postgres dbname=matches port=5432 sslmode=disable"
    database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        panic("No se puede conectar a la base de datos")
    }
    db = database
    db.AutoMigrate(&Match{})

    r := gin.Default()
	//r.GET("/swagger/*any", ginSwagger.WrapHandler(files.Handler))
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))


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

	//Nuevas rutas
	r.PATCH("/api/matches/:id/goals", UpdateGoals)
	r.PATCH("/api/matches/:id/yellowcards", UpdateYellowCards)
	r.PATCH("/api/matches/:id/redcards", UpdateRedCards)
	r.PATCH("/api/matches/:id/extratime", UpdateExtraTime)


    r.Run(":8080")
}

// @Summary Obtener todos los partidos
// @Description Devuelve un listado con todos los partidos registrados
// @Tags Matches
// @Produce json
// @Success 200 {array} Match
// @Router /matches [get]
func GetMatches(c *gin.Context) {
    var partidos []Match
    resultado := db.Find(&partidos)
    if resultado.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener los partidos"})
        return
    }
    c.JSON(http.StatusOK, partidos)
}

// @Summary Obtener un partido por ID
// @Description Devuelve la información de un partido específico
// @Tags Matches
// @Produce json
// @Param id path int true "ID del partido"
// @Success 200 {object} Match
// @Failure 404 {object} map[string]string
// @Router /matches/{id} [get]
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

// @Summary Crear un nuevo partido
// @Description Registra un nuevo partido con los datos proporcionados
// @Tags Matches
// @Accept json
// @Produce json
// @Param partido body Match true "Datos del partido"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /matches [post]
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

// @Summary Actualizar un partido
// @Description Modifica los datos de un partido existente
// @Tags Matches
// @Accept json
// @Produce json
// @Param id path int true "ID del partido"
// @Param partido body Match true "Datos actualizados del partido"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /matches/{id} [put]
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

// @Summary Eliminar un partido
// @Description Elimina un partido de la base de datos
// @Tags Matches
// @Param id path int true "ID del partido"
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /matches/{id} [delete]
func DeleteMatch(c *gin.Context) {
    id := c.Param("id")
    
    if resultado := db.Delete(&Match{}, "id = ?", id); resultado.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo eliminar el partido"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"mensaje": "Partido eliminado correctamente"})
}


// @Summary Actualizar goles de un partido
// @Description Modifica la cantidad de goles registrados en un partido
// @Tags Matches
// @Accept json
// @Produce json
// @Param id path int true "ID del partido"
// @Param goles body int true "Cantidad de goles"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /matches/{id}/goals [patch]
func UpdateGoals(c *gin.Context) {
    var partido Match
    id := c.Param("id")

    if resultado := db.First(&partido, "id = ?", id); resultado.Error != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Partido no encontrado"})
        return
    }

    var entrada struct {
        Goals int `json:"goals"`
    }

    if err := c.ShouldBindJSON(&entrada); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
        return
    }

    partido.Goals = entrada.Goals

    if resultado := db.Save(&partido); resultado.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron actualizar los goles"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"mensaje": "Goles actualizados correctamente", "partido": partido})
}

// @Summary Actualizar tarjetas amarillas
// @Description Modifica la cantidad de tarjetas amarillas en un partido
// @Tags Matches
// @Accept json
// @Produce json
// @Param id path int true "ID del partido"
// @Param yellowCards body int true "Cantidad de tarjetas amarillas"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /matches/{id}/yellowcards [patch]
func UpdateYellowCards(c *gin.Context) {
    var partido Match
    id := c.Param("id")

    if resultado := db.First(&partido, "id = ?", id); resultado.Error != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Partido no encontrado"})
        return
    }

    var entrada struct {
        YellowCards int `json:"yellowCards"`
    }

    if err := c.ShouldBindJSON(&entrada); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
        return
    }

    partido.YellowCards = entrada.YellowCards

    if resultado := db.Save(&partido); resultado.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron actualizar las tarjetas amarillas"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"mensaje": "Tarjetas amarillas actualizadas correctamente", "partido": partido})
}

// @Summary Actualizar tarjetas rojas
// @Description Modifica la cantidad de tarjetas rojas en un partido
// @Tags Matches
// @Accept json
// @Produce json
// @Param id path int true "ID del partido"
// @Param redCards body int true "Cantidad de tarjetas rojas"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /matches/{id}/redcards [patch]
func UpdateRedCards(c *gin.Context) {
    var partido Match
    id := c.Param("id")

    if resultado := db.First(&partido, "id = ?", id); resultado.Error != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Partido no encontrado"})
        return
    }

    var entrada struct {
        RedCards int `json:"redCards"`
    }

    if err := c.ShouldBindJSON(&entrada); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
        return
    }

    partido.RedCards = entrada.RedCards

    if resultado := db.Save(&partido); resultado.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron actualizar las tarjetas rojas"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"mensaje": "Tarjetas rojas actualizadas correctamente", "partido": partido})
}


// @Summary Registrar tiempo extra
// @Description Modifica la cantidad de tiempo extra en minutos en un partido
// @Tags Matches
// @Accept json
// @Produce json
// @Param id path int true "ID del partido"
// @Param extraTime body int true "Minutos de tiempo extra"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /matches/{id}/extratime [patch]
func UpdateExtraTime(c *gin.Context) {
    var partido Match
    id := c.Param("id")
    
    if resultado := db.First(&partido, "id = ?", id); resultado.Error != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Partido no encontrado"})
        return
    }
    
    var entrada struct {
        ExtraTime int `json:"extraTime"`
    }
    
    if err := c.ShouldBindJSON(&entrada); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
        return
    }
    
    partido.ExtraTime = entrada.ExtraTime
    
    if resultado := db.Save(&partido); resultado.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo registrar el tiempo extra"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"mensaje": "Tiempo extra registrado correctamente", "partido": partido})
}


