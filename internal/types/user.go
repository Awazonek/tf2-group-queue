package types

import "time"

type User struct {
	ID              string         `json:"id"`
	LastSeen        time.Time      `json:"timestamp"`         // Last seen, kicking them out after X minutes
	SessionCount    map[string]int `json:"session_data"`      // history of user connecting to server
	ConnectedServer string         `json:"server"`            // Server that user is connected to
	ServerToConnect string         `json:"connecting_server"` // Server that the user should connect to on next ping if they are not in a server
}
