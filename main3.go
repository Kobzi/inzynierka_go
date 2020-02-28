package main

import (
	"fmt"
	"os/exec"
	"log"
	"bytes"
//	"bufio"
)

func main() {
	//v, _ := mem.VirtualMemory()
//	fmt.Printf("Total: %v, Free:%v, UsedPercent:%f%%\n", v.Total, v.Free, v.UsedPercent)

    // convert to JSON. String() is also implemented
    fmt.Println("test")
		cmd := exec.Command("src/SteamCMD/steamcmd.exe", "-login anonymous -force_install_dir /server -app_update 346680 validate")
		var stdout, stderr bytes.Buffer
    cmd.Stdout = &stdout
    cmd.Stderr = &stderr
    err := cmd.Run()
    if err != nil {
        log.Fatalf("cmd.Run() failed with %s\n", err)
    }
    outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
    fmt.Printf("out:\n%s\nerr:\n%s\n", outStr, errStr)

}
