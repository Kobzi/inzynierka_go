package main

import (
    //"fmt"
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
	"log"
    //"strconv"
)

type UsersStruct struct {
    id   int    `json:"id"`
    name string `json:"name"`
    passwordHash string `json:"passwordHash"`
		email string `json:"email"`
		level int `json:"level"`
}

func main() {
  db, _ := sql.Open("sqlite3", "./database.db")

  results, err := db.Query("SELECT id, name, passwordHash, email, level FROM users")
   if err != nil {
       panic(err.Error()) // proper error handling instead of panic in your app
   }
   var usersResults[] UsersStruct




    // perform a db.Query insert



}
