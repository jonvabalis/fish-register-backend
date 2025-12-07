package handlers

import (
	"net/http"
	"sync"

	"fish-register-backend/internal/core"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	mu      sync.RWMutex
	clients = make(map[*websocket.Conn]bool)
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (app *FishApi) RunWebsocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	clients[conn] = true
	handleWebSocketConnection(conn)
}

func handleWebSocketConnection(conn *websocket.Conn) {
	defer func() {
		mu.Lock()
		delete(clients, conn)
		mu.Unlock()
		conn.Close()
	}()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

func SendNotification(notification core.Notification) {
	mu.Lock()
	defer mu.Unlock()

	for client := range clients {
		err := client.WriteJSON(notification)
		if err != nil {
			client.Close()
			delete(clients, client)
		}
	}
}
