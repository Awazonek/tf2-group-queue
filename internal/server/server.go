package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/awazonek/tf2-group-queue/internal/globals"
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
	s.router.HandleFunc("/list-groups", s.listGroups).Methods("GET")
	s.router.HandleFunc("/group-metadata", s.groupMetadata).Methods("GET")
	s.router.HandleFunc("/get-group/{groupID}", s.getGroup).Methods("GET")
	s.router.HandleFunc("/join-group/{groupID}", s.joinGroup).Methods("POST")
	s.router.HandleFunc("/create-group/{groupID}", s.createGroup).Methods("POST")
	s.router.HandleFunc("/update-group/{groupID}", s.updateGroup).Methods("POST")
	s.router.HandleFunc("/search-group/{groupID}", s.searchGroup).Methods("POST")
	s.router.HandleFunc("/user-ping/{groupID}", s.userPing).Methods("GET")
	fs := http.FileServer(http.Dir("./ui"))
	s.router.PathPrefix("/").Handler(fs)
}

func (s *Tf2GroupServer) groupMetadata(w http.ResponseWriter, r *http.Request) {
	// TODO:
	// Get all potential options for a group

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"status":         "ok",
		"group_metadata": globals.GetMaxedOutGroup(), // TODO: We could split this up better
	}
	json.NewEncoder(w).Encode(response)
}

func (s *Tf2GroupServer) listGroups(w http.ResponseWriter, r *http.Request) {
	var groupList []types.UserGroupData
	for _, group := range s.groups {
		groupList = append(groupList, types.GroupToUserGroup(group))
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"status": "ok",
		"groups": groupList,
	}
	json.NewEncoder(w).Encode(response)
}

func (s *Tf2GroupServer) getGroup(w http.ResponseWriter, r *http.Request) {
	// TODO
	// Get metadata from a group
	vars := mux.Vars(r)
	groupID := vars["groupID"]
	if group, ok := s.groups[groupID]; ok {

		w.Header().Set("Content-Type", "application/json")
		response := map[string]interface{}{
			"status": "ok",
			"group":  types.GroupToUserGroup(group),
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	http.Error(w, "Invalid group", http.StatusNotFound)
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
			response := map[string]interface{}{
				"status":   "ok",
				"in_group": false,
			}
			json.NewEncoder(w).Encode(response)
			return
		}
		user.LastSeen = time.Now().UTC()

		group.Users[userId] = user

		s.groups[groupID] = group
		w.Header().Set("Content-Type", "application/json")
		response := map[string]interface{}{
			"status":   "ok",
			"in_group": true,
			"group":    types.GroupToUserGroup(group),
		}
		json.NewEncoder(w).Encode(response)
	}
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{
		"status": "ok",
	}
	json.NewEncoder(w).Encode(response)
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

// TODO: Allow users to create groups
func (s *Tf2GroupServer) createGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupID := vars["groupID"]
	if _, exists := s.groups[groupID]; exists {
		http.Error(w, "Group already exists", http.StatusBadRequest)
		return
	}
	var groupData types.UserGroupData
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&groupData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Now we have a user group
	s.CreateGroup(groupData.ID, groupData.Parameters)
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{
		"status":        "ok",
		"group_created": "true",
	}
	json.NewEncoder(w).Encode(response)
}

func (s *Tf2GroupServer) updateGroup(w http.ResponseWriter, r *http.Request) {
	var groupData types.UserGroupData
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&groupData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if group, exists := s.groups[groupData.ID]; exists {
		// Now we have a user group
		group.Parameters = groupData.Parameters
		s.groups[group.ID] = group

		w.Header().Set("Content-Type", "application/json")
		response := map[string]string{
			"status":        "ok",
			"group_updated": "true",
		}
		json.NewEncoder(w).Encode(response)
	}
	http.Error(w, "Group does not exist", http.StatusBadRequest)
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
		group.QueueTries = 0
		s.groups[groupID] = group

		s.MatchGroup(group)
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
	ticker := time.NewTicker(4 * time.Second)
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
		// TODO: For now, we want to change this when we actually have this operational as we want to update the users map
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
