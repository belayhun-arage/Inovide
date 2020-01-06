package entity

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// Client is a middleman between the websocket connection and the hub.
// We Will Use this Struct Data in the RAM of the system For Accessing Different Web Socket Users online
// Using their Id Number
type Client struct {
	TheDistributor *Hub
	Conn           *websocket.Conn
	// Buffered channel of outbound messages.
	Send                 chan *Message
	IdentificationNumber int
}

func NewClient(theDistributor *Hub, conn *websocket.Conn, id int) *Client {

	return &Client{TheDistributor: theDistributor, Conn: conn, Send: make(chan *Message), IdentificationNumber: id}
}

/*
This Pump is Deliberately Made to Read Message From The Client Side of the System  Running infinitely
Each Online Connected Cient Will Have Will Have a single Infinitely Running Reading Pump
*/
func (c *Client) ReadPump() {

	fmt.Println("    Incomming Message  ...")

	defer func() {
		c.TheDistributor.Unregister <- c
		c.Conn.Close()
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	var mMessage *Message
	for {
		mMessage = &Message{}
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		err = json.Unmarshal(message, mMessage)
		if err != nil {
			continue
		}
		c.TheDistributor.Message <- mMessage
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			jsonMessage, err := json.Marshal(message)
			if err != nil {
				continue
			}
			w.Write(jsonMessage)
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				val, _ := json.Marshal(<-c.Send)
				w.Write(val)
			}
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
