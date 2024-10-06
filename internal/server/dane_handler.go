package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	types "github.com/awazonek/tf2-group-queue/internal/types"
)

// Get community servers from DANE
func (s *Tf2GroupServer) loadDaneServers() error {

	resp, err := http.Get("https://uncletopia.com/api/servers/state")
	if err != nil {
		return fmt.Errorf("failed to fetch server state: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var serverState types.DaneServerList
	if err := json.NewDecoder(resp.Body).Decode(&serverState); err != nil {
		return fmt.Errorf("failed to decode server response: %w", err)
	}

	// Update the global server list
	s.daneSeverList = serverState.Servers
	s.convertDaneToSearchable()
	log.Printf("Fetched %d Uncletopia servers", len(s.daneSeverList))

	return nil
}

// Turns the dane list into our general server list
func (s *Tf2GroupServer) convertDaneToSearchable() {

	s.serverList = []types.Tf2Server{}
	for _, ds := range s.daneSeverList {
		s.serverList = append(s.serverList, types.Tf2Server{
			ID:         fmt.Sprintf("%d", ds.ServerID),
			Name:       ds.Name,
			IP:         ds.IP,
			Port:       ds.Port,
			Mode:       types.GetGameMode(ds.MapName),
			Map:        ds.MapName,
			Players:    ds.Players,
			MaxPlayers: ds.MaxPlayers,
			Region:     ds.Region,
		})
	}
}
