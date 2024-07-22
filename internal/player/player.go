package player

type Player struct {
	ID        string
	FirstName string
	LastName  string
}

func (p Player) Name() string {
	return p.FirstName + " " + p.LastName
}
