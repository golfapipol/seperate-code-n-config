package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Port        int
	DatabaseURI string
}

func Connect(url string) (*sql.DB, error) {
	db, err := sql.Open("mysql", url)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func main() {
	env := os.Getenv("ENVIRONMENT") //flag.String("env", "development", "environment")
	var config Config
	fileName := fmt.Sprintf("./configs/%s.yaml", env)
	configData, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal("cannot config")
	}

	err = yaml.Unmarshal(configData, &config)
	if err != nil {
		log.Println("cannot unmarshal config")
	}

	databaseConnection, err := Connect(config.DatabaseURI)
	if err != nil {
		log.Fatal("cannot connect database ", err.Error())
	}
	defer databaseConnection.Close()

	http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {

		countTeamMember(databaseConnection, memberId)

	})

	api := API{
		DatabaseConnection: databaseConnection,
	}
	http.HandleFunc("/bar", api.handler)

	port := fmt.Sprintf(":%d", config.Port)
	log.Fatal(http.ListenAndServe(port, nil))
}

type API struct {
	DatabaseConnection *sql.DB
}

func (api API) handler(w http.ResponseWriter, r *http.Request) {

	countTeamMember(API.DatabaseConnection, memberId)

}

func countTeamMember(db *sql.DB, memberId int) TeamMember {
	statement, err := db.Prepare(`SELECT COUNT(id) AS HigherThanLevel FROM members WHERE leader_id=? AND level >= ?`)
	if err != nil {
		panic(err.Error())
	}
	resultCountCountHigherPearlPup := statement.QueryRow(memberId, levelPearlPup)
	resultCountHigherEmeraldPup := statement.QueryRow(memberId, levelEmeraldPup)
	resultCountHigherRubyPup := statement.QueryRow(memberId, levelRubyPup)

	var teamMember TeamMember

	err = resultCountCountHigherPearlPup.Scan(&teamMember.HigherPearl)
	if err != nil {
		panic(err.Error())
	}

	err = resultCountHigherEmeraldPup.Scan(&teamMember.HigherEmerald)
	if err != nil {
		panic(err.Error())
	}

	err = resultCountHigherRubyPup.Scan(&teamMember.HigherRuby)
	if err != nil {
		panic(err.Error())
	}

	return teamMember
}
