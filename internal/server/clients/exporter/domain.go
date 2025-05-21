package exporter

type Record struct {
	CreationDate int64  `json:"creationDate"`
	Key          string `json:"key"`
	Value        any    `json:"value,omitempty"`
	Variant      string `json:"variant,omitempty"`
	Reason       string `json:"reason,omitempty"`
	ErrorCode    string `json:"errorCode,omitempty"`
}
