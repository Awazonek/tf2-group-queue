package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	types "github.com/awazonek/tf2-group-queue/internal/types"
	"github.com/awazonek/tf2-group-queue/internal/util"
	"github.com/gorilla/mux"
)

type Tf2GroupServer struct {
	router *mux.Router

	groups        map[string]types.Group
	serverList    []types.Tf2Server
	daneSeverList []types.DaneServer
}

func NewServer() Tf2GroupServer {
	r := mux.NewRouter()
	server := Tf2GroupServer{
		router:        r,
		serverList:    make([]types.Tf2Server, 0),
		groups:        make(map[string]types.Group),
		daneSeverList: make([]types.DaneServer, 0),
	}
	server.setRoutes()
	return server
}

func (s *Tf2GroupServer) Start() {
	s.populateDefaultGroup()
	// On start, populate the list of servers.
	s.repeatedServerCall()

	// Start server last
	http.ListenAndServe(":8080", s.router)
}

// set URL routes for clients to join/search
func (s *Tf2GroupServer) setRoutes() {
	s.router.HandleFunc("/join-group/{groupID}", s.joinGroup).Methods("POST")
	s.router.HandleFunc("/search-group/{groupID}", s.searchGroup).Methods("POST")
	s.router.HandleFunc("/reset-user/{groupID}", s.resetUser).Methods("POST")
	s.router.HandleFunc("/user-ping/{groupID}", s.userPing).Methods("GET")
	fs := http.FileServer(http.Dir("./ui"))
	s.router.PathPrefix("/").Handler(fs)
}

// A user will ping if they are in a group, otherwise we don't care
func (s *Tf2GroupServer) resetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupID := vars["groupID"]
	if _, exists := s.groups[groupID]; exists {

		group := s.groups[groupID]
		userId := getIPWithoutPort(r.RemoteAddr)
		_, exists := group.Users[userId]
		if !exists {
			return
		}

		user := types.User{
			ID:           userId,
			LastSeen:     time.Now().UTC(),
			SessionCount: make(map[string]int),
		}

		w.Header().Set("Content-Type", "application/json")
		response := map[string]string{
			"status": "ok",
		}
		json.NewEncoder(w).Encode(response)

		group.Users[userId] = user
	}
}

// A user will ping if they are in a group, otherwise we don't care
func (s *Tf2GroupServer) userPing(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupID := vars["groupID"]
	if _, exists := s.groups[groupID]; exists {

		group := s.groups[groupID]
		userId := getIPWithoutPort(r.RemoteAddr)
		user, exists := group.Users[userId]
		if !exists {
			w.Header().Set("Content-Type", "application/json")
			response := map[string]string{
				"status": "ok",
			}
			json.NewEncoder(w).Encode(response)
			return
		}
		user.LastSeen = time.Now().UTC()

		if user.ConnectedServer != user.ServerToConnect {
			user.ConnectedServer = user.ServerToConnect
			serverAddr := user.ServerToConnect
			count, ok := user.SessionCount[serverAddr]
			if !ok {
				count = 2 // 1 and 2 are reserved for valve things
			}
			user.SessionCount[serverAddr] = count + 1

			w.Header().Set("Content-Type", "application/json")
			response := map[string]string{
				"group":     groupID,
				"server":    serverAddr,
				"quickpick": string(count),
			}
			json.NewEncoder(w).Encode(response)
		} else {

			w.Header().Set("Content-Type", "application/json")
			response := map[string]string{
				"status": "ok",
			}
			json.NewEncoder(w).Encode(response)
		}
		group.Users[userId] = user
	}
}

// A user will join this group
func (s *Tf2GroupServer) joinGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupID := vars["groupID"]

	if _, exists := s.groups[groupID]; exists {
		userId := getIPWithoutPort(r.RemoteAddr)
		user := types.User{
			ID:           userId,
			LastSeen:     time.Now().UTC(),
			SessionCount: make(map[string]int),
		}
		group := s.groups[groupID]
		group.Users[userId] = user

		w.Header().Set("Content-Type", "application/json")
		response := map[string]string{
			"group":        groupID,
			"group_status": "joined",
		}
		json.NewEncoder(w).Encode(response)
	} else {
		http.Error(w, "Invalid group", http.StatusNotFound)
	}
}

// This triggers the search for the group the user is in
func (s *Tf2GroupServer) searchGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupID := vars["groupID"]

	if _, exists := s.groups[groupID]; exists {
		group := s.groups[groupID]
		if group.Searching {
			http.Error(w, "Group already searching", http.StatusBadRequest)
			return // todo: write a confirm?
		}
		userId := getIPWithoutPort(r.RemoteAddr)

		// Only search if the user is in that group
		foundUser := false
		for _, user := range group.Users {
			if user.ID == userId {
				util.Log("\nUser %s is looking from %s", user.ID, userId)
				foundUser = true
				break
			}
		}
		if !foundUser {
			http.Error(w, "You are not in that group", http.StatusNotFound)
			return
		}

		group.Searching = true
		s.groups[groupID] = group

		s.MatchGroup(s.groups[groupID])
		w.Header().Set("Content-Type", "application/json")
		response := map[string]string{
			"group":        groupID,
			"group_status": "searching",
		}
		json.NewEncoder(w).Encode(response)
	} else {
		http.Error(w, "Invalid group", http.StatusNotFound)
	}
}

// Every X seconds we do stuff
func (s *Tf2GroupServer) repeatedServerCall() {
	var err error
	ticker := time.NewTicker(40 * time.Second)
	go func() {
		for range ticker.C {
			err = s.loadAllServers()
			if err != nil {
				log.Printf("Error loading servers: %v", err)
			}
			s.MatchGroups()

			fmt.Printf("Yay we would have retried!")
		}
	}()
}

func (s *Tf2GroupServer) loadAllServers() error {
	hasGroupsSearching := false
	for _, g := range s.groups {
		if g.Searching {
			hasGroupsSearching = true
		}
	}
	// If no groups searching, no reason to load servers
	if !hasGroupsSearching {
		return nil
	}
	err := s.loadDaneServers()
	if err != nil {
		return fmt.Errorf("error loading danes %w", err)
	}
	fmt.Printf("\nUncle dane info %-v", s.serverList)
	return nil
}

// Extracts only the IP address from "IP:Port"
func getIPWithoutPort(remoteAddr string) string {
	ip, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		// If there's an error parsing, fallback to the original address (unlikely case)
		return remoteAddr
	}
	return ip
}
