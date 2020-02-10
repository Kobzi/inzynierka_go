package main

import (
    "fmt"
    "log"
    //"net/http"
    //"github.com/gorilla/websocket"
    "os"
    "os/exec"
)


func gameServer(whatToDo string, path string, gameId string, gameParametrs string) int {
  processid:=0

  switch (whatToDo) {
  case "download" :
    cmd := exec.Command("./src/SteamCMD/steamcmd.exe", "+login anonymous", "+force_install_dir "+path, "+app_update "+gameId)
    cmd.Stdout = os.Stdout
    err := cmd.Start()
    if err != nil {
       log.Fatal(err)
    }
    //log.Printf("Just ran subprocess %d, exiting\n", cmd.Process.Pid)
    processid = cmd.Process.Pid

  case "runGame" :
    cmd := exec.Command(path, gameParametrs)
    cmd.Stdout = os.Stdout
    err := cmd.Start()
    cmd.Wait()
    if err != nil {
       log.Fatal(err)
    }

    //log.Printf("Just ran subprocess %d, exiting\n", cmd.Process.Pid)
    processid = cmd.Process.Pid
    }

return processid

}
func main() {
    fmt.Println("Hello World")
    /*output, err := exec.Command("./src/SteamCMD/steamcmd.exe").Output()
    if err!=nil {
        fmt.Println(err.Error())
    }
    fmt.Println(string(output))
    log.Printf("Just ran subprocess %d, exiting\n", output.Process.Pid)*/

    //id := gameServer("download", "../../servers", "90", "")
    id := gameServer("runGame", "start '' 'servers\hlds.exe'", "", "-console -game cstrike +maxplayers 20 +map de_dust2 -sv_lan 0 -port 27015")

    fmt.Println(id)
}
