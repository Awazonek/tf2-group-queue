package types

import "time"

type User struct {
	ID           string         `json:"id"`
	LastSeen     time.Time      `json:"timestamp"`    // Last seen, kicking them out after X minutes
	SessionCount map[string]int `json:"session_data"` // history of user connecting to server
}

type UserGroupData struct {
	// Data the user will receive about their group on ping
	ID         string           `json:"id"`
	Parameters ServerParameters `json:"server_parameters"`
	ServerInfo ServerInfo       `json:"server_info,omitempty"`
	Searching  bool             `json:"searching"`
	QueueTries int              `json:"query_tries"`
}

func GroupToUserGroup(group Group) UserGroupData {
	return UserGroupData{
		ID:         group.ID,
		ServerInfo: group.ServerInfo,
		Searching:  group.Searching,
		Parameters: group.Parameters,
		QueueTries: group.QueueTries,
	}
}
