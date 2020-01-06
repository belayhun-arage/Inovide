package entity

import (
	handler "github.com/Samuael/Projects/Inovide/controller"
)

type Hub struct {
	// Registered clients.
	Clients map[*Client]bool

	// Inbound messages from the clients.
	Message chan *Message

	// Register requests from the clients.
	Register chan *Client

	// Unregister requests from clients.
	Unregister chan *Client
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
var chatHandler *handler.ChatHandler

func SetChatControllHandler(chatServe *handler.ChatHandler) {
	chatHandler = chatServe
}

func NewHub() *Hub {
	return &Hub{
		Message:    make(chan *Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
	}
}
func (h *Hub) Run() {
	for {

		select {
		case client := <-h.Register:
			h.Clients[client] = true

		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
			}
		case message := <-h.Message:

			Mesage := chatHandler.SaveMesage(message)
			if Mesage != nil {
				continue
			}
			for client := range h.Clients {
				// select client.IdentificationNumber  {
				// 	case Mesage.SenderId:
				// 	default:
				// 	close(client.Send)
				// 	delete(h.Clients, client)
				// }
				if client.IdentificationNumber == message.SenderId || client.IdentificationNumber == message.RecieverId {

					client.Send <- Mesage
				}

			}
		}

	}
}
