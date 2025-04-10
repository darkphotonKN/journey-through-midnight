package game

import (
	"fmt"
	"time"
)

/**
* Engine page are the core inner-workings of the game, holding together
* the logic, events, turns, etc.
*
* The primary game cycle focuses on 4 "phases".
* Phase 1 Turn start - handles what each player gets at the beginning of each turn, and the initial event they are
* are thrown into. This is sent automatically at the start of the game and at the start of each round.
* Phase 2 Daytime Events - Each player interfaces with an event per "time interval" during the day, until
* they reach the midnight phase.
* Phase 3 Midnight Phase - Midnight is focused on battling and survival. Player's wait for each other before entering
* and exiting this phase.
* Phase 4 End Turn - The round ends and calculations are made.
**/

/**
* Handles the logic for starting the game round.
**/
func (g *Game) ProcessStartRound() {
	// check and process any player bonuses
}

/**
* Handles the logic for ending the game round.
**/
func (g *Game) ProcessEndRound() {
	// increment round
	g.Round++
}

/**
* Processes a player initiated event.
**/
func (g *Game) ProcessPlayerEvent(event GameEvent, playerState PlayerState) *PlayerState {
	// copy over player state, prevent any leaks or having to lock
	psCopy := playerState

	// pass in reference and have it mutated
	g.eventHandler.initiateEvent(event.Type, &psCopy)

	// update player in game instance after event
	g.Players[psCopy.ID] = &psCopy

	return &psCopy
}

/**
* Midnight Phase
**/

func (g *Game) ProcessPlayerMidnightEvent(game *Game, event GameEvent) {
}

/**
* High level method that wraps all internal game loop components together.
*
* Manages each unique game and it's coordinations.
**/
func (g *Game) ManageGameLoop() {
	serverTick := time.Second // one second per game "tick"
	ticker := time.NewTicker(serverTick)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			g.mu.Lock()

			fmt.Printf("\ngameId: %v\n\n", g.ID)
			defer g.mu.Unlock()

			fmt.Printf("game %s currently on round: %d.\n", g.ID, g.Round)

			// aggregate information for game loop
			playersAtMidnight := 0
			playersEndOfRound := 0

			for _, player := range g.Players {
				if player.Time.Hour == 24 {
					playersAtMidnight++
				}
				if player.Time.Hour == 5 {
					playersEndOfRound++
				}
			}

			// process game transition to midnight
			if playersAtMidnight == len(g.Players) {

			}

			// stop this goroutine, game ended or errored
		case <-g.CloseGameCh:
			fmt.Println("Game has been stopped or already stopped prior, exiting ManageGameLoop goroutine.")
			return
		}
	}
}
