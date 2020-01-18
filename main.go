package main

import (
	//"fmt"
	"html/template"
	//"log"
	"net/http"
	"os"
  "golang.org/x/crypto/bcrypt"

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

func findPassword(nameFromForm string) string{
		//b, err := sql.Open("mysql", "adminPanel:adminPanelPassword@tcp(127.0.0.1:3306)/adminPanel")
		db, err := sql.Open("sqlite3", "./database.db")
		if err != nil {
				panic(err.Error())
		}
		// defer the close till after the main function has finished
		// executing
		defer db.Close()

		var user UserStruct
		err = db.QueryRow("SELECT passwordHash FROM users WHERE name = ?", nameFromForm).Scan(&user.PasswordHash)
		if err != nil {
	    panic(err.Error()) // proper error handling instead of panic in your app
		}
		return user.PasswordHash
}

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

		switch (r.URL.Path) {
		case "/users" :
			db, _ := sql.Open("sqlite3", "./database.db")

		  results, err := db.Query("SELECT id, name, email, level FROM users")
		  if err != nil {
		  	panic(err.Error()) // proper error handling instead of panic in your app
		  }
		  var usersResults []UserStruct

		  for results.Next() {
 				var user UserStruct
				// for each row, scan the result into our tag composite object
				err = results.Scan(&user.Id, &user.Name, &user.Email, &user.Level)
				if err != nil {
					panic(err.Error()) // proper error handling instead of panic in your app
				}
				usersResults = append(usersResults,user)
		  }
			//fmt.Println(usersResults[0].Name)
			tpl = template.Must(template.ParseFiles("users.html"))
			tpl.Execute(w, usersResults)
		case "/logout" :
			// Revoke users authentication
			session.Values["authenticated"] = false
			session.Save(r, w)
			http.Redirect(w, r, "", 302)
		case "/" :
			tpl = template.Must(template.ParseFiles("index.html"))
			tpl.Execute(w, nil)
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
		//db, err := sql.Open("sqlite3", ":memory:")
	  //password := "admin"
	  //hash, _ := HashPassword(password)
		//log.Printf(hash)

		if r.FormValue("login") != "" && CheckPasswordHash(findPassword(r.FormValue("login")), r.FormValue("password")) {
			session.Values["authenticated"] = true
			session.Save(r, w)
	 	}
		http.Redirect(w, r, "", 302)
}
func main() {
	//v, _ := mem.VirtualMemory()
//	fmt.Printf("Total: %v, Free:%v, UsedPercent:%f%%\n", v.Total, v.Free, v.UsedPercent)

    // convert to JSON. String() is also implemented
  //  fmt.Println(v)
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
