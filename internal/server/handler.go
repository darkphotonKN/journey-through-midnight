package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/darkphotonKN/journey-through-midnight/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
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

	// find player with this unique connection
	targetPlayer, _ := s.findPlayerByConnection(conn)

	fmt.Printf("Starting listener for user %v\n", targetPlayer)

	for {
		_, message, err := conn.ReadMessage()

		if err != nil {
			// --- clean up connection ---

			// Unexpected Error
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {

				fmt.Printf("Abormal error occured with player %v. Closing connection.\n", targetPlayer)

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
* Handles adding clients and creating gameMsgChans for handling connection writes
* back to the connected client.
*
* NOTE: Gorilla Websocket package only allows ONE CONCURRENT WRITER
* at a time, meaning its best to utilize *unbuffered* channels to prevent
* a single client from locking the entire server.
**/
func (s *Server) setupClientWriter(conn *websocket.Conn) {

	// in the case the channel exists
	if msgChan := s.getGameMsgChan(conn); msgChan != nil {
		// concurrently listen to all incoming messages over the channel to write game actions
		// back to the client
		go func() {
			// reading from unbuffered channel to prevent more than one write
			// a time from ANY single connection
			for msg := range msgChan {
				err := conn.WriteJSON(msg)
				if err != nil {
					// TODO: remove connection from channel and close
					s.cleanUpClient(conn)
					break
				}
			}
		}()
	}
}

/**
* Adds a player to a list of online players via their unique connection.
**/

func (s *Server) addClient(conn *websocket.Conn, playerRequest model.PlayerRequest) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.playersOnline[playerRequest.ID] = model.Player{
		ID:       playerRequest.ID,
		UserName: playerRequest.UserName,
		Conn:     conn,
	}
}

/**
* Removes a player from the list of online players via their unique connection.
**/
func (s *Server) removeClient(conn *websocket.Conn) {
	// lock and unlock to prevent race conditions
	s.mu.Lock()
	defer s.mu.Unlock()

	// find corresponding player based on their connection
	player, err := s.findPlayerByConnection(conn)

	if err != nil {
		fmt.Printf("Error when attempting to find player with connection %s\n", err)
		return
	}

	// remove from list of connections
	delete(s.playersOnline, player.ID)

	fmt.Println("Player removed from server:", player)
}
