package main

import (
	"fmt"
	"html/template"
	//"log"
	"net/http"
	"os"
  "golang.org/x/crypto/bcrypt"
	"strings"
	"strconv"

	//"github.com/shirou/gopsutil/mem"

	"crypto/rand"
  "encoding/base64"
	"github.com/gorilla/sessions"

  "database/sql"
  //_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"

)

var tpl = template.Must(template.ParseFiles("index.html"))

var (
    // key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
		token=getToken(32)
    key = []byte(token)
    store = sessions.NewCookieStore(key)
)
type UserStruct struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
		Email string `json:"email"`
		PasswordHash string `json:"passwordHash"`
		Level int `json:"level"`
}

type GameServers struct {
		Id   int    `json:"id"`
		Type string `json:"type"`
		Name string `json:"name"`
		Localization string `json:"localization"`
		StartCommands string `json:"startcommands"`
		IsItOn bool `json:"isiton"`
}


func doQuery(query string, db *sql.DB){
	fmt.Println(query)
	statement, _ := db.Prepare(query)
	statement.Exec()
}

func getServersFromDataBase(db *sql.DB, where string) ([]GameServers, bool) {
	//db, _ := sql.Open("sqlite3", "./database.db")
	results, err := db.Query("SELECT * FROM servers" +where)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	var gameServers []GameServers
	var ifExists bool
	for results.Next() {
		var gameServer GameServers
		// for each row, scan the result into our tag composite object
		err = results.Scan(&gameServer.Id, &gameServer.Type, &gameServer.Name, &gameServer.Localization, &gameServer.StartCommands, &gameServer.IsItOn)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		} else {
			ifExists = true
		}
			gameServers = append(gameServers,gameServer)
	}
		return gameServers, ifExists
}

func getUsersFromDataBase(db *sql.DB, where string) ([]UserStruct, bool) {
	//db, _ := sql.Open("sqlite3", "./database.db")
	results, err := db.Query("SELECT * FROM users" +where)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	var usersResults []UserStruct
	var ifExists bool
	for results.Next() {
		var user UserStruct
		// for each row, scan the result into our tag composite object
		err = results.Scan(&user.Id, &user.Name, &user.PasswordHash, &user.Email, &user.Level)
		if err != nil {
		//	panic(err.Error()) // proper error handling instead of panic in your app
		ifExists = false
		} else {
			ifExists = true
		}
			usersResults = append(usersResults,user)
	}

		return usersResults, ifExists
}

//var user UserStruct
//err = db.QueryRow("SELECT passwordHash FROM users WHERE name = ?", strings.ToLower(nameFromForm)).Scan(&user.PasswordHash)

func getToken(length int) string {
    randomBytes := make([]byte, 32)
    _, err := rand.Read(randomBytes)
    if err != nil {
        panic(err)
    }
    return base64.StdEncoding.EncodeToString(randomBytes)[:length]
}

func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

func CheckPasswordHash(hash, password  string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	 //fmt.Println(r.URL.Path)
	 session, _ := store.Get(r, "cookie-name")

    // Check if user is authenticated
    if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
				if r.URL.Path != "/" {
	 		 		http.Redirect(w, r, "", 302)
	 	 		}
        login(w,r)
        return
    }
		fmt.Println(session.Values["name"])
		db, err := sql.Open("sqlite3", "./database.db")
		if err != nil {
				panic(err.Error())
		}

		switch (r.URL.Path) {
		case "/users" :

			var usersResults []UserStruct
			usersResults, _ = getUsersFromDataBase(db, "")

			if r.Method == http.MethodPost {
				s:= strings.Split(r.FormValue("submit"), "_")

				switch (s[0]) {
				case "add" :
					var ifExists bool
					_, ifExists = getUsersFromDataBase(db, (" WHERE name='"+ strings.ToLower(r.FormValue("name"))+ "'"))
					if (!ifExists) {
						hash, _ := HashPassword(r.FormValue("password"))
						doQuery("INSERT INTO users(name, passwordHash, email, level) VALUES('" +strings.ToLower(r.FormValue("name"))+ "', '" +hash+ "', '" +strings.ToLower(r.FormValue("email"))+ "', " +r.FormValue("level")+ ")", db )
					}
					http.Redirect(w, r, "/users", 302)

				case "edit":
					var hash string
					if (r.FormValue("password") != "" ) {
						hash, _ = HashPassword(r.FormValue("password"))
					} else {
						idFromSubmit, _ := strconv.ParseInt(s[1], 0, 64)
						hash = usersResults[idFromSubmit-1].PasswordHash
					}
					doQuery("UPDATE users SET name='"+strings.ToLower(r.FormValue("name"))+ "', passwordHash='"+hash+ "', email='" +strings.ToLower(r.FormValue("email"))+ "', level=" +r.FormValue("level")+ " WHERE id=" +s[1], db)
					http.Redirect(w, r, "/users", 302)

				case "delete" :
					if (session.Values["id"] != s[1]) {
						doQuery("DELETE FROM users WHERE id="+s[1], db)
					}
					http.Redirect(w, r, "/users", 302)

				default:
					http.Redirect(w, r, "/users", 302)
				}
			}
			// defer the close till after the main function has finished
			// executing
			defer db.Close()

			//fmt.Println(usersResults[0].Name)
			tpl = template.Must(template.ParseFiles("users.html"))
			tpl.Execute(w, usersResults)

		case "/logout" :
			// Revoke users authentication
			session.Values["authenticated"] = false
			session.Values["name"] = ""
			session.Values["id"] = ""
			session.Save(r, w)
			http.Redirect(w, r, "", 302)
		case "/" :
			var gameServersResults []GameServers
			gameServersResults, _ = getServersFromDataBase(db, "")




			tpl = template.Must(template.ParseFiles("index.html"))
			tpl.Execute(w, gameServersResults)
		default:
			http.Redirect(w, r, "", 302)
		}
}

func login(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			tpl = template.Must(template.ParseFiles("login.html"))
			tpl.Execute(w, nil)
			return
	  }

		session, _ := store.Get(r, "cookie-name")

		db, err := sql.Open("sqlite3", "./database.db")
		if err != nil {
				panic(err.Error())
		}

		var userResults []UserStruct
		var ifExists bool
		userResults, ifExists = getUsersFromDataBase(db, (" WHERE name='"+ strings.ToLower(r.FormValue("login"))+ "'"))
		//usersResults.PasswordHash
		if r.FormValue("login") != "" && ifExists && CheckPasswordHash(userResults[0].PasswordHash, r.FormValue("password")) {
			session.Values["authenticated"] = true
			session.Values["name"] = strings.ToLower(r.FormValue("login"))
			session.Values["id"] = strconv.Itoa(userResults[0].Id)
			session.Save(r, w)
	 	}

		http.Redirect(w, r, "", 302)
}
func main() {
	//v, _ := mem.VirtualMemory()
//	fmt.Printf("Total: %v, Free:%v, UsedPercent:%f%%\n", v.Total, v.Free, v.UsedPercent)

    // convert to JSON. String() is also implemented
  //  fmt.Println(v)

	//command := "Tell Application \"iTunes\" to playpause"

	    //c := exec.Command("/usr/bin/osascript", "-e", command)
		//	if err := c.Run(); err != nil {
	//		 fmt.Println(err.String())
	// }


	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))

	mux.HandleFunc("/", indexHandler)
	http.ListenAndServe(":"+port, mux)
}
