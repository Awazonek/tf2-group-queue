package types

type GameMode int

const (
	Unknown GameMode = iota
	Payload
	CaptureTheFlag
	ControlPoint
	KingOfTheHill
	AttackDefend
	PayloadRace
	PlayerDestruction
	Other
)

func (gm GameMode) String() string {
	switch gm {
	case Payload:
		return "Payload"
	case CaptureTheFlag:
		return "Capture the Flag"
	case ControlPoint:
		return "Control Point"
	case KingOfTheHill:
		return "King of the Hill"
	case AttackDefend:
		return "Attack/Defend"
	case PayloadRace:
		return "Payload Race"
	case PlayerDestruction:
		return "Player Destruction"
	default:
		return "Unknown"
	}
}

func GetGameMode(mapName string) GameMode {
	if len(mapName) < 4 {
		return Unknown
	}
	switch {
	case mapName[:3] == "pl_":
		return Payload
	case mapName[:4] == "ctf_":
		return CaptureTheFlag
	case mapName[:3] == "cp_":
		return ControlPoint
	case mapName[:5] == "koth_":
		return KingOfTheHill
	case mapName[:4] == "plr_":
		return PayloadRace
	case mapName[:3] == "pd_":
		return PlayerDestruction
	case mapName[:3] == "ad_":
		return AttackDefend
	default:
		return Unknown
	}
}
