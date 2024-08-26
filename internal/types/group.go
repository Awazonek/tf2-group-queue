package types

import (
	"fmt"
	"strings"
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
	if !util.Contains(g.Parameters.Regions, server.Region) {
		util.Log("Server %s is too far away, wanted region %s and we got %s", server.Name, strings.Join(g.Parameters.Regions, ","), server.Region)
		return false
	}

	if g.Parameters.MinPlayers > server.Players {
		util.Log("Server %s has too few players. Wanted at least %d but got %d", server.Name, g.Parameters.MinPlayers, server.Players)
		return false
	}

	if g.Parameters.MaxPlayers < server.MaxPlayers {
		util.Log("Server %s has too many players allowed. Wanted at most %d but got %d", server.Name, g.Parameters.MaxPlayers, server.MaxPlayers)
		return false
	}

	if g.Parameters.DisableThousandUncles && strings.Contains(server.Name, "One Thousand Uncles") {
		util.Log("Server is a Thousand Uncles server and that is disabled")
		return false
	}

	playerSlots := server.MaxPlayers - server.Players
	wantedSlots := len(g.Users) + g.Parameters.MinSpaceAvailable
	// Check player count
	if playerSlots < wantedSlots {
		util.Log("Server %s has too many players wanted %d slots open and got %d", server.Name, len(g.Users)+g.Parameters.MinSpaceAvailable, (server.MaxPlayers - server.Players))
		return false
	} else {
		util.Log("Server %s has %d slots open and we want %d", server.Name, playerSlots, wantedSlots)
	}

	// If server has correct game mode
	for _, mode := range g.Parameters.GameModes {
		if mode == server.Mode {
			return true
		}
	}
	util.Log("Server %s does not have the correct game mode, it has %d", server.Name, server.Mode)
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
