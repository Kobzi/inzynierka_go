package main

import (
    //"fmt"
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
	   "log"
    //"strconv"
)

type UserStruct struct {
    id   int    `json:"id"`
    name string `json:"name"`
    passwordHash string `json:"passwordHash"`
		email string `json:"email"`
		level int `json:"level"`
}

func main() {
  db, _ := sql.Open("sqlite3", "./database.db")

  results, err := db.Query("SELECT id, name, email, level FROM users")
   if err != nil {
       panic(err.Error()) // proper error handling instead of panic in your app
   }
   var usersResults []UserStruct

       for results.Next() {
           var user UserStruct
           // for each row, scan the result into our tag composite object
           err = results.Scan(&user.id, &user.name, &user.email, &user.level)
           if err != nil {
               panic(err.Error()) // proper error handling instead of panic in your app
           }
                   // and then print out the tag's Name attribute
          usersResults[1]=user
          log.Printf(user.name)
       }


    // perform a db.Query insert
//https://freshman.tech/web-development-with-go/
//https://www.thepolyglotdeveloper.com/2017/04/using-sqlite-database-golang-application/
//https://tutorialedge.net/golang/golang-mysql-tutorial/
//https://medium.com/@hugo.bjarred/mysql-and-golang-ea0d620574d2

}
