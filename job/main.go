package main

import (
	"fmt"
	"log"
	"os"

	"github.com/nvellon/weather-predict-go/solarsystem"
)

const (
	DryRun    = false
	TotalDays = 20 * 365
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
	counter := &solarsystem.Counter{
		Days:                            TotalDays,
		CountOther:                      0,
		CountDrought:                    0,
		CountOptimumTemperaturePressure: 0,
		CountRainSeason:                 0,
	}
	ss := solarsystem.NewSolarSystem()

	for i := 0; i < TotalDays; i++ {
		if !DryRun {
			record := &solarsystem.Record{
				Day:                          ss.GetDay(),
				IsDrought:                    ss.IsDrought(),
				IsOptimumTemperaturePressure: ss.IsOptimumTemperaturePressure(),
				IsRainSeason:                 ss.IsRainSeason(),
				FerengiLocation:              fmt.Sprintf("%f,%f", ss.Ferengi.Location.X, ss.Ferengi.Location.Y),
				BetasoideLocation:            fmt.Sprintf("%f,%f", ss.Betasoide.Location.X, ss.Betasoide.Location.Y),
				VulcanoLocation:              fmt.Sprintf("%f,%f", ss.Vulcano.Location.X, ss.Vulcano.Location.Y),
			}

			repo.Save(record)
		}

		isOther := true

		if ss.IsDrought() {
			counter.CountDrought++
			isOther = false
		}

		if ss.IsOptimumTemperaturePressure() {
			counter.CountOptimumTemperaturePressure++
			isOther = false
		}

		if ss.IsRainSeason() {
			counter.CountRainSeason++
			isOther = false
		}

		if isOther {
			counter.CountOther++
		}

		ss.NextDay()
	}

	repo.SaveCounter(counter)

	log.Printf("Days: %d", counter.Days)
	log.Printf("Others: %d", counter.CountOther)
	log.Printf("Drought: %d", counter.CountDrought)
	log.Printf("OptimumTemperaturePressure: %d", counter.CountOptimumTemperaturePressure)
	log.Printf("RainSeason: %d", counter.CountRainSeason)
}
