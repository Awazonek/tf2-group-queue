package types

import (
	"fmt"
	"time"

	"github.com/awazonek/tf2-group-queue/internal/util"
)

type Group struct {
	ID         string           `json:"id"`
	Parameters ServerParameters `json:"server_parameters"`
	Users      map[string]User  `json:"users"`
	Searching  bool             `json:"searching"`
}

// returns true if the server matches the groups requirements
func (g *Group) MatchesServer(server Tf2Server) bool {
	// Basic checks
	if g.Parameters.Region != server.Region {
		util.Log("\nServer %s is too far away, wanted region %s and we got %s", server.Name, g.Parameters.Region, server.Region)
		return false
	}

	if g.Parameters.MinPlayers > server.Players {
		util.Log("\nServer %s has too few players. Wanted at least %d but got %d", server.Name, g.Parameters.MinPlayers, server.Players)
		return false

	}

	// Check player count
	if server.MaxPlayers-server.Players < g.Parameters.MinSpaceAvailable {
		util.Log("\nServer %s has too many players wanted %d slots open and got %d", server.Name, g.Parameters.MinSpaceAvailable, (server.MaxPlayers - server.Players))
		return false
	}

	// If server has correct game mode
	for _, mode := range g.Parameters.GameModes {
		if mode == server.Mode {
			return true
		}
	}
	util.Log("\nServer %s does not have the correct game mode, it has %d", server.Name, server.Mode)
	return false
}

func (g *Group) ConnectUsers(server Tf2Server) {
	for key := range g.Users {
		user := g.Users[key]
		// delete user if they are older
		if time.Since(user.LastSeen) > 5*time.Minute {
			delete(g.Users, key)
			continue
		}

		user.ServerToConnect = fmt.Sprintf("%s:%d", server.IP, server.Port)
		g.Users[key] = user
	}
}
