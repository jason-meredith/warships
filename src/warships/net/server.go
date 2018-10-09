package net

import (
    //"bufio"
    "encoding/gob"
    "net"
    "fmt"
    "os"
)


func StartServer() {

    // Get port from args and listen
    service := ":" + os.Args[2]
    tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
    listener, err := net.Listen("tcp", tcpAddr.String())

    checkError(err)

    fmt.Printf("Listening on port %v\n", os.Args[2])

    // For every incoming request create new goroutine
    for {
        conn, err := listener.Accept()
        if err != nil {
            continue
        }

        go handleClient(conn)

    }
}

func handleClient(conn net.Conn) {

    // Get the address of the connecting player
    connectionAddress := conn.RemoteAddr().String();
    fmt.Printf("Incoming connection: %v\n", connectionAddress)

    decoder := gob.NewDecoder(conn)

    for{

        // Decode incoming gob'ed NetMessage struct
        incomingMsg := &NetMessage{}
        decoder.Decode(incomingMsg)

        // If the NetMessage is not null process the NetMessage and write send response to client
        if(incomingMsg.Command != "") {
            returnValue := processCommand(*incomingMsg)

            conn.Write([]byte(returnValue))
        }

    }

    conn.Close()
}

func processCommand(incomingMsg NetMessage) string {
    fmt.Printf("NetMessage: %v\n", incomingMsg.Command)

    return "Miss! Try again"
}
