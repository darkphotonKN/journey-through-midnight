package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

/**
* Handles a player searching for a match (incoming connection), upgrades them to websocket connections,
* and passes them off to individual goroutines to be handled concurrently.
**/
func (s *Server) HandleMatchConn(c *gin.Context) {
	conn, err := s.upgrader.Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		fmt.Println("Error establishing websocket connection.")
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to upgrade connection"})
		return
	}

	// handle each connected client's messages concurrently
	go s.ServeConnectedPlayer(conn)
}

/**
* Serves each individual connected player.
**/
func (s *Server) ServeConnectedPlayer(conn *websocket.Conn) {
	// removes client and closes connection
	defer func() {
		fmt.Println("Connection closed due to end of function.")
		s.removeClient(conn)
	}()

	fmt.Printf("Starting listener for user %v\n", s.playersOnline[conn])
	for {
		_, message, err := conn.ReadMessage()

		if err != nil {
			// --- clean up connection ---

			// Unexpected Error
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Printf("Abormal error occured with player %v. Closing connection.\n", s.playersOnline[conn])
				break
			}

			// Close Error
			if websocket.IsCloseError(err, websocket.CloseGoingAway) {
				fmt.Printf("Error on close, close going away, error: %s\n", err)
				break
			}

			// General Error
			fmt.Printf("General error occured during connection: %s\n", err)
			break
		}

		// decode message to pre-defined json structure "GameMessage"
		var decodedMsg GameMessage

		err = json.Unmarshal(message, &decodedMsg)

		if err != nil {
			fmt.Println("Error when decoding payload.")
			conn.WriteJSON(GameMessage{Action: "Error", Payload: "Your message to server was the incorrect format and could not be decoded as JSON."})
			continue
		}

		// handle concurrent writes back to clients
		s.setupClientWriter(conn)

		clientPackage := ClientPackage{GameMessage: decodedMsg, Conn: conn}

		// Send message to MessageHub via an *unbuffered channel* for handling based on type.
		s.serverChan <- clientPackage
	}
}

/**
* Removes a player from the list of client connections.
**/
func (s *Server) removeClient(conn *websocket.Conn) {
	// lock and unlock to prevent race conditions
	s.mu.Lock()
	defer s.mu.Unlock()

	if player, ok := s.playersOnline[conn]; ok {
		conn.Close()

		fmt.Println("Player removed from server:", player)

		// remove from list of connections
		delete(s.playersOnline, conn)
	}
}
