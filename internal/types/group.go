package types

import (
	"strings"

	"github.com/awazonek/tf2-group-queue/internal/util"
)

type Group struct {
	ID         string           `json:"id"`
	Parameters ServerParameters `json:"server_parameters"`
	Users      map[string]User  `json:"users"`
	ServerInfo ServerInfo       `json:"server_info"`
	QueueTries int              `json:"queue_tries"`
	Searching  bool             `json:"searching"`
}

type ServerInfo struct {
	IP       string   `json:"ip"`
	Port     int      `json:"port"`
	Map      string   `json:"map"`
	GameMode GameMode `json:"game_mode"`
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

	playerSlots := server.MaxPlayers - server.Players
	wantedSlots := len(g.Users) + g.Parameters.MinSpaceAvailable
	// Check player count
	if playerSlots < wantedSlots {
		util.Log("Server %s has too many players wanted %d slots open and got %d", server.Name, len(g.Users)+g.Parameters.MinSpaceAvailable, (server.MaxPlayers - server.Players))
		return false
	} else {
		util.Log("Server %s has %d slots open and we want %d", server.Name, playerSlots, wantedSlots)
	}

	if !util.Contains(g.Parameters.Maps, server.Map) {
		if util.Contains(g.Parameters.CustomRules.UnknownMapGameModes, GetGameMode(server.Map)) {
			util.Log("Server %s has map %s but we are allowing it because it is an approved game mode", server.Name, server.Map)
			return true
		}
	} else {
		return true // map is what we want!
	}
	util.Log("Server %s does not have the correct game mode, it has %d", server.Name, server.Mode)
	return false
}

func (g *Group) getUserFacingGroupData() UserGroupData {
	return UserGroupData{
		ID:         g.ID,
		Parameters: g.Parameters,
		ServerInfo: g.ServerInfo,
		Searching:  g.Searching,
	}
}
