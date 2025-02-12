package model

import "github.com/gorilla/websocket"

type GameMsgChan = chan string

/**
* Holds all the infromation for a specific game's meta data.
**/
type GameInformation struct {
	Conn    *websocket.Conn // unique connection to the game
	MsgChan *GameMsgChan    // message channel to communicate with game
}
