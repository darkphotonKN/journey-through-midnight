package server

import (
	"fmt"

	"github.com/darkphotonKN/journey-through-midnight/internal/model"
	"github.com/google/uuid"
)

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

	// --- event actions ---\
	event_choice Action = "event_choice"

	// -- shop event --
	buy_item   Action = "buy_item"
	leave_shop Action = "leave_shop"
)

type GameEventAction struct {
	GameID      uuid.UUID `json:"game_id"`
	EventID     uuid.UUID `json:"event_id"`
	EventChoice int       `json:"event_choice"`
}

/**
* Provides conversion of the Payload based on the action type.
**/
func (gm *GameMessage) ParsePayload() error {
	switch gm.Action {

	case event_choice:
		payloadMap, ok := gm.Payload.(map[string]interface{})

		if !ok {
			fmt.Println("Could not be parsed into a map.")
			return fmt.Errorf("Could not be parsed into a map.")
		}

		// custom conversion of incoming game data
		// Extract and convert each field
		gameIDStr, ok := payloadMap["game_id"].(string)
		if !ok {
			return fmt.Errorf("game_id is not a string")
		}

		eventIDStr, ok := payloadMap["event_id"].(string)
		if !ok {
			return fmt.Errorf("event_id is not a string")
		}

		eventChoice, ok := payloadMap["event_choice"].(float64)
		if !ok {
			return fmt.Errorf("event_choice is not a number")
		}

		// Parse UUIDs
		gameID, err := uuid.Parse(gameIDStr)
		if err != nil {
			return fmt.Errorf("invalid game_id: %v", err)
		}

		eventID, err := uuid.Parse(eventIDStr)
		if err != nil {
			return fmt.Errorf("invalid event_id: %v", err)
		}

		parsedGameEventAction := GameEventAction{
			GameID:      gameID,
			EventID:     eventID,
			EventChoice: int(eventChoice),
		}

		gm.Payload = parsedGameEventAction

	case find_match:

		payloadMap, ok := gm.Payload.(map[string]interface{})
		if !ok {
			fmt.Println("Payload is not of type map[string]interface{}")
			return fmt.Errorf("Payload is not of type map[string]interface{}")
		}

		fmt.Printf("Converting to player: %+v\n", payloadMap)

		// Extract the "id" and "name" fields from the map
		idStr, idOk := payloadMap["id"].(string)
		username, nameOk := payloadMap["username"].(string)

		if !idOk && !nameOk {
			fmt.Println("Error: The required fields 'id' or 'name' are not in the expected format.")
			return fmt.Errorf("Error: The required fields 'id' or 'name' are not in the expected format.")
		}

		if !idOk {
			fmt.Println("Error: The required field 'id' is not in the expected format.")
			return fmt.Errorf("Error: The required field 'id' is not in the expected format.")
		}

		if !nameOk {
			fmt.Println("Error: The required field 'name' is not in the expected format.")
			return fmt.Errorf("Error: The required field 'name' is not in the expected format.")
		}

		fmt.Println("idStr:", idStr)
		fmt.Println("name:", username)

		idUUID, err := uuid.Parse(idStr)

		if err != nil {
			fmt.Println("Could not parse id into a UUID.")
			return fmt.Errorf("Could not parse id into a UUID.")
		}

		convertedPlayer := model.PlayerRequest{
			ID:       idUUID,
			UserName: username,
		}

		fmt.Printf("Converted player: %+v\n", convertedPlayer)

		gm.Payload = convertedPlayer

	default:
		return fmt.Errorf("No valid Action from game message.")
	}

	return nil
}
