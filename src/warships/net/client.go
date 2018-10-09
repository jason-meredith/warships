package net

import (
    "bytes"
    "bufio"
    "fmt"
    "strings"

    "net"
    "os"
    "encoding/gob"
    "warships/game"
)

// const StopCharacter  = '\n'

type NetMessage struct {
    player *game.Player
    Command string
}

func StartClient() {

    // Retrieve the server address:port from args and attempt to connect
    service := os.Args[2]
    tcpAddr, _ := net.ResolveTCPAddr("tcp4", service)
    conn, _ := net.Dial("tcp", tcpAddr.String())

    // Reader for reading keyboard input
    reader := bufio.NewReader(os.Stdin)

    for {

        msg := NetMessage{}

        // Prompt for and read player command
        fmt.Print("Enter command > ")
        command, _ := reader.ReadString('\n')

        msg.Command = strings.TrimSpace(command)

        // Gob the command to send to server
        buf := new(bytes.Buffer)
        encoder := gob.NewEncoder(buf)
        err := encoder.Encode(msg)

        checkError(err)

        conn.Write(buf.Bytes())

        // Read server response
        buffer := make([]byte, 1024)
        resultLength, _ := conn.Read(buffer)

        fmt.Println(string(buffer[:resultLength]) + "\n")

    }

    os.Exit(0)


}