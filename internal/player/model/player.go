package model

type Player struct {
	ID        string
	FirstName string
	LastName  string
	TeamID    string
}

func (p Player) Name() string {
	return p.FirstName + " " + p.LastName
}

type Team struct {
	ID      string
	Player1 Player
	Player2 Player
}

func (t Team) Name() string {
	switch {
	case t.Player1.ID != "" && t.Player2.ID != "":
		return t.Player1.LastName + "/" + t.Player2.LastName
	case t.Player1.ID != "":
		return t.Player1.LastName
	case t.Player2.ID != "":
		return t.Player2.LastName
	default:
		return "Unnamed Team"
	}
}
