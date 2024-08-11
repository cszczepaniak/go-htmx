package players

type GetPlayerOpts struct {
	WithoutTeam bool
}

type GetPlayerOpt func(GetPlayerOpts) GetPlayerOpts

func WithoutTeam() GetPlayerOpt {
	return func(gpo GetPlayerOpts) GetPlayerOpts {
		gpo.WithoutTeam = true
		return gpo
	}
}

type GetTeamOpts struct {
	WithoutDivision bool
	DivisionID      string
}

type GetTeamOpt func(GetTeamOpts) GetTeamOpts

func WithoutDivision() GetTeamOpt {
	return func(gpo GetTeamOpts) GetTeamOpts {
		gpo.WithoutDivision = true
		return gpo
	}
}

func InDivision(id string) GetTeamOpt {
	return func(gpo GetTeamOpts) GetTeamOpts {
		gpo.WithoutDivision = false
		gpo.DivisionID = id
		return gpo
	}
}
