package server

import (
	"sort"

	"github.com/awazonek/tf2-group-queue/internal/types"
	"github.com/awazonek/tf2-group-queue/internal/util"
)

func (s *Tf2GroupServer) CreateGroup(ID string, Parameters types.ServerParameters) {
	s.groups[ID] = types.Group{
		ID:         ID,
		Parameters: Parameters,
		Searching:  false,
		Users:      make(map[string]types.User),
	}
}

func (s *Tf2GroupServer) MatchGroups() {
	// for each searching group
	// CHeck if a map matches parameters
	util.Log("Checking %d groups", len(s.groups))
	for _, g := range s.groups {
		s.MatchGroup(g)
	}
}

func (s *Tf2GroupServer) MatchGroup(group types.Group) {

	if group.Searching {
		util.Log("\nSearching for group %s", group.ID)
		validServers := []types.Tf2Server{}
		for _, srv := range s.serverList {
			if group.MatchesServer(srv) {
				util.Log("OMG A VALID SERVER! %s", srv.Name)
				validServers = append(validServers, srv)
			}
		}

		if len(validServers) > 0 {
			sort.Slice(validServers, func(i, j int) bool {
				return validServers[i].Players > validServers[j].Players
			})

			group.Searching = false
			group.ConnectUsers(validServers[0])
			s.groups[group.ID] = group
		} else {
			util.Log("No valid server found for group %s", group.ID)
		}
	} else {

		util.Log("\nGroup %s is not searching", group.ID)
	}
}

func (s *Tf2GroupServer) populateDefaultGroup() {
	s.CreateGroup("GuuzTesting", types.ServerParameters{
		Regions: []string{"us-east", "us-central"},
		GameModes: []types.GameMode{
			types.Payload,
			types.AttackDefend,
		},
		MinPlayers:            8,
		MaxPlayers:            32,
		DisableThousandUncles: true,
		MinSpaceAvailable:     2,
	})
}
