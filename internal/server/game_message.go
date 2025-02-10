package server

/**
* Game Message
* Responsible for handling all the message structure and formatting sent
* between the server and clients.
**/
type GameMessage struct {
	Action  Action      `json:"action"`
	Payload interface{} `json:"payload"`
}

// Enum of Action Types
type Action string

const (
	find_match  Action = "find_match"
	init_match  Action = "init_match"
	match_error Action = "match_error"
	choice      Action = "choice"
)
