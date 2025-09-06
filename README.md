# Journey Through Midnight - Game Server

A real-time websocket game server for an RPG autobattler, written in Go. Players connect, join a matchmaking queue, and battle through events as they try to survive the night.

## What's Working

### WebSocket Server
The server runs on Gin and uses Gorilla WebSocket for real-time connections. Each player gets their own goroutine to handle their messages, so everything stays responsive even with multiple players connected. When you connect at `/ws`, you're immediately ready to start communicating with the server.

### Matchmaking Queue
Got a pretty straightforward matchmaking system that pairs players up for games. Here's what happens:
- Players send a `find_match` message with their ID and username
- They get added to the queue and receive confirmation they've joined
- Every 15 seconds, the matchmaker checks if there's at least 2 players waiting
- When it finds a match, both players get notified and a new game instance spins up
- Players are removed from the queue once matched

The whole thing runs on its own goroutine so it doesn't block anything else.

### Message Hub
This is basically the brain of the server - all messages flow through here. It handles:
- Player actions coming in from websocket connections  
- Game initialization when the matchmaker creates a match
- Routing messages to the right game instances
- Broadcasting responses back to players

Each connection gets its own channel for writing messages back, which prevents any race conditions when multiple things try to talk to the same player.

### Game Structure
Games are set up with a full RPG system in mind:
- Each game has a unique ID and tracks all its players
- Players have heroes with classes like Fighter, Wizard, Rogue, Priest, Duelist, and Templar
- Full attribute system (Strength, Intelligence, Wisdom, Agility, Vitality, Faith, Charisma)
- Inventory and gold tracking
- Day/Night phase system with events

### Message Protocol
Communication happens through JSON messages with an action and payload structure. Current actions include:
- `find_match` - Join the matchmaking queue
- `init_match` - Server confirms a match was found
- `match_error` - Something went wrong with matchmaking
- `event_choice` - Player makes a choice during a game event
- `buy_item` / `leave_shop` - Shop interactions

### Player Management
The server keeps track of:
- All players currently online
- Mapping between websocket connections and player IDs
- Active game instances
- Automatic cleanup when players disconnect (removes them from games, closes empty games, cleans up channels)

## Tech Stack
- **Go 1.23.3**
- **Gin** for HTTP routing
- **Gorilla WebSocket** for real-time connections
- **PostgreSQL** with sqlx (database connection is initialized but not actively used yet)
- **UUID** for unique identifiers

## Running the Server
Make sure you have a `.env` file with your `PORT` defined, then just run:
```bash
go run cmd/main.go
```

The server will start up, initialize the database connection, and begin listening for websocket connections. The matchmaker starts automatically and begins checking for matches every 15 seconds.

## Project Structure
```
├── cmd/
│   └── main.go              # Entry point
├── internal/
│   ├── config/              # Database and routing setup
│   ├── game/                # Game logic and mechanics
│   ├── matchmaking/         # Queue and matching system
│   ├── model/               # Player data structures
│   └── server/              # WebSocket handling and message routing
```

The architecture uses channels extensively for communication between goroutines, keeping everything concurrent but safe. Each component is designed to work independently - the matchmaker doesn't know about websockets, the game engine doesn't know about the server implementation, and so on.