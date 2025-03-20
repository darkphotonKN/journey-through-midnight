package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

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

	for {
		fmt.Println("Listening for user messages...")
		_, message, err := conn.ReadMessage()

		fmt.Printf("\nMessage received from connected user: %s\n\n", string(message))

		// find player with this unique connection
		targetPlayer, playerErr := s.findPlayerByConnection(conn)

		if playerErr == nil {
			fmt.Printf("Message from registered player: %v\n", targetPlayer)
		} else {
			fmt.Println("Message from unregistered connection (normal for first messages), error:", playerErr)
		}

		fmt.Printf("Starting listener for user %v\n", targetPlayer)

		// --- clean up connection ---
		if err != nil {
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

		// --- Client Connection Handling ---
		// Decodes Incoming client message and serves their unique connection its own goroutine

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

		// Send message to MessageHub via an *unbuffered channel* for handling based on the type field.
		s.serverChan <- clientPackage
	}
}

/**
* Creates the unique game message channel for a specific connection for writing back
* from server to client.
**/
func (s *Server) createGameMsgChan(conn *websocket.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.gameMsgChan[conn] = make(chan GameMessage)
}

/**
* Gets the unique game message channel for a specific connection for writing back
* from server to client, validating that it exists.
**/
func (s *Server) getGameMsgChan(conn *websocket.Conn) (chan GameMessage, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	channel, exists := s.gameMsgChan[conn]

	if !exists {
		return nil, fmt.Errorf("Game message channel for this connection does not exist.")
	}

	return channel, nil
}

/**
* Handles adding clients and creating gameMsgChans for handling connection writes
* back to the connected client.
*
* NOTE: Gorilla Websocket package only allows ONE CONCURRENT WRITER
* at a time, meaning its best to utilize *unbuffered* channels to prevent
* a single client from locking the entire server, and prevent race conditions
* where multiple writes to the same connection.
**/
func (s *Server) setupClientWriter(conn *websocket.Conn) {
	// sets up this connection's personal game message channel
	s.createGameMsgChan(conn)

	// in the case the channel exists
	msgChan, err := s.getGameMsgChan(conn)

	if err != nil {
		fmt.Println(err)
		s.cleanUpClient(conn)
		return
	}

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

/**
* Cleans up the client connected to the online server from all relevant data structures.
**/
func (s *Server) cleanUpClient(conn *websocket.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// removes client from players online
	s.removeClient(conn)

	// close their channel
	channel, _ := s.getGameMsgChan(conn)
	close(channel)

	// removes their personal gameMsgChan
	delete(s.gameMsgChan, conn)
}
