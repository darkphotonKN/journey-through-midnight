package server

import (
	"fmt"

	"github.com/darkphotonKN/journey-through-midnight/internal/model"
)

/**
* Websocket Message Hub to handle all messages.
**/
func (s *Server) MessageHub() {
	fmt.Println("Starting Message Hub")

	for {
		fmt.Printf("Current client connections in session: %+v\n\n", s.playersOnline)

		select {
		case clientPackage := <-s.serverChan:
			fmt.Printf("Client Package received: %+v\n\n", clientPackage)

			// NOTE: parses payload to a specific type based on the action type
			// e.g. when payload is "find_match" the payload is converted from interface{} -> Player
			err := clientPackage.GameMessage.ParsePayload()

			if err != nil {
				fmt.Printf("Error occured when attempting to parse payload: %s\n", err)
				clientPackage.Conn.WriteJSON("Error attempting to parse payload.")
				continue
			}

			switch clientPackage.GameMessage.Action {
			case find_match:
				fmt.Println("Inside 'find match' case, payload:", clientPackage.GameMessage.Payload)

				// assert Payload type specific to gameMessage.Action == "find_match", which is Player
				player, ok := clientPackage.GameMessage.Payload.(model.PlayerRequest)

				if !ok {
					fmt.Printf("Error attempting to assert player from payload.\n")
					clientPackage.Conn.WriteJSON("Error attempting to assert player from payload.")
					continue
				}

				// add client to global struct connections pool
				s.addClient(clientPackage.Conn, player)

				// TODO: initiate matchmaking for the player concurrently
			}
		}
	}
}
