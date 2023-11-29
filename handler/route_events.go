package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/elizraa/chitchat/data"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow connections from any origin.
		},
	}
)

// Upgrade to a ws connection
// Add to active chat session
// GET /chats/{titleOrID}/ws
func webSocketHandler(w http.ResponseWriter, r *http.Request, ctxId string) (err error) {
	queries := mux.Vars(r)
	if ID, ok := queries["ID"]; ok {
		// cr, err := data.CS.Retrieve(titleOrID)
		intID, err := strconv.Atoi(ID)
		if err != nil {
			Info("ID not integer", r, err)
			return err
		}
		crDB, err := data.GetChatRoomByID(intID)
		if err != nil {
			Info("erroneous chats API request", r, err)
			return err
		}
		cr, _ := data.CS.Retrieve(crDB.Title)
		if cr == nil {
			fmt.Println("============================================")
			// Chat room doesn't exist in data.CS, create one and push it
			cr = &data.ChatRoom{
				ChatRoomDB: *crDB,
			}
			data.CS.PushCR(cr)
		}
		// Do stuff here:
		wsConn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			errorMessage(w, r, "Critical error creating WebSocket: "+err.Error())
			Danger("error creating WebSocket: ", err)
			return &data.APIError{Code: 301}
		}
		client := &data.Client{Room: cr, Conn: wsConn, Send: make(chan []byte)}
		client.Room.Broker.OpenClient <- client

		// Allow collection of memory referenced by the caller by doing all work in
		// new goroutines.
		go client.WritePump()
		go client.ReadPump()
	} else {
		return &data.APIError{Code: 101}
	}

	return
}
