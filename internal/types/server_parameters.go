package types

// Parameters that a Group is looking for
type ServerParameters struct {
	Regions           []string    `json:"region"`
	Maps              []string    `json:"maps"` // List of maps to pick
	MinPlayers        int         `json:"min_players"`
	MaxPlayers        int         `json:"max_players"`  // Used for blocking out 1000 uncles or 100 player servers
	CustomRules       CustomRules `json:"custom_rules"` // Stores special info like custom maps
	MinSpaceAvailable int         `json:"min_space"`    // Min number of players available in the server beyond group size (so min space 1 means a group of 2 will only join with 3 slots)
}

type CustomRules struct {
	DisableThousandUncles bool       `json:"thousand_uncles_disabled"` // Thousand Uncles is uncle danes silly server
	UnknownMapGameModes   []GameMode `json:"allow_unknown_maps"`       // Allow maps we don't recognize if they are using these game modes
}
