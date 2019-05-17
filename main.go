package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/nvellon/weather-predict-go/solarsystem"
)

var repo solarsystem.Repo

func init() {
	projectID := os.Getenv("DATASTORE_PROJECT_ID")
	if projectID == "" {
		log.Fatal(`You need to set the environment variable "DATASTORE_PROJECT_ID"`)
	}

	var err error
	repo, err = solarsystem.NewRepo(projectID)
	if err != nil {
		log.Fatalf("Could not create repo: %v", err)
	}
}

func main() {
	r := gin.Default()

	r.GET("/clima", GetWeather)
	r.GET("/totales", GetTotals)

	r.Run() // listen and serve on 0.0.0.0:8080
}

// GetWeather godoc
func GetWeather(c *gin.Context) {
	day, err := strconv.Atoi(c.DefaultQuery("dia", "1"))
	if err != nil {
		log.Printf("Failed to convert query param: %v", err)
		c.Status(http.StatusBadRequest)
		return
	}

	record, err := repo.GetByDay(day)
	if err != nil {
		log.Fatalf("Failed to fetch record by day: %v", err)
		c.Status(http.StatusNotFound)
		return
	}

	clima := "otro"
	if record.IsDrought {
		clima = "sequía"
	}
	if record.IsOptimumTemperaturePressure {
		clima = "OTPC"
	}
	if record.IsRainSeason {
		clima = "lluvia"
	}

	c.JSON(http.StatusOK, gin.H{
		"dia":   day,
		"clima": clima,
	})
}

// GetTotals godoc
func GetTotals(c *gin.Context) {
	counter, err := repo.GetCounter()
	if err != nil {
		log.Fatalf("Failed to fetch counter: %v", err)
		c.Status(http.StatusServiceUnavailable)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"dias":   counter.Days,
		"sequia": counter.CountDrought,
		"lluvia": counter.CountRainSeason,
		"otpc":   counter.CountOptimumTemperaturePressure,
		"otros":  counter.CountOther,
	})
}
