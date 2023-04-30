package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

type Recipe struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	EnoughFor   string `json:"enough_for"`
	Origin      string `json:"origin"`
	Ingredients string `json:"ingredients"`
	Description string `json:"description"`
	Kind        string `json:"kind"`
	PrepTime    string `json:"prep_time"`
	Difficulty  string `json:"difficulty"`
	Notes       string `json:"notes"`
	CookTime    string `json:"cook_time"`
	ServingSize string `json:"serving_size"`
	Rating      string `json:"rating"`
}

var db *sql.DB

func initDB() {
	user := viper.GetString("mysql.user")
	password := viper.GetString("mysql.password")
	host := viper.GetString("mysql.host")
	port := viper.GetInt("mysql.port")
	database := viper.GetString("mysql.database")

	dataSourceName := buildDSN(user, password, host, port, database)

	var err error
	db, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		fmt.Printf("Error opening database connection: %s\n", err)
		redactedDSN := buildDSN(user, "REDACTED", host, port, database)
		fmt.Printf("DSN: %s\n", redactedDSN)
		os.Exit(1)
	}

	err = db.Ping()
	if err != nil {
		fmt.Printf("Error pinging database: %s\n", err)
		redactedDSN := buildDSN(user, "REDACTED", host, port, database)
		fmt.Printf("DSN: %s\n", redactedDSN)
		os.Exit(1)
	}
}

func buildDSN(user, password, host string, port int, database string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", user, password, host, port, database)
}
