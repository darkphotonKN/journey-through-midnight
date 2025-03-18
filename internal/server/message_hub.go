package server

import (
	"fmt"

	"github.com/darkphotonKN/journey-through-midnight/internal/game"
	"github.com/darkphotonKN/journey-through-midnight/internal/model"
	"github.com/google/uuid"
)

/**
* Websocket Message Hub to handle all messages.
**/
func (s *Server) MessageHub() {
	fmt.Println("Starting Message Hub")

	for {

		select {

		// --- handling client sent messages passed on from the server ---
		case clientPackage := <-s.serverChan:
			fmt.Printf("Client Package received: %+v\n\n", clientPackage)

			fmt.Printf("------> Current client connections in session: %+v\n\n", s.playersOnline)

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
				err := s.matchMaker.JoinMatchMaking(&model.Player{
					ID:       player.ID,
					UserName: player.UserName,
					Conn:     clientPackage.Conn,
				})

				// responsd error back to client
				if err != nil {
					// game client-specific channel
					gameMsgChan, chanErr := s.getGameMsgChan(clientPackage.Conn)

					if chanErr != nil {
						fmt.Printf("Error when attempting to send message back to player: %s.", err)
						continue
					}

					gameMsgChan <- GameMessage{Action: "error", Payload: struct {
						Message string `json:"message"`
					}{Message: err.Error()}}

				}

			case buy_item:
				fmt.Printf("Player %+v is attempting to buy an item.\n", clientPackage)

			// TODO: Move this to its own instance or implement a way to validate game instance
			case match_error:
				fmt.Println("Match errored from client side.")
				continue
			}

		// --- handling new game initializations ---
		case newGame := <-s.matchMaker.GetNewGameChan():

			fmt.Printf("\nNew game was started, players:\n\n")
			for _, player := range newGame.Players {
				fmt.Printf("\nPlayer: %s\n", player.UserName)
			}

			// store new game on server and start initializing a unique goroutine for
			// the respective players
			err := s.addGameToServer(newGame)

			if err != nil {
				// get each game client-specific channel and send them the error
				playerGameMsgChans := make(map[uuid.UUID]chan GameMessage)

				for _, player := range newGame.Players {
					gameMsgChan, chanErr := s.getGameMsgChan(player.Conn)

					// skip non-existant channels
					if chanErr != nil {
						continue
					}

					playerGameMsgChans[player.ID] = gameMsgChan
				}

				// game already exists
				if err == game.ErrGameExists {
					// broadcast to each player the error

					for _, player := range newGame.Players {
						// get that player's channel
						gameMsgChan := playerGameMsgChans[player.ID]

						gameMsgChan <- GameMessage{Action: "error", Payload: struct {
							Message string `json:"message"`
						}{Message: err.Error()}}
					}

				}

				for _, player := range newGame.Players {
					// get that player's channel
					gameMsgChan := playerGameMsgChans[player.ID]

					gameMsgChan <- GameMessage{Action: "error", Payload: struct {
						Message string `json:"message"`
					}{Message: "Unknown error occured when attempting to start game."}}
				}

			}

			// -- start game management goroutine --

			// NOTE:
			// Player 1 WebSocket <--> Read Goroutine ---> |
			//																				 Message Hub  <--> Game Goroutine
			// Player 2 WebSocket <--> Read Goroutine ---> |
			//                                             |
			//                  Write Goroutines <----     |
			// Server owns the connection, sessions, and coordination so this lives in server
			// Game owns the logic, rules, and state of the game

			go s.manageGameLoop(newGame.ID)

		}
	}
}
