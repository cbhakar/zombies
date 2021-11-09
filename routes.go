package main

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"net/http"
	"strconv"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) initializeRoutes() {
	a.Router = mux.NewRouter()

	a.Router.HandleFunc("/survivors", a.getSurvivors).Methods("GET")
	a.Router.HandleFunc("/survivor", a.addSurvivor).Methods("POST")
	a.Router.HandleFunc("/survivor/{id:[0-9]+}", a.getSurvivor).Methods("GET")
	a.Router.HandleFunc("/survivor/{id:[0-9]+}", a.updateSurvivor).Methods("PUT")
	a.Router.HandleFunc("/survivor/{id:[0-9]+}", a.deleteSurvivor).Methods("DELETE")
	a.Router.HandleFunc("/report/survivor/{id:[0-9]+}", a.flagInfectedSurvivors).Methods("POST")
	a.Router.HandleFunc("/report", a.getReports).Methods("GET")
}

func (a *App) getReports(w http.ResponseWriter, r *http.Request) {
	survivors, err := getSurvivors(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	notInfected, infected, notInfectedCount, infectedCount := filterInfectedSurvivors(survivors)
	robots, err := hackIntoDbForRobotDetails()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	report := report{
		PercentageOfInfected:    findPercentage(infectedCount, notInfectedCount+infectedCount),
		PercentageOfNonInfected: findPercentage(notInfectedCount, notInfectedCount+infectedCount),
		InfectedSurvivors:       infected,
		NonInfectedSurvivors:    notInfected,
		Robots:                  robots,
	}

	respondWithJSON(w, http.StatusOK, report)
}

func (a *App) getSurvivor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Survivor ID")
		return
	}

	s := survivorIn{ID: id}
	if err := s.getSurvivor(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "survivor not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, filterSurvivorForOutput(s))
}

func (a *App) getSurvivors(w http.ResponseWriter, r *http.Request) {
	survivors, err := getSurvivors(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	var outSurvivors []survivorOut
	for _, survivor := range survivors {
		s := filterSurvivorForOutput(survivor)
		outSurvivors = append(outSurvivors, s)
	}
	if len(outSurvivors) == 0{
		respondWithError(w, http.StatusNotFound, "survivors not found")
	}
	respondWithJSON(w, http.StatusOK, outSurvivors)
}

func (a *App) flagInfectedSurvivors(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Survivor ID")
		return
	}

	s := survivorIn{ID: id}
	if err := s.getSurvivor(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "survivor not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	if s.IsInfected {
		respondWithError(w, http.StatusNotFound, "survivor is already infected")
		return
	}
	s.Reported += 1
	if s.Reported >= 3 {
		s.IsInfected = true
	}
	if err := s.updateInfectedSurvivor(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "survivor flagged successfully"})
}

func (a *App) addSurvivor(w http.ResponseWriter, r *http.Request) {
	var s survivorIn
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&s); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := s.addSurvivor(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]int{"survivor_id": s.ID})
}

func (a *App) updateSurvivor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Survivor ID")
		return
	}

	var s survivorIn
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&s); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	s.ID = id

	if err := s.updateSurvivor(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) deleteSurvivor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Survivor ID")
		return
	}

	s := survivorIn{ID: id}
	if err := s.getSurvivor(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "survivor not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	if err := s.deleteSurvivor(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
