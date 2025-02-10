package server

import (
	"github.com/darkphotonKN/journey-through-midnight/internal/model"
	"github.com/gorilla/websocket"
	uuid "github.com/jackc/pgx/pgtype/ext/gofrs-uuid"
)

type Server struct {
	// players concurrently online
	players map[uuid.UUID]model.Player

	// all current game connections
	games map[*websocket.Conn]model.GameInformation
}
