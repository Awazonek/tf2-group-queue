package types

type DaneServerList struct {
	Servers []DaneServer `json:"servers"`
}

type DaneServer struct {
	ServerID   int      `json:"server_id"`
	Host       string   `json:"host"`
	Port       int      `json:"port"`
	IP         string   `json:"ip"`
	Name       string   `json:"name"`
	NameShort  string   `json:"name_short"`
	Region     string   `json:"region"`
	CC         string   `json:"cc"`
	Players    int      `json:"players"`
	MaxPlayers int      `json:"max_players"`
	Bots       int      `json:"bots"`
	MapName    string   `json:"map"`
	GameTypes  []string `json:"game_types"`
	Latitude   float32  `json:"latitude"`
	Longitude  float32  `json:"longitude"`
	Distance   float64  `json:"distance"`
}
