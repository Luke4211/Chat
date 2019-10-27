package main

import (
    "fmt"
    "net"
    "bufio"
    "sync"
    "strconv"
)



func main() {
    numClients := 0
    listener, _ := net.Listen("tcp", ":10500")
    
    var clientChans = struct {
        sync.RWMutex
        m map[int]chan string
    } {m: make(map[int]chan string)}

    for {
        connection, _ := listener.Accept()
        ch := make(chan string)

        clientChans.Lock()
        clientChans.m[numClients] = ch
        clientChans.Unlock()

        go handleClient(connection, clientChans, numClients)
        go getMessages(connection, clientChans, numClients)
        numClients++
    }
}
func handleClient(connection net.Conn, clientChans struct {sync.RWMutex; m map[int]chan string}, clientNum int) {

    fmt.Println("Connected to client #" + strconv.Itoa(clientNum) )
    for msg := range clientChans.m[clientNum] {
        if msg == "!quit" {
            clientChans.Lock()
            close(clientChans.m[clientNum])
            delete(clientChans.m, clientNum)
            clientChans.Unlock()
            connection.Close()
            return
        } else {
            clientChans.Lock()
            connection.Write([]byte(msg + "\n"))
            clientChans.Unlock()
        }
    }

}

func getMessages(connection net.Conn, clientChans struct {sync.RWMutex; m map[int]chan string}, clientNum int) {
    running := true
    for running {
        message, _ := bufio.NewReader(connection).ReadString('\n')
        if message == "!quit\n" {
            running = false
            clientChans.Lock()
            clientChans.m[clientNum]<- "!quit"
            clientChans.Unlock()
        } else {
            fmt.Print(message)
            go broadcast(message, clientChans, clientNum)
        }

    }
}

func broadcast(msg string, clientChans struct { sync.RWMutex; m map[int]chan string} , clientNum int) {
    clientChans.Lock()
    for key, client := range clientChans.m {
        if key != clientNum {
            client<- msg
        }
    }
    clientChans.Unlock()
}
