/*
** Client.go ~ A simple chat client.
** Written by Luke Gorski
**/
package main

import (
    "fmt"
    "net"
    "bufio"
    "os"
)

func main() {

    var conn net.Conn
    var err error
    defaultPort := ":10500"
    reader := bufio.NewReader(os.Stdin)
    connected := false


    for {
        /* If we are not currently connected to the server,
         * prompt user for IP to connect.
         */
        if !connected {
            fmt.Println("Enter IP in order to connect, or type exit to terminate.")
            ip, _ := reader.ReadString('\n')

            //Remove the newline character from IP.
            ip = ip[0:len(ip)-1]
            if ip == "exit" {
                return
            }

            //Attempt to connect to server.
            conn, err = net.Dial("tcp", ip + defaultPort)
            if err != nil {
                fmt.Println("Error connecting, verify IP address.")
            } else {
                connected = true
                fmt.Println("Connected, type messages below. Type !quit to disconnect.")

                //Upon succesfull connection, launch getMessages
                //in a seperate goroutine to fetch messages from server
                go getMessages(conn)
            }
        }
        if connected {
            msg, _ := reader.ReadString('\n')
            if msg == ("!quit\n") {
                connected = false
            }
            //Send message to the server
            fmt.Fprintf(conn, msg + "\n")
        }
    }
    conn.Close()
}

/*
 * Reads messages from server and prints to terminal.
 */
func getMessages(connection net.Conn) {
    for {
        reader := bufio.NewReader(connection)
        msg, _ := reader.ReadString('\n')
        fmt.Print(msg)
    }

}