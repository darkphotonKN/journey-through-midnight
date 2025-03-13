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

			// deduce player from package
			player, err := s.findPlayerByConnection(clientPackage.Conn)

			if err != nil {
				fmt.Println("Player hasn't joined in any game yet.")
			} else {
				fmt.Printf("Player %+v sending action.\n", player)
			}

			// NOTE: parses payload to a specific type based on the action type
			// e.g. when payload is "find_match" the payload is converted from interface{} -> Player
			err = clientPackage.GameMessage.ParsePayload()

			if err != nil {
				fmt.Printf("Error occured when attempting to parse payload: %s\n", err)
				clientPackage.Conn.WriteJSON(fmt.Sprintf("Error attempting to parse payload: %s", err))
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

				// add player to join the match making queue
				s.matchMaker.JoinMatchMaking(&model.Player{
					ID:       player.ID,
					UserName: player.UserName,
					Conn:     clientPackage.Conn,
				})

			case buy_item:
				fmt.Printf("Player %+v is attempting to buy an item.\n", clientPackage)

			// TODO: Move this to its own instance or implement a way to validate game instance
			case match_error:
				fmt.Println("Match errored from client side.")
				break
			}

		case newGame := <-s.matchMaker.GetNewGameChan():
			fmt.Printf("\nNew game was started, players:\n\n")
			for _, player := range newGame.Players {
				fmt.Printf("\nPlayer: %s:\n", player.UserName)
			}

		}
	}
}
