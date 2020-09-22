package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"

	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
	"github.com/tobgu/qframe"
	_ "github.com/tobgu/qframe"               // import qframe library
	qsql "github.com/tobgu/qframe/config/sql" // import qFrame with SQL
	_ "gopkg.in/yaml.v2"                      // Import yaml library
)

type conf struct {
	SQLType      string `yaml:"type"`
	SQLStatememt string `yaml:"sql"`
}

func main() {
	//os.Remove("sqlite-database.db") // I delete the file to avoid duplicated records.
	// SQLite is a file based database.

	log.Println("Creating sqlite-database.db...")
	if _, err := os.Stat("sqlite-database.db"); err == nil || os.IsExist(err) {
		log.Println("Database file already exisits")
	} else {
		file, err := os.Create("sqlite-database.db") // Create SQLite file
		if err != nil {
			log.Fatal(err.Error())
		}
		file.Close()
		log.Println("sqlite-database.db created")
	}

	sqliteDatabase, _ := sql.Open("sqlite3", "./sqlite-database.db") // Open the created SQLite File
	defer sqliteDatabase.Close()                                     // Defer Closing the database
	createTable(sqliteDatabase)                                      // Create Database Tables
	// INSERT RECORDS
	insertStudent(sqliteDatabase, "0001", "Liana Kim", "Bachelor")
	insertStudent(sqliteDatabase, "0002", "Glen Rangel", "Bachelor")

	// DISPLAY INSERTED RECORDS
	displayStudents(sqliteDatabase)

	dat, err := ioutil.ReadFile("query.sql")
	check(err)
	fmt.Println(string(dat))

	c, err := readConf("conf.yaml")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(c.SQLStatememt) // check it is only select not /\b(insert|update|drop|delete|create)\b/i
	re := regexp.MustCompile(`(?i).*(insert|update|drop|delete|create).*`)
	x := re.FindStringIndex(c.SQLStatememt)
	if x != nil {
		fmt.Printf("Not pure Select, do not execute it\n")
	} else {
		fmt.Println(c.SQLStatememt)
	}
	// Display output as data frame
	tx, _ := sqliteDatabase.Begin()
	Qf := qframe.ReadSQL(tx, qsql.Query(c.SQLStatememt), qsql.SQLite())
	fmt.Println(Qf)
}
