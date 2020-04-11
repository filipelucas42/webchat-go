package main

import (
	"net/http"

	"golang.org/x/net/websocket"
)

type Connection struct {
	conn *websocket.Conn
	name string
}

type Message struct {
	name    string
	message []byte
}

var NewUser chan Connection
var NewMessage chan Message
var DeleteUser chan string

func ProcessSockets(newUser <-chan Connection, newMessage <-chan Message, deleteUser <-chan string) {
	for {
		select {
		case createUser := <-newUser:
			connections[createUser.name] = createUser.conn

			for _, content := range connections {
				aux := []byte(createUser.name + " joined chat")
				content.Write(aux)
			}
		case message := <-newMessage:
			for _, content := range connections {
				aux := []byte("From " + message.name + ": ")
				content.Write(append(aux[:], message.message[:]...))
			}
		case delUser := <-deleteUser:

			connections[delUser].Close()
			for _, content := range connections {
				aux := []byte(delUser + " disconnected from chat")
				content.Write(aux)
			}
			delete(connections, delUser)
		}
	}
}

var connections map[string]*websocket.Conn

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./public/index.html")
}

func websocketHandler(ws *websocket.Conn) {
	var message []byte
	var name []byte
	websocket.Message.Receive(ws, &name)
	connection := Connection{
		name: string(name),
		conn: ws,
	}
	if _, ok := connections[string(name)]; ok {
		ws.Write([]byte("Username already in use"))
		ws.Close()
		return
	}
	NewUser <- connection
	newMessage := Message{
		name: string(name),
	}
	for {
		err := websocket.Message.Receive(ws, &message)
		if err != nil {
			break
		}
		//fmt.Println("New message received from " + string(name))
		newMessage.message = message
		NewMessage <- newMessage
	}
	//delete(connections, string(name))
	DeleteUser <- string(name)
}

func main() {
	NewUser = make(chan Connection)
	NewMessage = make(chan Message)
	DeleteUser = make(chan string)
	connections = make(map[string]*websocket.Conn)
	go ProcessSockets(NewUser, NewMessage, DeleteUser)
	fs := http.FileServer(http.Dir(("./public")))
	http.Handle("/public/", http.StripPrefix("/public/", fs))
	http.HandleFunc("/", indexHandler)
	http.Handle("/socket", websocket.Handler(websocketHandler))
	http.ListenAndServe(":8082", nil)
}
