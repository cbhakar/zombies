package main

import (
	"encoding/json"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func filterSurvivorForOutput(s survivorIn) survivorOut {
	return survivorOut{
		ID:     s.ID,
		Name:   s.Name,
		Age:    s.Age,
		Gender: s.Gender,
		Location: coordinates{
			Latitude:  s.Latitude,
			Longitude: s.Longitude,
		},
		Resources:  s.Resources,
		IsInfected: s.IsInfected,
	}
}

func filterInfectedSurvivors(s []survivorIn) (notInfected, infected []survivorOut, notInfectedCount, infectedCount int) {
	notInfectedCount = 0
	infectedCount = 0
	for _, survivor := range s {
		if survivor.IsInfected {
			infectedCount++
			infected = append(infected, filterSurvivorForOutput(survivor))
		}else {
			notInfectedCount++
			notInfected = append(notInfected, filterSurvivorForOutput(survivor))
		}
	}
	return
}

func findPercentage(part, total int) float64 {
	return (float64(part) * float64(100)) / float64(total)
}
