package types

type Tf2Server struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	IP         string   `json:"ip"`
	Port       int      `json:"port"`
	Mode       GameMode `json:"gamemode"`
	Map        string   `json:"map"`
	Players    int      `json:"players"`
	MaxPlayers int      `json:"max_players"`
	Region     string   `json:"region"`
}
