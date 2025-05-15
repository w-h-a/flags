package cache

type AllFlags struct {
	Flags        []FlagState `json:"flags"`
	ErrorCode    string      `json:"errorCode,omitempty"`
	ErrorMessage string      `json:"errorMessage,omitempty"`
}

func (a *AllFlags) AddFlag(state FlagState) {
	a.Flags = append(a.Flags, state)
}

func NewAllFlags() AllFlags {
	return AllFlags{
		Flags: []FlagState{},
	}
}

type FlagState struct {
	Key          string `json:"key"`
	Value        any    `json:"value,omitempty"`
	Variant      string `json:"variant,omitempty"`
	Reason       string `json:"reason,omitempty"`
	ErrorCode    string `json:"errorCode,omitempty"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}
