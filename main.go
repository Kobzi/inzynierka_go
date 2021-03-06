package main

import (
	"fmt"
	"html/template"
	//"log"
	"net/http"
	"os"
	"os/exec"
  "golang.org/x/crypto/bcrypt"
	"strings"
	"strconv"
	"encoding/json"
	"time"

	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/cpu"

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

type GamesStruct struct {
		IdGame   int    `json:"idgame"`
		GamePlatform string `json:"gameplatform"`
		GamePlatformId int `json:"gameplatformid"`
		NameShort string `json:"nameshort"`
		NameFull string `json:"namefull"`
		StandardCommands string `json:"standardcommands"`
		FileToRun string `json:"filetorun"`
}

type GameServers struct {
		Id   int    `json:"id"`
		Type string `json:"type"`
		Name string `json:"name"`
		Localization string `json:"localization"`
		StartCommands string `json:"startcommands"`
		IsItOn bool `json:"isiton"`
		AlreadyDownloaded bool `json:"alreadydownloaded"`
		Owner string `json:"owner"`
}

type HardwareStruct struct {
		//ID string `json:"id"`
		MemoryTotal string `json:"memorytotal"`
		MemoryUsed string `json:"memoryused"`

		MemoryPercentOfUsed string `json:"memorypercentofused"`
		HDDTotal string `json:"hddtotal"`
		HDDUsed string `json:"hddused"`
		HDDPercentOfUsed string `json:"hddpercentofused"`
	//	OS string `json:"os"`
		CPUUsage string `json:"cpuusage"`
}

func getHardwareInfo() HardwareStruct {
	v, _ := mem.VirtualMemory()
	hdd, _ := disk.Usage("/")
	cpuPercent, _ :=cpu.Percent(time.Second,false)

	hw := HardwareStruct {
		strconv.FormatUint(v.Total, 10),
		strconv.FormatUint(v.Used, 10),
		strconv.FormatFloat(v.UsedPercent, 'f', 2, 64),
		strconv.FormatUint(hdd.Total, 10),
		strconv.FormatUint(hdd.Used, 10),
		strconv.FormatFloat(hdd.UsedPercent, 'f', 2, 64),
		strconv.FormatFloat(cpuPercent[0], 'f', 2, 64) }
		return hw
}

func gameServer(whatToDo string, path string, gameId string, gameParametrs string, fileToRun string) int {
  processid:=0

  switch (whatToDo) {
  case "download" :
    cmd := exec.Command("./src/SteamCMD/steamcmd.exe", "+login anonymous", "+force_install_dir "+path, "+app_update "+gameId)
    cmd.Stdout = os.Stdout
    err := cmd.Start()
    if err != nil {
       panic(err.Error())
    }
    //log.Printf("Just ran subprocess %d, exiting\n", cmd.Process.Pid)
    processid = cmd.Process.Pid

  case "runGame" :
    cmd := exec.Command("cmd","/c", "cd /d "+path+" && start " +gameId+ " " +gameParametrs)
    cmd.Stdout = os.Stdout
    err := cmd.Start()
    cmd.Wait()
    if err != nil {
       panic(err.Error())
    }

    processid = cmd.Process.Pid
  }
	return processid
}

func doQuery(query string, db *sql.DB){
	fmt.Println(query)
	statement, _ := db.Prepare(query)
	statement.Exec()
}

func getServersFromDataBase(db *sql.DB, where string) ([]GameServers, bool) {
	//db, _ := sql.Open("sqlite3", "./database.db")
	results, err := db.Query("SELECT servers.id, servers.type, servers.nameServer, servers.localization, servers.startCommands, servers.isiton, servers.alreadydownloaded, organization.userId FROM servers INNER JOIN organization ON servers.id=organization.serverId " +where)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	var gameServers []GameServers
	var ifExists bool
	for results.Next() {
		var gameServer GameServers
		// for each row, scan the result into our tag composite object
		err = results.Scan(&gameServer.Id, &gameServer.Type, &gameServer.Name, &gameServer.Localization, &gameServer.StartCommands, &gameServer.IsItOn, &gameServer.AlreadyDownloaded, &gameServer.Owner)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		} else {
			ifExists = true
		}
			gameServers = append(gameServers,gameServer)
	}
		return gameServers, ifExists
}

func getGamesFromDataBase(db *sql.DB, where string) ([]GamesStruct, bool){
	results, err := db.Query("SELECT * FROM games " +where)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	var gamesResults []GamesStruct
	var ifExists bool

	for results.Next() {
		var games GamesStruct
		// for each row, scan the result into our tag composite object
		err = results.Scan(&games.IdGame, &games.GamePlatform, &games.GamePlatformId, &games.NameShort, &games.NameFull, &games.StandardCommands, &games.FileToRun )
		if err != nil {
		//	panic(err.Error()) // proper error handling instead of panic in your app
		ifExists = false
		} else {
			ifExists = true
		}
			gamesResults = append(gamesResults,games)
	}

		return gamesResults, ifExists
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
func apiHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("API")
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
	 case "/api/hardware" :
		 json.NewEncoder(w).Encode(getHardwareInfo())
	 }
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
		//fmt.Println(session.Values["name"])
		db, err := sql.Open("sqlite3", "./database.db")
		if err != nil {
				panic(err.Error())
		}

		switch (r.URL.Path) {
		case "/users" :

			var usersResults []UserStruct
			usersResults, _ = getUsersFromDataBase(db, (" WHERE level>" +session.Values["levelOfUser"].(string)))

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
			tpl.Execute(w, struct{
				Name string
				Users []UserStruct
			}{session.Values["name"].(string), usersResults})

		case "/logout" :
			// Revoke users authentication
			session.Values["authenticated"] = false
			session.Values["name"] = ""
			session.Values["id"] = ""
			session.Values["level"] = ""
			session.Save(r, w)
			http.Redirect(w, r, "", 302)


		case "/myprofile":

			var usersResults []UserStruct
			usersResults, _ = getUsersFromDataBase(db, (" WHERE name='" +session.Values["name"].(string)+ "'"))

			tpl = template.Must(template.ParseFiles("myprofile.html"))
			tpl.Execute(w, struct{
				Name string
				UserStruct UserStruct
			}{session.Values["name"].(string),usersResults[0]})


		case "/" :
			gameServersResults, _ := getServersFromDataBase(db, "WHERE organization.userId=" +session.Values["id"].(string))
			gamesResults, _ := getGamesFromDataBase(db,"")


			if r.Method == http.MethodPost {
				s:= strings.Split(r.FormValue("submit"), "_")
				fmt.Println(r.FormValue("submit"))
				switch (s[0]) {
					case "add" :
						var ifExists bool
						_, ifExists = getServersFromDataBase(db, (" WHERE localization='"+ strings.ToLower(r.FormValue("localization"))+ "'"))
						if (!ifExists) {
							doQuery("INSERT INTO servers(type, nameServer, localization, startCommands, isiton, alreadydownloaded) VALUES('" +strings.ToLower(r.FormValue("game"))+ "', '" +strings.ToLower(r.FormValue("name"))+ "', '" +strings.ToLower(r.FormValue("localization"))+ "', '" +strings.ToLower(r.FormValue("commandsToRunServer"))+ "', 0, 0)", db )
						}
						http.Redirect(w, r, "/", 302)
					case "edit" :
						doQuery("UPDATE servers SET nameServer='" +r.FormValue("name")+"', startCommands='" +r.FormValue("commandsToRunServer")+ "' WHERE id=" +s[1], db)
						http.Redirect(w, r, "/", 302)
					case "delete" :
						doQuery("DELETE FROM servers WHERE id="+s[1], db)
						http.Redirect(w, r, "/", 302)
					case "download" :
						server,_ := getServersFromDataBase(db, " WHERE id="+s[1])
						game,_ := getGamesFromDataBase(db, " WHERE nameShort='"+server[0].Type+ "'")
						gameServer("download", server[0].Localization, strconv.Itoa(game[0].GamePlatformId), server[0].StartCommands, game[0].FileToRun)
				default:
					http.Redirect(w, r, "/", 302)
				}
			}

			tpl = template.Must(template.ParseFiles("index.html"))
			tpl.Execute(w, struct{
				Name string
				Servers []GameServers
				Games []GamesStruct
			}{session.Values["name"].(string),gameServersResults,gamesResults})

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

		userResults, ifExists := getUsersFromDataBase(db, (" WHERE name='"+ strings.ToLower(r.FormValue("login"))+ "'"))
		//usersResults.PasswordHash
		if r.FormValue("login") != "" && ifExists && CheckPasswordHash(userResults[0].PasswordHash, r.FormValue("password")) {
			session.Values["authenticated"] = true
			session.Values["name"] = strings.ToLower(r.FormValue("login"))
			session.Values["id"] = strconv.Itoa(userResults[0].Id)
			session.Values["levelOfUser"] = strconv.Itoa(userResults[0].Level)
			session.Save(r, w)
	 	}

		http.Redirect(w, r, "", 302)
}
func main() {

	//id := gameServer("download", "../../servers", "90", "")
//	id := gameServer("runGame", "./servers/hlds.exe", "", "-console -game cstrike +maxplayers 20 +map de_dust2 -sv_lan 0 -port 27015")

		port := os.Getenv("PORT")
		if port == "" {
			port = "3000"
		}

		mux := http.NewServeMux()

		fs := http.FileServer(http.Dir("assets"))
		mux.Handle("/assets/", http.StripPrefix("/assets/", fs))

		mux.HandleFunc("/api/", apiHandler)
		mux.HandleFunc("/", indexHandler)
		http.ListenAndServe(":"+port, mux)
}
