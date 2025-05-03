package cache

type AllFlags struct {
	Flags map[string]FlagState `json:"flags"`
}

func (a *AllFlags) AddFlag(key string, state FlagState) {
	a.Flags[key] = state
}

func NewAllFlags() AllFlags {
	return AllFlags{
		Flags: map[string]FlagState{},
	}
}

type FlagState struct {
	Value   any    `json:"value"`
	Variant string `json:"variant"`
	Reason  string `json:"reason"`
}
