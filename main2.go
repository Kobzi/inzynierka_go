package main

import (
    "fmt"
    "log"
    //"net/http"
    //"github.com/gorilla/websocket"
    "os"
    "os/exec"
    "time"
)


func gameServer(whatToDo string, path string, gameId string, gameParametrs string) *exec.Cmd {
  var test *exec.Cmd

  switch (whatToDo) {
  case "download" :
    cmd := exec.Command("./src/SteamCMD/steamcmd.exe", "+login anonymous", "+force_install_dir "+path, "+app_update "+gameId)
    cmd.Stdout = os.Stdout
    err := cmd.Start()
    if err != nil {
       log.Fatal(err)
    }
    //log.Printf("Just ran subprocess %d, exiting\n", cmd.Process.Pid)
    //processid = cmd.Process.Pid

  case "runGame" :
    //cmd := exec.Command(path, gameParametrs)
    cmd := exec.Command("cmd","/c", "cd /d "+path+" && start " +gameId+ " " +gameParametrs)
    //cmd /c "cd /d ./servers && start hlds.exe -console -game cstrike +maxplayers 20 +map de_dust2 -sv_lan 0 -port 27015"
    //cmd := exec.Command("C:/Users/aizda/Go/servers/hlds.exe","-console -game cstrike +maxplayers 20 +map de_dust2 -sv_lan 0 -port 27015")
    //cmd.Dir = "C:/Users/aizda/Go/servers"
    cmd.Stdout = os.Stdout
    err := cmd.Start()
    cmd.Wait()
    if err != nil {
       log.Fatal(err)
    }
    test = cmd

    log.Printf("Just ran subprocess %d, exiting\n", cmd.Process.Pid)
    //processid = cmd.Process.Pid
    }
return test

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
    id := gameServer("runGame", "./servers", "hlds.exe", "-console -game cstrike +maxplayers 20 +map de_dust2 -sv_lan 0 -port 27015")

    fmt.Println(id)
    duration := time.Duration(10)*time.Second // Pause for 10 seconds
    time.Sleep(duration)
    fmt.Println(id.Process.Pid)
    id.Process.Kill()
}
