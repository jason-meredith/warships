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

        // Decode incoming gob'ed Command struct
        incomingCmd := &Command{}
        decoder.Decode(incomingCmd)

        // If the Command is not null process the Command and write send response to client
        if(incomingCmd.CmdString != "") {
            returnValue := processCommand(*incomingCmd)

            conn.Write([]byte(returnValue))
        }

    }

    conn.Close()
}

func processCommand(incoming_cmd Command) string {
    fmt.Printf("Command: %v\n", incoming_cmd.CmdString)

    return "Miss! Try again"
}
