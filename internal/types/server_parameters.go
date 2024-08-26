package types

// Parameters that a Group is looking for
type ServerParameters struct {
	Region                string     `json:"region"`
	GameModes             []GameMode `json:"game_modes"`
	MinPlayers            int        `json:"min_players"`
	MaxPlayers            int        `json:"max_players"`              // Used for blocking out 1000 uncles or 100 player servers
	DisableThousandUncles bool       `json:"thousand_uncles_disabled"` // Thousand Uncles is uncle danes silly server
	MinSpaceAvailable     int        `json:"min_space"`                // Min number of players available in the server beyond group size (so min space 1 means a group of 2 will only join with 3 slots)
}
