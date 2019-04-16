package main

import (
	"log"
	"time"
	"strings"
	"net/http"
	"github.com/gorilla/websocket"
)

/* 
server:
once connected, client should send wechat id/sex/nickname/desc/avatar(base64) to server,
server store each person's info in global map, only online people, so no pressure of mem.
there are 3 table in mem:
{id: [sex, nickname, desc, avatar]}
{id: *conn}
{id: partner_id}

client: 
login -> click start -> websocket conected -> send user data to register self
-> waiting partner -> chatting
*/
var userglobalmap = make(map[string][]string)
var connglobalmap = make(map[string]*websocket.Conn)
var partnerglobalmap = make(map[string]string)

func main() {	
	http.HandleFunc("/", handler)
	log.Println("ws server start at :35000...")
	err := http.ListenAndServe(":35000", nil)
	checkError(err)	
	log.Println("ws server stopped.")
}

var upgrader = websocket.Upgrader{}
func handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	checkError(err)	
	defer conn.Close()
	//now ws connected.

	// handle first message: 
	// must be user's id/sex/nickname/desc/avatar(send by js on_connected)
	// client js must do strip, then join them with '\n'
	_, message, err := conn.ReadMessage()
	checkError(err)
	tmp := strings.Split(string(message), "\n")
	// register self into global map
	id := tmp[0]
	userglobalmap[id] = tmp[1:]
	connglobalmap[id] = conn
	partnerglobalmap[id] = ""
	defer func () {
		// empty user info when offline
		connglobalmap[id] = nil
		partnerglobalmap[partnerglobalmap[id]] = ""
		partnerglobalmap[id] = ""
	}()

	// finding partner until found.
	for {
		time.Sleep(3*time.Second)
		// i have no partner, finding
		if partnerglobalmap[id] == "" {
			for k, v := range partnerglobalmap {
				// not myself and someone online
				if id != k && connglobalmap[k] != nil {
					// found someone no partner too or his partner is me
					if  v == "" || v == id {
						// i am already found by A when i found B
						if partnerglobalmap[id] != "" {
							break
						}
						partnerglobalmap[id] = k
						partnerglobalmap[k] = id
						break
					}
				}		
			}
		} else {
			break
		}
	}	

	// start chatting with partner.
	// A <-> server <-> B
	for {
		// server wait here until client send message
		mt, message, err := conn.ReadMessage()
		checkError(err)
		// send message to partner
		partnerconn := connglobalmap[partnerglobalmap[id]]
		// if partner offline
		if partnerconn == nil {
			partnerconn.WriteMessage(mt, []byte("[offline] sorry, your partner is offline."))
		}
		err = partnerconn.WriteMessage(mt, message)
		checkError(err)
	}
}

func checkError(err error) {
	if err != nil {
		log.Printf("Error => %s", err.Error())
	}
}




