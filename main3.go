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

  //"database/sql"
  //     _ "github.com/mattn/go-sqlite3"
)

var tpl = template.Must(template.ParseFiles("index.html"))

var (
    // key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
		token=getToken(32)
    key = []byte(token)
    store = sessions.NewCookieStore(key)
)

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
	session, _ := store.Get(r, "cookie-name")

    // Check if user is authenticated
    if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
        login(w,r)
        return
    }
	tpl = template.Must(template.ParseFiles("index.html"))
	tpl.Execute(w, nil)
}

func login(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		tpl = template.Must(template.ParseFiles("login.html"))
		tpl.Execute(w, nil)
            return
  }
	session, _ := store.Get(r, "cookie-name")


//db, err := sql.Open("sqlite3", ":memory:")


  password := "admin"
  hash, _ := HashPassword(password)

 if r.FormValue("login") == "admin" && CheckPasswordHash(hash, r.FormValue("password")) {
	 session.Values["authenticated"] = true
	 session.Save(r, w)
}

	http.Redirect(w, r, "", 302)
}

func logout(w http.ResponseWriter, r *http.Request) {
    session, _ := store.Get(r, "cookie-name")

    // Revoke users authentication
			session.Values["authenticated"] = false
	    session.Save(r, w)

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
	mux.HandleFunc("/logout", logout)
	http.ListenAndServe(":"+port, mux)
}
