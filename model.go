package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
)

type coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type survivorOut struct {
	ID         int         `json:"id"`
	Name       string      `json:"name"`
	Age        int         `json:"age"`
	Gender     string      `json:"gender"`
	Location   coordinates `json:"last_known_location"`
	Resources  []string    `json:"resources"`
	IsInfected bool        `json:"is_infected"`
}

type survivorIn struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	Age        int     `json:"age"`
	Gender     string  `json:"gender"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	Resources  []string  `json:"resources"`
	IsInfected bool    `json:"is_infected"`
	Reported   int     `json:"reported"`
}

type report struct {
	PercentageOfInfected    float64       `json:"percentage_of_infected_survivors"`
	PercentageOfNonInfected float64       `json:"percentage_of_non_infected_survivors"`
	InfectedSurvivors       []survivorOut `json:"infected_survivors"`
	NonInfectedSurvivors    []survivorOut `json:"non_infected_survivors"`
	Robots                  []robot       `json:"robots"`
}

func NewPostgresConnection(host, port, user, password, dbname string) *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	conn, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	return conn
}


func (s *survivorIn) getSurvivor(db *sql.DB) error {
	var res string
	err :=  db.QueryRow("SELECT name, age, gender, latitude, longitude, resources, is_infected, reported   FROM survivors WHERE id=$1",
		s.ID).Scan(&s.Name, &s.Age, &s.Gender, &s.Latitude, &s.Longitude, &res, &s.IsInfected, &s.Reported)
	s.Resources = strings.Split(res, ",")
	return err
}

func (s *survivorIn) updateSurvivor(db *sql.DB) error {
	_, err := db.Exec("UPDATE survivors SET latitude=$1, longitude=$2 WHERE id=$3",
		s.Latitude, s.Longitude, s.ID)

	return err
}

func (s *survivorIn) updateInfectedSurvivor(db *sql.DB) error {
	_, err := db.Exec("UPDATE survivors SET is_infected=$1, reported=$2 WHERE id=$3",
		s.IsInfected, s.Reported, s.ID)

	return err
}

func (s *survivorIn) deleteSurvivor(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM survivors WHERE id=$1", s.ID)

	return err
}

func (s *survivorIn) addSurvivor(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO survivors(name, age, gender, latitude, longitude, resources, is_infected) VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id",
		s.Name, s.Age, s.Gender, s.Latitude, s.Longitude, strings.Join(s.Resources, ","), s.IsInfected).Scan(&s.ID)

	if err != nil {
		return err
	}

	return nil
}

func getSurvivors(db *sql.DB) ([]survivorIn, error) {
	rows, err := db.Query(
		"SELECT id, name, age, gender, latitude, longitude, resources, is_infected, reported FROM survivors")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	survivors := []survivorIn{}

	for rows.Next() {
		var s survivorIn
		var res string
		if err := rows.Scan(&s.ID, &s.Name, &s.Age, &s.Gender, &s.Latitude, &s.Longitude, &res, &s.IsInfected, &s.Reported); err != nil {
			return nil, err
		}
		s.Resources = strings.Split(res, ",")
		survivors = append(survivors, s)
	}

	return survivors, nil
}

func (a *App) ensureTableExists() {
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS survivors
(
    id SERIAL,
    name TEXT NOT NULL UNIQUE,
    age int NOT NULL,
    gender TEXT NOT NULL,
	latitude float NOT NULL,
	longitude float NOT NULL,
    resources TEXT NOT NULL,
	is_infected BOOLEAN NOT NULL DEFAULT false,
    reported int NOT NULL DEFAULT 0,
    CONSTRAINT survivors_pkey PRIMARY KEY (id)
)`
