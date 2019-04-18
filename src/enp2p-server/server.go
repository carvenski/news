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
    // set log format
    log.SetFlags(log.Lshortfile | log.LstdFlags)    
    http.HandleFunc("/", handler)
    log.Println("ws server start at :35000...")
    err := http.ListenAndServe(":35000", nil)
    if err != nil {
        log.Printf("Error => %s", err.Error())
    }    
    log.Println("ws server stopped.")
}

var upgrader = websocket.Upgrader{}
func handler(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Printf("Error => %s", err.Error())
    }    
    defer func ()  {
        // close conn
        conn.Close()
        // try-catch any exception, do nothing just return
        if err := recover(); err != nil {
            log.Println("ERROR => ", err)
        }
    }()
    //now ws connected.
    log.Println("1 ws client connected from: ", conn.RemoteAddr())

    // handle first message: 
    // must be user's id/sex/nickname/desc/avatar(send by js on_connected)
    // client js must do strip, then join them with '\n'
    _, message, err := conn.ReadMessage()
    if err != nil {
        log.Printf("Error => %s", err.Error())
    }    
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
    FINDING_PARTNER:
    for {
        log.Printf( "client %s is finding partner...", id)
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
        log.Println("chatting...")
        // server wait here until client send message
        mt, message, err := conn.ReadMessage()
        if err != nil {
            // this client offline
            log.Printf("=> client %s offline: %s", id, err.Error())
            break
        }        
        // if partner offline
        if partnerglobalmap[id] == "" {
            conn.WriteMessage(mt, []byte("[WARNING] sorry, your partner is offline."))
            goto FINDING_PARTNER
        }
        // send message to partner
        partnerconn := connglobalmap[partnerglobalmap[id]]
        err = partnerconn.WriteMessage(mt, message)
        if err != nil {
            log.Printf("Error => %s", err.Error())
        }    
    }
}



