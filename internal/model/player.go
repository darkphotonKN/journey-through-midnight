package model

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Player struct {
	ID       uuid.UUID
	UserName string
	Conn     *websocket.Conn
}
