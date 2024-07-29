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
