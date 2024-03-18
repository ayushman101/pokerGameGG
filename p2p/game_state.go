package p2p

type GameVariant uint8

const (
	TexasHoldem GameVariant = iota
	Other
)

func (g GameVariant) String() string {
	switch g {
	case TexasHoldem:
		return "TEXASHOLDEM"
	case Other:
		return "OTHER"
	default:
		return "UNKNOWN"
	}
}

type GameStatus uint32

const (
	GameStatusWaiting GameStatus = iota
	GameStatusDealing
	GameStatusPreFlop
	GameStatusFlop
	GameStatusTurn
	GameStatusRiver
)

func (g GameStatus) String() string {
	switch g {
	case GameStatusDealing:
		return "DEALING"
	case GameStatusWaiting:
		return "WAITING"
	case GameStatusPreFlop:
		return "PRE FLOP"
	case GameStatusFlop:
		return "FLOP"
	case GameStatusTurn:
		return "TURN"
	case GameStatusRiver:
		return "RIVER"
	default:
		return "unknown"
	}
}
