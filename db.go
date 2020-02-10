package main

import (
	"fmt"
	//"html/template"
	//"log"

  "database/sql"
  //_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"

)
func main() {
	fmt.Println("Test")

  db, _ := sql.Open("sqlite3", "./database2.db")
  statement, _ := db.Prepare("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, name TEXT, passwordHash TEXT, email TEXT, level INTEGER)")
  statement.Exec()
  statement, _ = db.Prepare("INSERT INTO users(name, passwordHash, email, level) VALUES('Kobzi', '$2a$14$F4QCrBRiz7mZh6/O2NmHa.D0lruHy6A7BZNMoXwWDJym20x3eJ1O2', 'test@test.com', 1)")
  statement.Exec()
//("INSERT INTO users(name, passwordHash, email, level) VALUES('Kobzi', '$2a$14$F4QCrBRiz7mZh6/O2NmHa.D0lruHy6A7BZNMoXwWDJym20x3eJ1O2', 'test@test.com', 1);")

}
