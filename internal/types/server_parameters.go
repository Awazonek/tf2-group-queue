package types

// Parameters that a Group is looking for
type ServerParameters struct {
	Region            string     `json:"region"`
	GameModes         []GameMode `json:"game_modes"`
	MinPlayers        int        `json:"min_players"`
	MinSpaceAvailable int        `json:"min_space"` // Min number of players available in the server beyond group size (so min space 1 means a group of 2 will only join with 3 slots)
}
