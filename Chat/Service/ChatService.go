package ChatService 


import (
//	"bytes"
	// "flag"
	// "log"
	 "net/http"
	// "time"

	"crypto/rand"
	"encoding/base64"
	"github.com/gorilla/websocket"
	"io"
	entity "github.com/Samuael/Projects/Inovide/models"
	"log"
)


var keyGUID = []byte("258EAFA5-E914-47DA-95CA-C5AB0DC85B11")

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}


type ChatService struct {
	userrepo *repository.UserRepo
}


func NewChatService(){}





func generateKey() (string, error) {
	p := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, p); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(p), nil
}


func (chatService *ChatService) serveWs(hub *entity.Hub , w http.ResponseWriter, r *http.Request) {

	wsKey, _ := generateKey()
	r.Header.Add("Connection", "Upgrade")
	r.Header.Add("Upgrade", "websocket")
	r.Header.Add("Sec-WebSocket-Version", "13")
	r.Header.Add("Sec-WebSocket-Key", wsKey)
	log.Printf("ws key '%v' ----  ", wsKey)

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
		return
	}
	client := &entity.Client{TheDistributor : hub, Conn: conn, Send: make(chan []byte, 256)}
	client.TheDistributor.Register <- client
	go client.writePump()
	go client.readPump()
}