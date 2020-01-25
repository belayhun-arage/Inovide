package entity

import (
	"encoding/json"
	"fmt"
	"time"

	// "github.com/gorilla/websocket"
	"github.com/jinzhu/gorm"
	"golang.org/x/net/websocket" /// the Web SOcket Inmplementation usign the Built in Web Socket LIbrary of golang
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
	// c.Conn.SetReadLimit(maxMessageSize)
	// c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	// c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	var mMessage *Message
	for {
		mMessage = &Message{}
		// _, message, err := c.Conn.ReadMessage()

		erro := websocket.Message.Receive(c.Conn, &mMessage)

		if erro != nil {
			if erro.Error() == websocket.ErrBadClosingStatus.Error() || erro.Error() == websocket.ErrBadClosingStatus.ErrorString || erro.Error() == websocket.ErrBadFrame.Error() {
				break
			}
		}
		// message = bytes.Replace(message, newline, space, -1)
		// err = json.Unmarshal(message, mMessage)
		// if err != nil {
		// continue
		// }

		year, month, day := gorm.NowFunc().Date()
		nowIs := fmt.Sprintf("%d/%d/%d", day, month, year)
		mMessage.Dateofcreation = nowIs
		c.TheDistributor.Messages <- mMessage
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
			fmt.Println(message.Messagedata)
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))

			if !ok {
				// 	c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			// w, err := c.Conn.NextWriter(websocket.TextMessage)
			// if err != nil {
			// 	return
			// }
			jsonMessage, err := json.Marshal(message)
			if err != nil {
				continue
			}

			if err := websocket.Message.Send(c.Conn, jsonMessage); err != nil {

				fmt.Println("Can't Send the Mesage ")
				break
			}

			// w.Write(jsonMessage)
			// n := len(c.Send)
			// for i := 0; i < n; i++ {
			// 	w.Write(newline)
			// 	val, _ := json.Marshal(<-c.Send)
			// 	w.Write(val)
			// }
			// if err := w.Close(); err != nil {
			// 	return
			// }
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			// if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
			// 	return
			// // }
			// if err := websocket.Message.Send(c.Conn  , )
			fmt.Println("We Will Be Ticking as the time goes By   For ticking ")

		}
	}
}
