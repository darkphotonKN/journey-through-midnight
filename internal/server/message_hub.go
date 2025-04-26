package server

import (
	"fmt"
	"strconv"

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

		// --- handling messages sent from server over serverChan ---
		case clientPackage := <-s.serverChan:
			fmt.Printf("Client Package received: %+v\n\n", clientPackage)

			fmt.Printf("------> Current client connections in session: %+v\n\n", s.playersOnline)

			// TODO: Remove after debugging
			fmt.Printf("\n------> Current games information session: %+v\n\n", s.games)

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

			// -- serving messages based on action type --
			switch clientPackage.GameMessage.Action {

			case find_match:
				// assert Payload type specific to gameMessage.Action == "find_match", which is PlayerRequest
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

				// response error back to client
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

			// individual game-specific actions
			case event_choice:
				gameEventAction := clientPackage.GameMessage.Payload.(GameEventAction)

				// find corresponding game
				currentGame, exists := s.games[gameEventAction.GameID]

				if !exists {
					// get the player's write-back game message channel
					playerMsgChan, err := s.getGameMsgChan(clientPackage.Conn)

					playerMsgChan <- GameMessage{Action: "error", Payload: struct {
						Message string `json:"message"`
					}{Message: err.Error()}}
				}

				// casting event choice to a string and sending it to the original game
				currentGame.MsgCh <- strconv.Itoa(gameEventAction.EventChoice)

			case buy_item:
				fmt.Printf("Player %+v is attempting to buy an item.\n", clientPackage)

				// send instruction over to game loop

			case match_error:
				fmt.Println("Match errored from client side.")
				continue

			}

		// --- handling new game initializations ---
		case newGame := <-s.matchMaker.GetNewGameChan():

			fmt.Printf("\nNew game was started. Game Info\n\n%+v\n\nplayers:\n", newGame)
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
			//															Server.Message Hub Goroutine  <-->  Game.ManageGameLoop Goroutine
			// Player 2 WebSocket <--> Read Goroutine ---> |
			//                                             |
			//                  Write Goroutines <----     |

			go newGame.ManageGameLoop()

		// TOOD: communicate events to an existing game loop
		// TODO: Need to move this to the serverChan send above
		case message := <-s.serverChan:

			err := message.GameMessage.ParsePayload()

			if err != nil {
				fmt.Printf("Error occured when attempting to parse payload: %s\n", err)
				message.Conn.WriteJSON(fmt.Sprintf("Error attempting to parse payload: %s", err))
				continue
			}

			// TODO: game communication
			// in this situation the action payload will contain the game id
			// message.GameMessage.Payload
			//
			// // find game message channel
			// ch := s.games[3].MsgCh
			//
			// ch <- eventChoice.GameMessage.Payload.(string)

		}
	}
}
